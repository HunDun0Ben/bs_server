package iexcel

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

// IExcelBuilder 定义了一套用于以链式方式创建和操作 Excel 工作簿的通用接口。
//
// 它旨在提供一个比底层库更高层次、更易用的 API。
// 错误处理机制：链式调用中的方法不直接返回错误，错误会被记录在内部。
// 用户可以在调用链的任何地方，或在最后通过 Error() 方法检查是否有错误发生。
type IExcelBuilder interface {
	// AddSheet 是一个高层方法，用于添加一个新工作表并直接用结构体切片填充数据。
	// 这是最常用的快捷操作。
	// records 参数必须是结构体切片。
	AddSheet(sheetName string, records any) IExcelBuilder

	// RemoveSheet 从工作簿中移除一个工作表。
	RemoveSheet(sheetName string) IExcelBuilder

	// SetCellValue 提供了一个更细粒度的操作，用于设置特定单元格的值。
	// axis 参数是单元格坐标，例如 "A1", "B2"。
	SetCellValue(sheetName, axis string, value any) IExcelBuilder

	// SetActiveSheet 设置打开 Excel 文件时默认显示的活动工作表。
	SetActiveSheet(sheetName string) IExcelBuilder

	// Save 将工作簿保存到指定的文件路径。
	// 这是一个终端操作，会结束链式调用并返回最终的错误（如果有）。
	Save(filepath string) error

	// Write 将工作簿内容写入到一个 io.Writer。
	// 这是一个终端操作，适用于网络传输等场景。
	Write(writer io.Writer) error

	// Error 返回在链式调用过程中累计发生的第一个错误。
	// 如果没有错误，则返回 nil。
	Error() error
}

// excelBuilder 是 IExcelBuilder 接口的内部具体实现。
// 它不被导出，以鼓励外部代码仅依赖于接口。
type excelBuilder struct {
	file       *excelize.File
	err        error
	isPristine bool // 用于跟踪是否是第一次添加工作表
}

// New 是 Excel 构建器的构造函数。
// 它返回 IExcelBuilder 接口，作为所有链式调用的起点。
func New() IExcelBuilder {
	f := excelize.NewFile()
	// excelize 默认会创建一个 "Sheet1"。我们不删除它，而是在第一次 AddSheet 时重命名它。
	return &excelBuilder{
		file:       f,
		err:        nil,
		isPristine: true, // 初始状态为 pristine，表示是全新的构建器
	}
}

// Error 返回在链式调用过程中累计发生的第一个错误。
func (eb *excelBuilder) Error() error {
	return eb.err
}

// AddSheet 实现 IExcelBuilder 接口的 AddSheet 方法
func (eb *excelBuilder) AddSheet(sheetName string, records any) IExcelBuilder {
	// 如果构建器已经处于错误状态，则直接返回，不再执行后续操作。
	if eb.err != nil {
		return eb
	}

	// --- 1. 输入校验 ---
	v := reflect.ValueOf(records)
	if v.Kind() != reflect.Slice {
		eb.err = errors.New("records must be a slice of structs")
		return eb
	}

	elemType := v.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		eb.err = errors.New("records must be a slice of structs")
		return eb
	}

	// --- 2. 创建或重命名工作表 ---
	if eb.isPristine {
		// 如果是第一次添加，直接重命名默认的 "Sheet1"
		eb.err = eb.file.SetSheetName("Sheet1", sheetName)
		if eb.err == nil {
			// 成功后，将状态标记为不再是 pristine
			eb.isPristine = false
		}
	} else {
		// 否则，创建一个新的工作表
		_, eb.err = eb.file.NewSheet(sheetName)
	}
	if eb.err != nil {
		return eb
	}

	// --- 3. 提取并排序要导出的字段 ---
	type fieldInfo struct {
		Index int
		Tag   string
	}
	var orderedFields []fieldInfo
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("xlsx")
		if tag == "" || tag == "-" {
			continue // 跳过不需要写入的字段
		}
		orderedFields = append(orderedFields, fieldInfo{Index: i, Tag: tag})
	}

	// --- 4. 写入表头 ---
	for colIndex, field := range orderedFields {
		var colName string
		colName, eb.err = excelize.ColumnNumberToName(colIndex + 1)
		if eb.err != nil {
			return eb
		}
		position := colName + "1"
		eb.err = eb.file.SetCellValue(sheetName, position, field.Tag)
		if eb.err != nil {
			return eb
		}
	}

	// --- 5. 写入数据行 ---
	if v.Len() == 0 {
		return eb
	}

	for rowIndex := 0; rowIndex < v.Len(); rowIndex++ {
		elemValue := v.Index(rowIndex)
		for colIndex, field := range orderedFields {
			var colName string
			// 错误在上面已经处理过，这里可以忽略
			colName, _ = excelize.ColumnNumberToName(colIndex + 1)
			position := fmt.Sprintf("%s%d", colName, rowIndex+2) // rowIndex+2 是因为第1行是表头
			eb.err = eb.file.SetCellValue(sheetName, position, elemValue.Field(field.Index).Interface())
			if eb.err != nil {
				return eb
			}
		}
	}

	return eb
}

