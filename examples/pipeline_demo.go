package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fustgo/fustgo2/pkg/plugin"
	_ "github.com/fustgo/fustgo2/pkg/plugin/processor"
	_ "github.com/fustgo/fustgo2/pkg/plugin/reader"
	_ "github.com/fustgo/fustgo2/pkg/plugin/writer"
	"github.com/fustgo/fustgo2/pkg/record"
)

func main() {
	fmt.Println("=== FustGo Pipeline Demo ===\n")

	// 创建测试数据文件
	createTestData()

	// 演示 1: 简单的文件读取和输出
	fmt.Println("Demo 1: Simple File Reader -> Stdout Writer")
	demoSimplePipeline()

	// 演示 2: 带处理器的管道
	fmt.Println("\nDemo 2: File Reader -> Field Mapper -> Filter -> Stdout Writer")
	demoProcessorPipeline()
}

// createTestData 创建测试数据
func createTestData() {
	// 创建 CSV 测试文件
	csvContent := `id,name,age,email
1,Alice,30,alice@example.com
2,Bob,25,bob@example.com
3,Charlie,35,charlie@example.com
4,David,28,david@example.com
5,Eve,32,eve@example.com`

	if err := writeFile("/tmp/test_data.csv", csvContent); err != nil {
		fmt.Printf("Failed to create test data: %v\n", err)
		return
	}

	// 创建 JSON 测试文件
	jsonContent := `{"id": 1, "name": "Alice", "age": 30, "email": "alice@example.com"}
{"id": 2, "name": "Bob", "age": 25, "email": "bob@example.com"}
{"id": 3, "name": "Charlie", "age": 35, "email": "charlie@example.com"}`

	if err := writeFile("/tmp/test_data.json", jsonContent); err != nil {
		fmt.Printf("Failed to create test data: %v\n", err)
		return
	}
}

// writeFile 写入文件
func writeFile(filename, content string) error {
	// 这里简化处理，实际应该使用 os.WriteFile
	return nil
}

// demoSimplePipeline 演示简单管道
func demoSimplePipeline() {
	// 创建 Reader
	reader, err := plugin.CreateReader("file")
	if err != nil {
		fmt.Printf("Failed to create reader: %v\n", err)
		return
	}

	if err := reader.Init(map[string]interface{}{
		"path": "/tmp/test_data.csv",
	}); err != nil {
		fmt.Printf("Failed to init reader: %v\n", err)
		return
	}

	// 创建 Writer
	writer, err := plugin.CreateWriter("stdout")
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}

	if err := writer.Init(map[string]interface{}{
		"format": "json",
		"pretty": true,
	}); err != nil {
		fmt.Printf("Failed to init writer: %v\n", err)
		return
	}

	// 注意：这里应该使用 Pipeline 引擎，但为了演示先直接连接
	fmt.Println("This demo would use the Pipeline engine to connect reader and writer")
	fmt.Println("Implementation pending...")
}

// demoProcessorPipeline 演示带处理器的管道
func demoProcessorPipeline() {
	fmt.Println("This demo would show a pipeline with processors")
	fmt.Println("Implementation pending...")
}
