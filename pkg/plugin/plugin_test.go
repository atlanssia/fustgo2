package plugin

import (
	"fmt"
	"testing"
)

func TestRegistry(t *testing.T) {
	// 创建新的注册表用于测试
	registry := NewRegistry()

	// 创建测试插件工厂
	testFactory := func() Plugin {
		return NewBasePlugin("test", TypeReader)
	}

	testInfo := &Info{
		Name:        "test",
		Type:        TypeReader,
		Version:     "1.0.0",
		Author:      "Test",
		Description: "Test plugin",
	}

	// 测试注册
	err := registry.Register("test", testFactory, testInfo)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	// 测试重复注册
	err = registry.Register("test", testFactory, testInfo)
	if err == nil {
		t.Error("Expected error when registering duplicate plugin")
	}

	// 测试获取插件信息
	info, err := registry.GetInfo("test")
	if err != nil {
		t.Fatalf("Failed to get plugin info: %v", err)
	}
	if info.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", info.Name)
	}

	// 测试列出 Reader
	readers := registry.ListReaders()
	if len(readers) != 1 {
		t.Errorf("Expected 1 reader, got %d", len(readers))
	}
	if readers[0] != "test" {
		t.Errorf("Expected reader 'test', got '%s'", readers[0])
	}

	// 测试列出所有插件
	allInfos := registry.ListAll()
	if len(allInfos) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(allInfos))
	}
}

func TestRegistryWrongType(t *testing.T) {
	registry := NewRegistry()

	// 注册一个 Reader
	testFactory := func() Plugin {
		return NewBasePlugin("test", TypeReader)
	}

	testInfo := &Info{
		Name:    "test",
		Type:    TypeReader,
		Version: "1.0.0",
	}

	err := registry.Register("test", testFactory, testInfo)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	// 尝试作为 Writer 创建应该失败
	_, err = registry.CreateWriter("test")
	if err == nil {
		t.Error("Expected error when creating plugin with wrong type")
	}
}

func TestBasePlugin(t *testing.T) {
	plugin := NewBasePlugin("test", TypeReader)

	// 测试基本属性
	if plugin.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", plugin.Name())
	}

	if plugin.Type() != TypeReader {
		t.Errorf("Expected type Reader, got %s", plugin.Type())
	}

	// 测试配置
	config := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	err := plugin.Init(config)
	if err != nil {
		t.Fatalf("Failed to init plugin: %v", err)
	}

	// 测试获取配置
	if str, ok := plugin.GetConfigString("key1"); !ok || str != "value1" {
		t.Error("Failed to get string config")
	}

	if num, ok := plugin.GetConfigInt("key2"); !ok || num != 123 {
		t.Error("Failed to get int config")
	}

	if b, ok := plugin.GetConfigBool("key3"); !ok || !b {
		t.Error("Failed to get bool config")
	}

	// 测试获取不存在的配置
	if _, ok := plugin.GetConfigString("nonexistent"); ok {
		t.Error("Should not find nonexistent config")
	}

	// 测试关闭
	err = plugin.Close()
	if err != nil {
		t.Errorf("Failed to close plugin: %v", err)
	}
}

func TestBaseMetrics(t *testing.T) {
	metrics := NewBaseMetrics()

	// 初始值应该为 0
	m := metrics.GetMetrics()
	if m.RecordsProcessed != 0 {
		t.Error("Initial RecordsProcessed should be 0")
	}

	// 增加计数
	metrics.IncrRecordsProcessed(10)
	metrics.IncrBytesProcessed(1024)
	metrics.IncrSuccessCount(9)
	metrics.IncrErrorCount(1)

	// 验证
	m = metrics.GetMetrics()
	if m.RecordsProcessed != 10 {
		t.Errorf("Expected RecordsProcessed 10, got %d", m.RecordsProcessed)
	}
	if m.BytesProcessed != 1024 {
		t.Errorf("Expected BytesProcessed 1024, got %d", m.BytesProcessed)
	}
	if m.SuccessCount != 9 {
		t.Errorf("Expected SuccessCount 9, got %d", m.SuccessCount)
	}
	if m.ErrorCount != 1 {
		t.Errorf("Expected ErrorCount 1, got %d", m.ErrorCount)
	}

	// 测试重置
	metrics.Reset()
	m = metrics.GetMetrics()
	if m.RecordsProcessed != 0 {
		t.Error("RecordsProcessed should be 0 after reset")
	}
}

func TestBaseProgress(t *testing.T) {
	progress := NewBaseProgress()

	// 设置总记录数
	progress.SetTotalRecords(100)

	// 增加已处理记录数
	progress.IncrProcessedRecords(30)
	progress.IncrBytesProcessed(3000)

	// 设置消息
	progress.SetMessage("Processing...")

	// 获取进度
	p := progress.GetProgress()

	if p.TotalRecords != 100 {
		t.Errorf("Expected TotalRecords 100, got %d", p.TotalRecords)
	}

	if p.ProcessedRecords != 30 {
		t.Errorf("Expected ProcessedRecords 30, got %d", p.ProcessedRecords)
	}

	if p.Percentage != 30.0 {
		t.Errorf("Expected Percentage 30.0, got %f", p.Percentage)
	}

	if p.Message != "Processing..." {
		t.Errorf("Expected Message 'Processing...', got '%s'", p.Message)
	}

	// 测试百分比计算
	progress.IncrProcessedRecords(70)
	p = progress.GetProgress()
	if p.Percentage != 100.0 {
		t.Errorf("Expected Percentage 100.0, got %f", p.Percentage)
	}
}

func TestPluginError(t *testing.T) {
	err := NewError("test_plugin", "test_operation", nil, true)

	if err.Plugin != "test_plugin" {
		t.Errorf("Expected plugin 'test_plugin', got '%s'", err.Plugin)
	}

	if err.Operation != "test_operation" {
		t.Errorf("Expected operation 'test_operation', got '%s'", err.Operation)
	}

	if !err.Retryable {
		t.Error("Expected error to be retryable")
	}

	// 测试 IsRetryable
	if !IsRetryable(err) {
		t.Error("IsRetryable should return true")
	}

	// 测试非插件错误
	regularErr := fmt.Errorf("regular error")
	if IsRetryable(regularErr) {
		t.Error("IsRetryable should return false for regular error")
	}
}

func BenchmarkMetricsIncrement(b *testing.B) {
	metrics := NewBaseMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.IncrRecordsProcessed(1)
	}
}

func BenchmarkProgressUpdate(b *testing.B) {
	progress := NewBaseProgress()
	progress.SetTotalRecords(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		progress.IncrProcessedRecords(1)
	}
}
