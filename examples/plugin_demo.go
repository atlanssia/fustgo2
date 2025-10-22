package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fustgo/fustgo2/pkg/plugin"
	_ "github.com/fustgo/fustgo2/pkg/plugin/processor" // 注册 processors
	_ "github.com/fustgo/fustgo2/pkg/plugin/reader"    // 注册 readers
	_ "github.com/fustgo/fustgo2/pkg/plugin/writer"    // 注册 writers
	"github.com/fustgo/fustgo2/pkg/record"
)

func main() {
	fmt.Println("=== FustGo Plugin System Demo ===\n")

	// 1. 列出所有已注册的插件
	listPlugins()

	// 2. 演示 FileReader + StdoutWriter
	fmt.Println("\n=== Demo 1: FileReader + StdoutWriter ===")
	demoFileReader()

	// 3. 演示 FieldMapper
	fmt.Println("\n=== Demo 2: FieldMapper Processor ===")
	demoFieldMapper()

	// 4. 演示 Filter
	fmt.Println("\n=== Demo 3: Filter Processor ===")
	demoFilter()

	// 5. 演示完整管道: Reader -> Processor -> Writer
	fmt.Println("\n=== Demo 4: Complete Pipeline ===")
	demoPipeline()
}

// listPlugins 列出所有注册的插件
func listPlugins() {
	fmt.Println("Registered Readers:")
	for _, name := range plugin.ListReaders() {
		if info, err := plugin.GetInfo(name); err == nil && info != nil {
			fmt.Printf("  - %s (v%s): %s\n", info.Name, info.Version, info.Description)
		}
	}

	fmt.Println("\nRegistered Writers:")
	for _, name := range plugin.ListWriters() {
		if info, err := plugin.GetInfo(name); err == nil && info != nil {
			fmt.Printf("  - %s (v%s): %s\n", info.Name, info.Version, info.Description)
		}
	}

	fmt.Println("\nRegistered Processors:")
	for _, name := range plugin.ListProcessors() {
		if info, err := plugin.GetInfo(name); err == nil && info != nil {
			fmt.Printf("  - %s (v%s): %s\n", info.Name, info.Version, info.Description)
		}
	}
}

