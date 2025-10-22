package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// mockReader 模拟 Reader
type mockReader struct {
	*plugin.BasePlugin
	metrics  *plugin.BaseMetrics
	progress *plugin.BaseProgress
	records  []*record.Record
}

func newMockReader(records []*record.Record) *mockReader {
	return &mockReader{
		BasePlugin: plugin.NewBasePlugin("mock_reader", plugin.TypeReader),
		metrics:    plugin.NewBaseMetrics(),
		progress:   plugin.NewBaseProgress(),
		records:    records,
	}
}

func (r *mockReader) Read(ctx context.Context, output chan<- *record.Record) error {
	defer close(output)
	
	for _, rec := range r.records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- rec:
			r.metrics.IncrRecordsProcessed(1)
			r.progress.IncrProcessedRecords(1)
		}
	}
	
	return nil
}

func (r *mockReader) GetProgress() *plugin.Progress {
	return r.progress.GetProgress()
}

func (r *mockReader) Close() error {
	return nil
}

// mockProcessor 模拟 Processor
type mockProcessor struct {
	*plugin.BasePlugin
	metrics *plugin.BaseMetrics
	name    string
	processFunc func(*record.Record) *record.Record
}

func newMockProcessor(name string, processFunc func(*record.Record) *record.Record) *mockProcessor {
	return &mockProcessor{
		BasePlugin:  plugin.NewBasePlugin(name, plugin.TypeProcessor),
		metrics:     plugin.NewBaseMetrics(),
		name:        name,
		processFunc: processFunc,
	}
}

func (p *mockProcessor) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
	for rec := range input {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			processed := rec
			if p.processFunc != nil {
				processed = p.processFunc(rec)
			}
			
			// 如果处理函数返回 nil，表示过滤掉该记录
			if processed == nil {
				p.metrics.IncrFilteredCount(1)
				continue
			}
			
			select {
			case output <- processed:
				p.metrics.IncrRecordsProcessed(1)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	
	// 处理器应该关闭输出通道
	close(output)
	return nil
}

func (p *mockProcessor) GetMetrics() *plugin.Metrics {
	return p.metrics.GetMetrics()
}

func (p *mockProcessor) Close() error {
	return nil
}

// mockFilterProcessor 模拟过滤处理器
type mockFilterProcessor struct {
	*mockProcessor
}

func (f *mockFilterProcessor) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
	for rec := range input {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 过滤掉 id=2 的记录
			if id, ok := rec.Data["id"].(int); ok && id == 2 {
				f.metrics.IncrFilteredCount(1)
				continue
			}
			
			select {
			case output <- rec:
				f.metrics.IncrRecordsProcessed(1)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	
	// 处理器应该关闭输出通道
	close(output)
	return nil
}

// mockWriter 模拟 Writer
type mockWriter struct {
	*plugin.BasePlugin
	metrics *plugin.BaseMetrics
	records []*record.Record
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		BasePlugin: plugin.NewBasePlugin("mock_writer", plugin.TypeWriter),
		metrics:    plugin.NewBaseMetrics(),
		records:    make([]*record.Record, 0),
	}
}

func (w *mockWriter) Write(ctx context.Context, input <-chan *record.Record) error {
	for rec := range input {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			w.records = append(w.records, rec)
			w.metrics.IncrRecordsProcessed(1)
		}
	}
	
	return nil
}

func (w *mockWriter) Flush() error {
	return nil
}

func (w *mockWriter) GetMetrics() *plugin.Metrics {
	return w.metrics.GetMetrics()
}

func (w *mockWriter) Close() error {
	return nil
}