// Save 实现 IExcelBuilder 接口的 Save 方法
func (eb *excelBuilder) Save(filepath string) error {
	// 如果链式调用过程中已发生错误，直接返回该错误
	if eb.err != nil {
		return eb.err
	}

	sheetList := eb.file.GetSheetList()
	// 如果用户没有添加任何 sheet，或所有 sheet 都被删除了，则创建一个默认的
	if len(sheetList) == 0 {
		eb.file.NewSheet("Sheet1")
	} else {
		// 否则，将活动工作表设置为第一个，以保证良好的用户体验
		eb.file.SetActiveSheet(0)
	}

	return eb.file.SaveAs(filepath)
}

// RemoveSheet 实现 IExcelBuilder 接口的 RemoveSheet 方法
func (eb *excelBuilder) RemoveSheet(sheetName string) IExcelBuilder {
	if eb.err != nil {
		return eb
	}

	// excelize 不允许删除最后一个工作表。
	// 如果我们正要删除最后一个工作表，我们必须先创建一个新的，
	// 然后删除目标工作表，最后将新的重命名为默认名称。
	sheetList := eb.file.GetSheetList()
	if len(sheetList) == 1 && sheetList[0] == sheetName {
		// 这是最后一个工作表，应用变通方案。
		// 为了避免超出31个字符的长度限制，我们替换前缀而不是追加后缀。
		var tempSheetName string
		if len(sheetName) > 3 {
			tempSheetName = "tmp" + sheetName[3:]
		} else {
			tempSheetName = "tmp" + sheetName
		}

		_, eb.err = eb.file.NewSheet(tempSheetName)
		if eb.err != nil {
			return eb
		}

		// 现在有两个工作表，我们可以安全地删除目标工作表。
		eb.file.DeleteSheet(sheetName)

		// 将临时工作表重命名为默认的 "Sheet1"。
		eb.err = eb.file.SetSheetName(tempSheetName, "Sheet1")
		if eb.err != nil {
			return eb
		}

		// 此操作后，构建器处于等同于 New() 的状态。
		// 下一个 AddSheet 应该重命名这个 "Sheet1"。
		eb.isPristine = true
	} else {
		// 如果有多个工作表，或者要删除的工作表不存在，
		// 直接删除即可。
		eb.file.DeleteSheet(sheetName)
	}

	return eb
}

// SetCellValue 实现 IExcelBuilder 接口的 SetCellValue 方法
func (eb *excelBuilder) SetCellValue(sheetName, axis string, value any) IExcelBuilder {
	if eb.err != nil {
		return eb
	}
	eb.err = eb.file.SetCellValue(sheetName, axis, value)
	return eb
}

// SetActiveSheet 实现 IExcelBuilder 接口的 SetActiveSheet 方法
func (eb *excelBuilder) SetActiveSheet(sheetName string) IExcelBuilder {
	if eb.err != nil {
		return eb
	}
	var index int
	index, eb.err = eb.file.GetSheetIndex(sheetName)
	if eb.err != nil {
		return eb
	}
	eb.file.SetActiveSheet(index)
	return eb
}

// Write 实现 IExcelBuilder 接口的 Write 方法
func (eb *excelBuilder) Write(writer io.Writer) error {
	if eb.err != nil {
		return eb.err
	}

	sheetList := eb.file.GetSheetList()
	if len(sheetList) == 0 {
		eb.file.NewSheet("Sheet1")
	} else {
		eb.file.SetActiveSheet(0)
	}

	return eb.file.Write(writer)
}
