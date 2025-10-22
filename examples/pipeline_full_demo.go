package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fustgo/fustgo2/pkg/plugin"
	_ "github.com/fustgo/fustgo2/pkg/plugin/processor"
	_ "github.com/fustgo/fustgo2/pkg/plugin/reader"
	_ "github.com/fustgo/fustgo2/pkg/plugin/writer"
	"github.com/fustgo/fustgo2/pkg/pipeline"
	"github.com/fustgo/fustgo2/pkg/record"
)

func main() {
	fmt.Println("=== FustGo Pipeline Engine Demo ===\n")

	// 创建测试数据文件
	createTestData()

	// 演示 1: 简单的文件读取和输出
	fmt.Println("Demo 1: Simple File Reader -> Stdout Writer")
	demoSimplePipeline()

	// 演示 2: 带处理器的管道
	fmt.Println("\nDemo 2: File Reader -> Field Mapper -> Filter -> Stdout Writer")
	demoProcessorPipeline()

	// 演示 3: 管道控制（暂停/恢复）
	fmt.Println("\nDemo 3: Pipeline Control (Pause/Resume)")
	demoPipelineControl()
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

	if err := os.WriteFile("/tmp/test_data.csv", []byte(csvContent), 0644); err != nil {
		log.Printf("Failed to create test data: %v", err)
		return
	}

	// 创建 JSON 测试文件
	jsonContent := `{"id": 1, "name": "Alice", "age": 30, "email": "alice@example.com"}
{"id": 2, "name": "Bob", "age": 25, "email": "bob@example.com"}
{"id": 3, "name": "Charlie", "age": 35, "email": "charlie@example.com"}`

	if err := os.WriteFile("/tmp/test_data.json", []byte(jsonContent), 0644); err != nil {
		log.Printf("Failed to create test data: %v", err)
		return
	}

	fmt.Println("Created test data files:")
	fmt.Println("  - /tmp/test_data.csv")
	fmt.Println("  - /tmp/test_data.json")
	fmt.Println()
}

// demoSimplePipeline 演示简单管道
func demoSimplePipeline() {
	// 创建 Reader
	reader, err := plugin.CreateReader("file")
	if err != nil {
		log.Printf("Failed to create reader: %v", err)
		return
	}

	if err := reader.Init(map[string]interface{}{
		"path": "/tmp/test_data.csv",
	}); err != nil {
		log.Printf("Failed to init reader: %v", err)
		return
	}

	// 创建 Writer
	writer, err := plugin.CreateWriter("stdout")
	if err != nil {
		log.Printf("Failed to create writer: %v", err)
		return
	}

	if err := writer.Init(map[string]interface{}{
		"format": "json",
		"pretty": true,
	}); err != nil {
		log.Printf("Failed to init writer: %v", err)
		return
	}

	// 构建管道
	p, err := pipeline.NewBuilder("simple-pipeline", "Simple Pipeline").
		SetReader(reader).
		SetWriter(writer).
		WithDescription("Read CSV file and output as JSON").
		WithConfig(&pipeline.Config{
			ChannelSize:  100,
			ErrorStrategy: pipeline.StrategyContinue,
		}).
		Build()

	if err != nil {
		log.Printf("Failed to build pipeline: %v", err)
		return
	}

	// 创建执行器
	executor := pipeline.NewExecutor(p)

	// 启动执行
	fmt.Println("Starting pipeline execution...")
	if err := executor.Start(); err != nil {
		log.Printf("Failed to start pipeline: %v", err)
		return
	}

	// 等待执行完成
	if err := executor.Wait(); err != nil {
		log.Printf("Failed to wait for pipeline: %v", err)
		return
	}

	// 显示指标
	metrics := p.GetMetrics()
	fmt.Printf("Pipeline completed:\n")
	fmt.Printf("  Records Read: %d\n", metrics.RecordsRead)
	fmt.Printf("  Records Written: %d\n", metrics.RecordsWritten)
	fmt.Printf("  Duration: %v\n", metrics.Duration)
	fmt.Printf("  TPS: %.2f\n", metrics.TPS)
}