// TestPipelineBuilder 测试管道构建器
func TestPipelineBuilder(t *testing.T) {
	// 创建模拟插件
	reader := newMockReader([]*record.Record{
		record.NewRecord(map[string]interface{}{"id": 1, "name": "Alice"}),
	})
	
	processor := newMockProcessor("mock_processor", func(rec *record.Record) *record.Record {
		// 简单处理：添加 processed 字段
		data := make(map[string]interface{})
		for k, v := range rec.Data {
			data[k] = v
		}
		data["processed"] = true
		return record.NewRecord(data)
	})
	
	writer := newMockWriter()
	
	// 测试正常构建
	pipeline, err := NewBuilder("test-pipeline", "Test Pipeline").
		SetReader(reader).
		AddProcessor(processor).
		SetWriter(writer).
		WithDescription("A test pipeline").
		Build()
	
	if err != nil {
		t.Fatalf("Failed to build pipeline: %v", err)
	}
	
	if pipeline == nil {
		t.Fatal("Pipeline should not be nil")
	}
	
	if pipeline.ID != "test-pipeline" {
		t.Errorf("Expected ID 'test-pipeline', got '%s'", pipeline.ID)
	}
	
	if pipeline.Name != "Test Pipeline" {
		t.Errorf("Expected name 'Test Pipeline', got '%s'", pipeline.Name)
	}
	
	if pipeline.Description != "A test pipeline" {
		t.Errorf("Expected description 'A test pipeline', got '%s'", pipeline.Description)
	}
	
	// 测试缺少 Reader
	_, err = NewBuilder("test-pipeline", "Test Pipeline").
		SetWriter(writer).
		Build()
	
	if err == nil {
		t.Error("Expected error for missing reader, got nil")
	}
	
	// 测试缺少 Writer
	_, err = NewBuilder("test-pipeline", "Test Pipeline").
		SetReader(reader).
		Build()
	
	if err == nil {
		t.Error("Expected error for missing writer, got nil")
	}
}

// TestPipelineExecution 测试管道执行
func TestPipelineExecution(t *testing.T) {
	// 创建测试数据
	testRecords := []*record.Record{
		record.NewRecord(map[string]interface{}{"id": 1, "name": "Alice"}),
		record.NewRecord(map[string]interface{}{"id": 2, "name": "Bob"}),
		record.NewRecord(map[string]interface{}{"id": 3, "name": "Charlie"}),
	}
	
	// 创建模拟插件
	reader := newMockReader(testRecords)
	
	// 创建处理器：添加 processed 字段
	processor := newMockProcessor("processor", func(rec *record.Record) *record.Record {
		data := make(map[string]interface{})
		for k, v := range rec.Data {
			data[k] = v
		}
		data["processed"] = true
		return record.NewRecord(data)
	})
	
	writer := newMockWriter()
	
	// 构建管道
	pipeline, err := NewBuilder("test-execution", "Test Execution").
		SetReader(reader).
		AddProcessor(processor).
		SetWriter(writer).
		WithConfig(&Config{
			Parallelism:   2,
			ChannelSize:   10,
			ErrorStrategy: StrategyContinue,
		}).
		Build()
	
	if err != nil {
		t.Fatalf("Failed to build pipeline: %v", err)
	}
	
	// 创建执行器
	executor := NewExecutor(pipeline)
	
	// 启动执行
	if err := executor.Start(); err != nil {
		t.Fatalf("Failed to start pipeline: %v", err)
	}
	
	// 等待执行完成
	time.Sleep(100 * time.Millisecond) // 给 goroutine 一些时间来执行
	if err := executor.Wait(); err != nil {
		t.Fatalf("Failed to wait for pipeline: %v", err)
	}
	
	// 验证结果
	if len(writer.records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(writer.records))
		for i, rec := range writer.records {
			t.Logf("Record %d: %+v", i, rec.Data)
		}
	}
	
	// 验证处理后的数据
	for _, rec := range writer.records {
		if _, ok := rec.Data["processed"]; !ok {
			t.Error("Expected 'processed' field in record")
		}
	}
	
	// 验证指标
	metrics := pipeline.GetMetrics()
	if metrics.RecordsRead != 3 {
		t.Errorf("Expected 3 records read, got %d", metrics.RecordsRead)
	}
	
	if metrics.RecordsProcessed != 3 {
		t.Errorf("Expected 3 records processed, got %d", metrics.RecordsProcessed)
	}
	
	if metrics.RecordsWritten != 3 {
		t.Errorf("Expected 3 records written, got %d", metrics.RecordsWritten)
	}
}