// demoFileReader 演示 FileReader
func demoFileReader() {
	// 创建测试 CSV 文件
	testFile := "/tmp/test_demo.csv"
	content := `id,name,age
1,Alice,30
2,Bob,25
3,Charlie,35
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// 创建 FileReader
	reader, err := plugin.CreateReader("file")
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	// 初始化配置
	config := map[string]interface{}{
		"path":   testFile,
		"format": "csv",
	}

	if err := reader.Init(config); err != nil {
		log.Fatalf("Failed to init reader: %v", err)
	}

	// 创建输出通道
	output := make(chan *record.Record, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动读取
	done := make(chan error, 1)
	go func() {
		done <- reader.Read(ctx, output)
	}()

	// 读取并打印记录
	recordCount := 0
	for rec := range output {
		recordCount++
		fmt.Printf("  Record %d: %+v\n", recordCount, rec.Data)
	}

	// 等待完成
	if err := <-done; err != nil {
		log.Printf("Read error: %v", err)
	}

	fmt.Printf("Total records read: %d\n", recordCount)
}

// demoFieldMapper 演示 FieldMapper
func demoFieldMapper() {
	// 创建 FieldMapper
	mapper, err := plugin.CreateProcessor("field_mapper")
	if err != nil {
		log.Fatalf("Failed to create mapper: %v", err)
	}

	// 配置字段映射
	config := map[string]interface{}{
		"mappings": map[string]interface{}{
			"user_id":    "id",
			"user_name":  "name",
			"user_email": "email",
		},
		"drop_unmapped": false,
	}

	if err := mapper.Init(config); err != nil {
		log.Fatalf("Failed to init mapper: %v", err)
	}

	// 创建测试数据
	input := make(chan *record.Record, 1)
	output := make(chan *record.Record, 1)

	testData := record.NewRecord(map[string]interface{}{
		"user_id":    123,
		"user_name":  "John Doe",
		"user_email": "john@example.com",
		"extra":      "keep this",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动处理
	go func() {
		mapper.(interface {
			Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
		}).Process(ctx, input, output)
	}()

	// 发送数据
	input <- testData
	close(input)

	// 接收结果
	result := <-output
	fmt.Printf("  Original: %+v\n", testData.Data)
	fmt.Printf("  Mapped:   %+v\n", result.Data)
}

// demoFilter 演示 Filter
func demoFilter() {
	// 创建 Filter
	filter, err := plugin.CreateProcessor("filter")
	if err != nil {
		log.Fatalf("Failed to create filter: %v", err)
	}

	// 配置过滤规则: age >= 30
	config := map[string]interface{}{
		"mode": "include",
		"rules": []interface{}{
			map[string]interface{}{
				"field":    "age",
				"operator": "gte",
				"value":    30,
			},
		},
	}

	if err := filter.Init(config); err != nil {
		log.Fatalf("Failed to init filter: %v", err)
	}

	// 创建测试数据
	input := make(chan *record.Record, 3)
	output := make(chan *record.Record, 3)

	testRecords := []*record.Record{
		record.NewRecord(map[string]interface{}{"name": "Alice", "age": 30}),
		record.NewRecord(map[string]interface{}{"name": "Bob", "age": 25}),
		record.NewRecord(map[string]interface{}{"name": "Charlie", "age": 35}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动处理
	go func() {
		filter.(interface {
			Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
		}).Process(ctx, input, output)
	}()

	// 发送数据
	for _, rec := range testRecords {
		input <- rec
		fmt.Printf("  Input: %+v\n", rec.Data)
	}
	close(input)

	// 接收结果
	fmt.Println("\n  Filtered results (age >= 30):")
	for result := range output {
		fmt.Printf("    -> %+v\n", result.Data)
	}

	// 获取指标
	metrics := filter.(interface {
		GetMetrics() *plugin.Metrics
	}).GetMetrics()
	fmt.Printf("\n  Metrics: Processed=%d, Filtered=%d\n",
		metrics.RecordsProcessed, metrics.FilteredCount)
}

// demoPipeline 演示完整的数据管道
func demoPipeline() {
	// 创建测试 JSON 文件
	testFile := "/tmp/pipeline_demo.json"
	content := `{"user_id": 1, "user_name": "Alice", "age": 30}
{"user_id": 2, "user_name": "Bob", "age": 25}
{"user_id": 3, "user_name": "Charlie", "age": 35}
{"user_id": 4, "user_name": "David", "age": 40}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	fmt.Printf("Pipeline: FileReader -> FieldMapper -> Filter -> StdoutWriter\n\n")

	// 1. 创建 FileReader
	reader, _ := plugin.CreateReader("file")
	reader.Init(map[string]interface{}{
		"path":   testFile,
		"format": "json",
	})

	// 2. 创建 FieldMapper
	mapper, _ := plugin.CreateProcessor("field_mapper")
	mapper.Init(map[string]interface{}{
		"mappings": map[string]interface{}{
			"user_id":   "id",
			"user_name": "name",
		},
		"drop_unmapped": false,
	})

	// 3. 创建 Filter (age >= 30)
	filter, _ := plugin.CreateProcessor("filter")
	filter.Init(map[string]interface{}{
		"mode": "include",
		"rules": []interface{}{
			map[string]interface{}{
				"field":    "age",
				"operator": "gte",
				"value":    30,
			},
		},
	})

	// 4. 创建 StdoutWriter
	writer, _ := plugin.CreateWriter("stdout")
	writer.Init(map[string]interface{}{
		"format": "json",
		"pretty": true,
	})

	// 创建通道连接各个组件
	readerOutput := make(chan *record.Record, 10)
	mapperOutput := make(chan *record.Record, 10)
	filterOutput := make(chan *record.Record, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 启动管道
	go reader.Read(ctx, readerOutput)
	go mapper.(interface {
		Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
	}).Process(ctx, readerOutput, mapperOutput)
	go filter.(interface {
		Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
	}).Process(ctx, mapperOutput, filterOutput)

	// 写入输出
	if err := writer.Write(ctx, filterOutput); err != nil {
		log.Printf("Write error: %v", err)
	}

	// 等待片刻确保所有输出完成
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n✅ Pipeline completed!")

	// 打印最终指标
	fmt.Println("\nPipeline Metrics:")
	
	readerProgress := reader.GetProgress()
	fmt.Printf("  Reader: %d records read\n", readerProgress.ProcessedRecords)

	mapperMetrics := mapper.(interface {
		GetMetrics() *plugin.Metrics
	}).GetMetrics()
	fmt.Printf("  Mapper: %d records processed\n", mapperMetrics.RecordsProcessed)

	filterMetrics := filter.(interface {
		GetMetrics() *plugin.Metrics
	}).GetMetrics()
	fmt.Printf("  Filter: %d records processed, %d filtered out\n",
		filterMetrics.RecordsProcessed, filterMetrics.FilteredCount)

	writerMetrics := writer.GetMetrics()
	fmt.Printf("  Writer: %d records written\n", writerMetrics.RecordsProcessed)
}

func init() {
	// 确保示例目录存在
	exampleDir := filepath.Dir("/tmp/test_demo.csv")
	os.MkdirAll(exampleDir, 0755)
}