// demoProcessorPipeline 演示带处理器的管道
func demoProcessorPipeline() {
	// 创建 Reader
	reader, err := plugin.CreateReader("file")
	if err != nil {
		log.Printf("Failed to create reader: %v", err)
		return
	}

	if err := reader.Init(map[string]interface{}{
		"path": "/tmp/test_data.json",
	}); err != nil {
		log.Printf("Failed to init reader: %v", err)
		return
	}

	// 创建 Field Mapper Processor
	mapper, err := plugin.CreateProcessor("field_mapper")
	if err != nil {
		log.Printf("Failed to create mapper: %v", err)
		return
	}

	if err := mapper.Init(map[string]interface{}{
		"mappings": map[string]interface{}{
			"id":    "user_id",
			"name":  "full_name",
			"email": "email_address",
		},
		"drop_unmapped": true,
	}); err != nil {
		log.Printf("Failed to init mapper: %v", err)
		return
	}

	// 创建 Filter Processor
	filter, err := plugin.CreateProcessor("filter")
	if err != nil {
		log.Printf("Failed to create filter: %v", err)
		return
	}

	if err := filter.Init(map[string]interface{}{
		"mode": "include",
		"rules": []interface{}{
			map[string]interface{}{
				"field":    "age",
				"operator": "gte",
				"value":    30,
			},
		},
	}); err != nil {
		log.Printf("Failed to init filter: %v", err)
		return
	}

	// 创建 Writer
	writer, err := plugin.CreateWriter("stdout")
	if err != nil {
		log.Printf("Failed to create writer: %v", err)
		return
	}

	if err := writer.Init(map[string]interface{}{
		"format": "table",
	}); err != nil {
		log.Printf("Failed to init writer: %v", err)
		return
	}

	// 构建管道
	p, err := pipeline.NewBuilder("processor-pipeline", "Processor Pipeline").
		SetReader(reader).
		AddProcessor(mapper).
		AddProcessor(filter).
		SetWriter(writer).
		WithDescription("Read JSON, map fields, filter by age, output as table").
		WithConfig(&pipeline.Config{
			ChannelSize:  100,
			ErrorStrategy: pipeline.StrategyContinue,
		}).
		Build()

	if err != nil {
		log.Printf("Failed to build pipeline: %v", err)
		return
	}

	// 创建执行器
	executor := pipeline.NewExecutor(p)

	// 启动执行
	fmt.Println("Starting pipeline execution...")
	if err := executor.Start(); err != nil {
		log.Printf("Failed to start pipeline: %v", err)
		return
	}

	// 等待执行完成
	if err := executor.Wait(); err != nil {
		log.Printf("Failed to wait for pipeline: %v", err)
		return
	}

	// 显示指标
	metrics := p.GetMetrics()
	fmt.Printf("Pipeline completed:\n")
	fmt.Printf("  Records Read: %d\n", metrics.RecordsRead)
	fmt.Printf("  Records Processed: %d\n", metrics.RecordsProcessed)
	fmt.Printf("  Records Filtered: %d\n", metrics.RecordsFiltered)
	fmt.Printf("  Records Written: %d\n", metrics.RecordsWritten)
	fmt.Printf("  Duration: %v\n", metrics.Duration)
	fmt.Printf("  TPS: %.2f\n", metrics.TPS)
}

// demoPipelineControl 演示管道控制
func demoPipelineControl() {
	// 创建大量测试数据
	var testRecords []*record.Record
	for i := 1; i <= 100; i++ {
		testRecords = append(testRecords, record.NewRecord(map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("User%d", i),
			"age":  20 + (i % 50),
		}))
	}

	// 创建模拟 Reader
	reader := &mockSlowReader{
		records: testRecords,
		delay:   50 * time.Millisecond, // 每条记录延迟50ms
	}

	// 创建 Writer
	writer, err := plugin.CreateWriter("stdout")
	if err != nil {
		log.Printf("Failed to create writer: %v", err)
		return
	}

	if err := writer.Init(map[string]interface{}{
		"format": "text",
	}); err != nil {
		log.Printf("Failed to init writer: %v", err)
		return
	}

	// 构建管道
	p, err := pipeline.NewBuilder("control-pipeline", "Control Pipeline").
		SetReader(reader).
		SetWriter(writer).
		WithDescription("Demonstrate pipeline control").
		WithConfig(&pipeline.Config{
			ChannelSize:  10,
			ErrorStrategy: pipeline.StrategyContinue,
		}).
		Build()

	if err != nil {
		log.Printf("Failed to build pipeline: %v", err)
		return
	}

	// 创建执行器
	executor := pipeline.NewExecutor(p)

	// 启动执行
	fmt.Println("Starting pipeline execution...")
	if err := executor.Start(); err != nil {
		log.Printf("Failed to start pipeline: %v", err)
		return
	}

	// 在后台监控进度
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				metrics := p.GetMetrics()
				fmt.Printf("Progress: %d records processed, TPS: %.2f\n", 
					metrics.RecordsProcessed, metrics.TPS)
			case <-time.After(10 * time.Second):
				return
			}
		}
	}()

	// 演示控制功能
	time.Sleep(2 * time.Second)
	fmt.Println("Pausing pipeline...")
	if err := executor.Pause(); err != nil {
		log.Printf("Failed to pause pipeline: %v", err)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Resuming pipeline...")
	if err := executor.Resume(); err != nil {
		log.Printf("Failed to resume pipeline: %v", err)
	}

	// 等待执行完成
	if err := executor.Wait(); err != nil {
		log.Printf("Failed to wait for pipeline: %v", err)
		return
	}

	fmt.Println("Pipeline execution completed!")
}

// mockSlowReader 模拟慢速 Reader
type mockSlowReader struct {
	*plugin.BasePlugin
	records []*record.Record
	delay   time.Duration
}

func (r *mockSlowReader) Read(ctx context.Context, output chan<- *record.Record) error {
	defer close(output)

	for _, rec := range r.records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.delay):
			// 模拟处理延迟
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- rec:
			// 发送记录
		}
	}

	return nil
}

func (r *mockSlowReader) GetProgress() *plugin.Progress {
	return &plugin.Progress{
		ProcessedRecords: int64(len(r.records)),
	}
}

func (r *mockSlowReader) Close() error {
	return nil
}

func init() {
	// 注册模拟 Reader
	plugin.Register("mock_slow_reader", func() plugin.Plugin {
		return &mockSlowReader{
			BasePlugin: plugin.NewBasePlugin("mock_slow_reader", plugin.TypeReader),
		}
	}, &plugin.Info{
		Name:        "mock_slow_reader",
		Type:        plugin.TypeReader,
		Version:     "1.0.0",
		Description: "Mock slow reader for demo",
	})
}