// TestPipelinePauseResume 测试管道暂停和恢复
func TestPipelinePauseResume(t *testing.T) {
	// 创建测试数据
	testRecords := []*record.Record{
		record.NewRecord(map[string]interface{}{"id": 1, "name": "Alice"}),
		record.NewRecord(map[string]interface{}{"id": 2, "name": "Bob"}),
	}
	
	// 创建模拟插件
	reader := newMockReader(testRecords)
	writer := newMockWriter()
	
	// 构建管道
	pipeline, err := NewBuilder("test-pause", "Test Pause").
		SetReader(reader).
		SetWriter(writer).
		Build()
	
	if err != nil {
		t.Fatalf("Failed to build pipeline: %v", err)
	}
	
	// 创建执行器
	executor := NewExecutor(pipeline)
	
	// 测试未启动时的操作
	if err := executor.Pause(); err == nil {
		t.Error("Expected error when pausing unstarted pipeline")
	}
	
	if err := executor.Resume(); err == nil {
		t.Error("Expected error when resuming unstarted pipeline")
	}
	
	// 启动执行
	if err := executor.Start(); err != nil {
		t.Fatalf("Failed to start pipeline: %v", err)
	}
	
	// 测试暂停
	if err := executor.Pause(); err != nil {
		t.Errorf("Failed to pause pipeline: %v", err)
	}
	
	// 测试重复暂停
	if err := executor.Pause(); err == nil {
		t.Error("Expected error when pausing already paused pipeline")
	}
	
	// 测试恢复
	if err := executor.Resume(); err != nil {
		t.Errorf("Failed to resume pipeline: %v", err)
	}
	
	// 测试重复恢复
	if err := executor.Resume(); err == nil {
		t.Error("Expected error when resuming already running pipeline")
	}
	
	// 等待执行完成
	if err := executor.Wait(); err != nil {
		t.Fatalf("Failed to wait for pipeline: %v", err)
	}
}

// TestPipelineStop 测试管道停止
func TestPipelineStop(t *testing.T) {
	// 创建测试数据
	testRecords := []*record.Record{
		record.NewRecord(map[string]interface{}{"id": 1, "name": "Alice"}),
		record.NewRecord(map[string]interface{}{"id": 2, "name": "Bob"}),
	}
	
	// 创建模拟插件
	reader := newMockReader(testRecords)
	writer := newMockWriter()
	
	// 构建管道
	pipeline, err := NewBuilder("test-stop", "Test Stop").
		SetReader(reader).
		SetWriter(writer).
		Build()
	
	if err != nil {
		t.Fatalf("Failed to build pipeline: %v", err)
	}
	
	// 创建执行器
	executor := NewExecutor(pipeline)
	
	// 测试未启动时停止
	if err := executor.Stop(); err == nil {
		t.Error("Expected error when stopping unstarted pipeline")
	}
	
	// 启动执行
	if err := executor.Start(); err != nil {
		t.Fatalf("Failed to start pipeline: %v", err)
	}
	
	// 等待一小段时间确保开始执行
	time.Sleep(10 * time.Millisecond)
	
	// 停止管道
	if err := executor.Stop(); err != nil {
		t.Errorf("Failed to stop pipeline: %v", err)
	}
	
	// 测试重复停止
	if err := executor.Stop(); err == nil {
		t.Error("Expected error when stopping already stopped pipeline")
	}
}

// TestPipelineErrorHandling 测试错误处理
func TestPipelineErrorHandling(t *testing.T) {
	// 测试不同的错误策略
	
	// 1. StrategyFailFast
	t.Run("FailFast", func(t *testing.T) {
		// TODO: 实现失败快速策略测试
	})
	
	// 2. StrategyContinue
	t.Run("Continue", func(t *testing.T) {
		// TODO: 实现继续策略测试
	})
	
	// 3. StrategyRetry
	t.Run("Retry", func(t *testing.T) {
		// TODO: 实现重试策略测试
	})
}

// BenchmarkPipelineThroughput 测试管道吞吐量
func BenchmarkPipelineThroughput(b *testing.B) {
	// 创建大量测试数据
	var testRecords []*record.Record
	for i := 0; i < 10000; i++ {
		testRecords = append(testRecords, record.NewRecord(map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("User%d", i),
		}))
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// 创建模拟插件
		reader := newMockReader(testRecords)
		writer := newMockWriter()
		
		// 构建管道
		pipeline, err := NewBuilder("benchmark", "Benchmark").
			SetReader(reader).
			SetWriter(writer).
			Build()
		
		if err != nil {
			b.Fatalf("Failed to build pipeline: %v", err)
		}
		
		// 创建执行器
		executor := NewExecutor(pipeline)
		
		// 启动执行
		if err := executor.Start(); err != nil {
			b.Fatalf("Failed to start pipeline: %v", err)
		}
		
		// 等待执行完成
		if err := executor.Wait(); err != nil {
			b.Fatalf("Failed to wait for pipeline: %v", err)
		}
	}
}
