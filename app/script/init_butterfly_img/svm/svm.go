package main

/*
#cgo CFLAGS: -I/home/ben/workspace/go_wks/libsvm
#cgo LDFLAGS: -L/home/ben/workspace/go_wks/libsvm -lsvm
#include <stdlib.h>
#include "svm.h"
*/
import "C"

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

// 读取CSV文件，返回二维特征矩阵和标签数组.
func readCSV(file string) ([][]float64, []float64, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	features := make([][]float64, 0, len(lines)) // 预分配容量
	labels := make([]float64, 0, len(lines))

	for _, line := range lines {
		n := len(line)
		if n < 2 {
			continue
		}
		feat := make([]float64, n-1)
		for i := 0; i < n-1; i++ {
			fv, err := strconv.ParseFloat(line[i], 64)
			if err != nil {
				return nil, nil, err
			}
			feat[i] = fv
		}
		label, err := strconv.ParseFloat(line[n-1], 64)
		if err != nil {
			return nil, nil, err
		}

		features = append(features, feat)
		labels = append(labels, label)
	}
	return features, labels, nil
}

// 创建 libsvm 的 svm_node 数组，末尾 index=0 作为结束标志.
func createSvmNode(feature []float64) *C.struct_svm_node {
	n := len(feature)
	size := C.size_t(n+1) * C.size_t(unsafe.Sizeof(C.struct_svm_node{}))
	ptr := C.malloc(size)
	nodes := (*[1 << 30]C.struct_svm_node)(ptr)[: n+1 : n+1]
	for i := 0; i < n; i++ {
		nodes[i].index = C.int(i + 1)
		nodes[i].value = C.double(feature[i])
	}
	nodes[n].index = -1
	return (*C.struct_svm_node)(ptr)
}

func main() {
	features, labels, err := readCSV("data.csv")
	if err != nil {
		panic(err)
	}
	l := C.int(len(labels))
	// 转换标签
	y := (*C.double)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(C.double(0)))))
	defer C.free(unsafe.Pointer(y))
	yArr := (*[1 << 30]C.double)(unsafe.Pointer(y))[:l:l]
	for i := 0; i < int(l); i++ {
		yArr[i] = C.double(labels[i])
	}

	// 转换特征
	x := (**C.struct_svm_node)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(uintptr(0)))))
	defer C.free(unsafe.Pointer(x))

	xArr := (*[1 << 30]*C.struct_svm_node)(unsafe.Pointer(x))[:l:l]
	for i := 0; i < int(l); i++ {
		xArr[i] = createSvmNode(features[i])
	}

	// 构造问题结构体
	var prob C.struct_svm_problem
	prob.l = l
	prob.y = y
	prob.x = x

	nuValues := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
	gammaValues := []float64{0.01, 0.1, 1.0, 10.0}

	for _, nu := range nuValues {
		for _, gamma := range gammaValues {
			param := newSvmParameter(nu, gamma)
			traning(&prob, param, features, labels)
		}
	}
}

func traning(prob *C.struct_svm_problem, param *C.struct_svm_parameter, features [][]float64, labels []float64) {
	// 检查参数有效性
	errStr := C.svm_check_parameter(prob, param)
	if errStr != nil {
		fmt.Println("参数错误:", C.GoString(errStr))
		return
	}
	// 训练模型
	model := C.svm_train(prob, param)
	defer C.svm_free_and_destroy_model(&model) //nolint:gocritic
	var ct int
	for i := 0; i < 120; i++ {
		testNode := createSvmNode(features[0+i])
		pred := C.svm_predict(model, testNode)
		if float64(pred) == labels[0+i] {
			ct++
		}
	}
	fmt.Printf("样本预测类别: %d", ct)
}

// 伪代码，示范参数结构体初始化.
func newSvmParameter(nu, gamma float64) *C.struct_svm_parameter {
	var param C.struct_svm_parameter
	param.svm_type = C.C_SVC
	param.kernel_type = 2
	param.nu = C.double(nu)
	param.gamma = C.double(gamma)
	param.C = 1
	param.cache_size = 100
	param.eps = 0.001
	param.shrinking = 1
	param.probability = 0
	return &param
}
