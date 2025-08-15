package iexcel

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

// TestUser 是用于测试的示例结构体
type TestUser struct {
	ID   int    `xlsx:"用户ID"`
	Name string `xlsx:"姓名"`
	Age  int    `xlsx:"年龄"`
	// 这个字段会被忽略
	Password string `xlsx:"-"`
	Remark   string `xlsx:"备注"`
}

// TestProduct 是另一个用于测试的示例结构体
type TestProduct struct {
	SKU   string  `xlsx:"产品SKU"`
	Price float64 `xlsx:"价格"`
}

// assertHeaders 是一个测试辅助函数，用于验证工作表的表头
func assertHeaders(t *testing.T, f *excelize.File, sheetName string, expectedHeaders []string) {
	t.Helper()
	for i, expectedHeader := range expectedHeaders {
		colName, err := excelize.ColumnNumberToName(i + 1)
		assert.NoError(t, err)
		cellAddress := colName + "1"
		actualHeader, err := f.GetCellValue(sheetName, cellAddress)
		assert.NoError(t, err)
		assert.Equal(t, expectedHeader, actualHeader, "Header in cell %s should match", cellAddress)
	}
}

// assertRow 是一个测试辅助函数，用于验证特定行的数据
func assertRow(t *testing.T, allRows [][]string, rowIndex int, expectedData []string) {
	t.Helper()
	assert.GreaterOrEqualf(t, len(allRows), rowIndex+1, "Should have at least %d rows", rowIndex+1)
	actualRow := allRows[rowIndex]
	// Pad the actual row with empty strings to match the expected length.
	// This is to compensate for excelize's GetRows behavior which omits trailing empty cells.
	if len(actualRow) < len(expectedData) {
		paddedActual := make([]string, len(expectedData))
		copy(paddedActual, actualRow)
		actualRow = paddedActual
	}
	assert.Equal(t, expectedData, actualRow, "Row data at index %d should match", rowIndex)
}

// assertSheets 是一个辅助函数，用于验证工作簿包含且仅包含指定名称的工作表
func assertSheets(t *testing.T, f *excelize.File, expectedSheetNames []string) {
	t.Helper()
	sheetList := f.GetSheetList()
	assert.ElementsMatch(t, expectedSheetNames, sheetList, "The workbook's sheet list should match the expected list")
}

// TestNew 测试构造函数
func TestNew(t *testing.T) {
	builder := New()
	assert.NotNil(t, builder, "New() should return a non-nil builder")
	assert.NoError(t, builder.Error(), "A new builder should have no error")
}

// TestAddSheetAndSave_Success 测试最核心的成功路径
func TestAddSheetAndSave_Success(t *testing.T) {
	testData := []TestUser{
		{ID: 1, Name: "Alice", Age: 30, Password: "123", Remark: "VIP"},
		{ID: 2, Name: "Bob", Age: 25, Password: "456", Remark: ""},
	}
	filepath := "test_output_success.xlsx"
	defer os.Remove(filepath)

	err := New().AddSheet("用户列表", testData).Save(filepath)
	assert.NoError(t, err)
	assert.FileExists(t, filepath)

	f, err := excelize.OpenFile(filepath)
	assert.NoError(t, err)
	defer f.Close()

	assertSheets(t, f, []string{"用户列表"})

	expectedHeaders := []string{"用户ID", "姓名", "年龄", "备注"}
	assertHeaders(t, f, "用户列表", expectedHeaders)

	rows, err := f.GetRows("用户列表")
	assert.NoError(t, err)
	assert.Len(t, rows, len(testData)+1, "Should be 1 header row + 2 data rows")

	assertRow(t, rows, 1, []string{"1", "Alice", "30", "VIP"})
	assertRow(t, rows, 2, []string{"2", "Bob", "25", ""})
}

// TestAddSheet_EmptySlice 测试传入空切片的情况
func TestAddSheet_EmptySlice(t *testing.T) {
	testData := []TestUser{}
	filepath := "test_output_empty.xlsx"
	defer os.Remove(filepath)

	err := New().AddSheet("空列表", testData).Save(filepath)
	assert.NoError(t, err)
	assert.FileExists(t, filepath)

	f, err := excelize.OpenFile(filepath)
	assert.NoError(t, err)
	defer f.Close()

	assertSheets(t, f, []string{"空列表"})

	expectedHeaders := []string{"用户ID", "姓名", "年龄", "备注"}
	assertHeaders(t, f, "空列表", expectedHeaders)

	rows, err := f.GetRows("空列表")
	assert.NoError(t, err)
	assert.Len(t, rows, 1, "Should only have 1 header row")
}

// TestAddSheet_NonSliceInput 测试传入非切片类型的错误情况
func TestAddSheet_NonSliceInput(t *testing.T) {
	builder := New().AddSheet("无效数据", TestUser{ID: 1, Remark: "test"})
	assert.Error(t, builder.Error(), "Should return an error for non-slice input")
	assert.Contains(t, builder.Error().Error(), "records must be a slice of structs")
}

// TestChain_WithModifications 测试链式调用和修改操作
func TestChain_WithModifications(t *testing.T) {
	users := []TestUser{{ID: 1, Name: "User"}}
	filepath := "test_output_mods.xlsx"
	defer os.Remove(filepath)

	// 测试一个更简单的场景：添加一个工作表然后将其删除
	builder := New().
		AddSheet("用户", users).
		RemoveSheet("用户")

	err := builder.Save(filepath)
	assert.NoError(t, err)
	assert.NoError(t, builder.Error())

	f, err := excelize.OpenFile(filepath)
	assert.NoError(t, err)
	defer f.Close()

	// 验证工作表列表
	// 移除添加的表后，应该只剩下 excelize 默认的 "Sheet1"
	assertSheets(t, f, []string{"Sheet1"})
}

// TestWrite_Success 测试写入到 io.Writer
func TestWrite_Success(t *testing.T) {
	testData := []TestUser{{ID: 1, Name: "Writer"}}
	var buffer bytes.Buffer

	err := New().AddSheet("BufferTest", testData).Write(&buffer)
	assert.NoError(t, err)
	assert.Greater(t, buffer.Len(), 0, "Buffer should contain data after write")

	f, err := excelize.OpenReader(&buffer)
	assert.NoError(t, err)
	defer f.Close()

	assertSheets(t, f, []string{"BufferTest"})

	expectedHeaders := []string{"用户ID", "姓名", "年龄", "备注"}
	assertHeaders(t, f, "BufferTest", expectedHeaders)

	rows, err := f.GetRows("BufferTest")
	assert.NoError(t, err)
	assert.Len(t, rows, 2, "Should be 1 header row + 1 data row")
	assertRow(t, rows, 1, []string{"1", "Writer", "0", ""})
}
