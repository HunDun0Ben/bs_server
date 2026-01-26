package main

/*
#cgo CFLAGS: -I/home/ben/workspace/go-wks/libsvm
#cgo LDFLAGS: -L/home/ben/workspace/go-wks/libsvm -lsvm
#include <stdlib.h>
#include "svm.h"
*/
import "C"

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unsafe"
)

// 读取CSV文件，返回 bow 矩阵和标签数组.
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

	// 预分配容量到训练大小, 避免出现一直扩容.
	features := make([][]float64, 0, len(lines))
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

// MinMaxScaling performs min-max scaling on features to range [0, 1]
func MinMaxScaling(features [][]float64) ([][]float64, []float64, []float64) {
	if len(features) == 0 {
		return features, nil, nil
	}
	dim := len(features[0])
	minVals := make([]float64, dim)
	maxVals := make([]float64, dim)

	// Initialize min/max
	for i := 0; i < dim; i++ {
		minVals[i] = features[0][i]
		maxVals[i] = features[0][i]
	}

	// Find min/max
	for _, feat := range features {
		for i, v := range feat {
			if v < minVals[i] {
				minVals[i] = v
			}
			if v > maxVals[i] {
				maxVals[i] = v
			}
		}
	}

	// Scale
	scaled := make([][]float64, len(features))
	for i, feat := range features {
		scaled[i] = make([]float64, dim)
		for j, v := range feat {
			if maxVals[j] == minVals[j] {
				scaled[i][j] = 0
			} else {
				scaled[i][j] = (v - minVals[j]) / (maxVals[j] - minVals[j])
			}
		}
	}
	return scaled, minVals, maxVals
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
	rand.Seed(time.Now().UnixNano())

	features, labels, err := readCSV("../script/data.csv")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d samples.\n", len(labels))

	// 1. 数据归一化 (Min-Max Scaling)
	// SVM 对数据缩放非常敏感，通常缩放到 [0,1] 或 [-1,1]
	features, _, _ = MinMaxScaling(features)
	fmt.Println("Data scaled to [0, 1].")

	// 2. 数据打乱 (Shuffle)
	perm := rand.Perm(len(features))
	shuffledFeatures := make([][]float64, len(features))
	shuffledLabels := make([]float64, len(labels))
	for i, v := range perm {
		shuffledFeatures[i] = features[v]
		shuffledLabels[i] = labels[v]
	}
	features = shuffledFeatures
	labels = shuffledLabels
	fmt.Println("Data shuffled.")

	// 3. 划分训练集和验证集 (80% Train, 20% Validation)
	trainSize := int(float64(len(labels)) * 0.8)
	trainFeatures := features[:trainSize]
	trainLabels := labels[:trainSize]
	valFeatures := features[trainSize:]
	valLabels := labels[trainSize:]

	fmt.Printf("Train set size: %d, Validation set size: %d\n", len(trainLabels), len(valLabels))

	l := C.int(len(trainLabels))
	// 转换标签
	y := (*C.double)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(C.double(0)))))
	defer C.free(unsafe.Pointer(y))
	yArr := (*[1 << 30]C.double)(unsafe.Pointer(y))[:l:l]
	for i := 0; i < int(l); i++ {
		yArr[i] = C.double(trainLabels[i])
	}

	// 转换特征
	x := (**C.struct_svm_node)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(uintptr(0)))))
	defer C.free(unsafe.Pointer(x))

	// 用于跟踪所有分配的 svm_node，以便最后释放
	var allocatedNodes []*C.struct_svm_node
	defer func() {
		for _, node := range allocatedNodes {
			C.free(unsafe.Pointer(node))
		}
	}()

	xArr := (*[1 << 30]*C.struct_svm_node)(unsafe.Pointer(x))[:l:l]
	for i := 0; i < int(l); i++ {
		node := createSvmNode(trainFeatures[i])
		xArr[i] = node
		allocatedNodes = append(allocatedNodes, node)
	}

	// 构造问题结构体
	var prob C.struct_svm_problem
	prob.l = l
	prob.y = y
	prob.x = x

	CValues := []float64{0.1, 1, 10, 100, 1000}
	gammaValues := []float64{0.01, 0.1, 1.0, 10.0}

	var bestModel *C.struct_svm_model
	var bestAccuracy float64 = -1.0
	var bestC, bestGamma float64

	for _, c := range CValues {
		for _, gamma := range gammaValues {
			// fmt.Printf("Training with C=%.2f, gamma=%.2f... ", c, gamma)
			param := newSvmParameter(c, gamma)
			model := traning(&prob, param)
			if model == nil {
				continue
			}

			// 使用验证集评估模型
			var correctCount int
			for i := 0; i < len(valLabels); i++ {
				predictedLabel := predict(model, valFeatures[i])
				if predictedLabel == valLabels[i] {
					correctCount++
				}
			}
			accuracy := float64(correctCount) / float64(len(valLabels))
			// fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)

			if accuracy > bestAccuracy {
				fmt.Printf("New best found: C=%.2f, gamma=%.2f, Acc=%.2f%%\n", c, gamma, accuracy*100)
				// 释放之前保存的最佳模型
				if bestModel != nil {
					C.svm_free_and_destroy_model(&bestModel)
				}
				bestAccuracy = accuracy
				bestC = c
				bestGamma = gamma
				bestModel = model
			} else {
				// 如果当前模型不是最佳模型，则立即释放它
				C.svm_free_and_destroy_model(&model)
			}
		}
	}

	fmt.Printf("\nGrid search finished.\n")
	fmt.Printf("Best C=%.2f, gamma=%.2f with accuracy=%.2f%%\n", bestC, bestGamma, bestAccuracy*100)

	// // 保存最佳模型
	// if bestModel != nil {
	// 	modelFileName := C.CString("best_model.svm")
	// 	defer C.free(unsafe.Pointer(modelFileName))
	// 	if C.svm_save_model(modelFileName, bestModel) != 0 {
	// 		fmt.Println("Failed to save model")
	// 	} else {
	// 		fmt.Println("Best model saved to best_model.svm")
	// 	}
	// 	C.svm_free_and_destroy_model(&bestModel)
	// }
}

func predict(model *C.struct_svm_model, feature []float64) float64 {
	// 1. 为待预测样本创建 svm_node
	node := createSvmNode(feature)
	defer C.free(unsafe.Pointer(node)) // 关键：确保释放内存

	// 2. 执行预测
	predictedLabel := C.svm_predict(model, node)

	return float64(predictedLabel)
}

func traning(prob *C.struct_svm_problem, param *C.struct_svm_parameter) *C.struct_svm_model {
	// 检查参数有效性
	errStr := C.svm_check_parameter(prob, param)
	if errStr != nil {
		fmt.Println("参数错误:", C.GoString(errStr))
		return nil
	}
	// 训练模型
	// C.svm_train returns a pointer to a model allocated by malloc, we must free it later.
	model := C.svm_train(prob, param)
	return model
}

// 伪代码，示范参数结构体初始化.
func newSvmParameter(c, gamma float64) *C.struct_svm_parameter {
	var param C.struct_svm_parameter
	param.svm_type = C.C_SVC
	param.kernel_type = 2 // RBF
	param.gamma = C.double(gamma)
	param.C = C.double(c)
	param.cache_size = 100
	param.eps = 0.001
	param.shrinking = 1
	param.probability = 0
	return &param
}
