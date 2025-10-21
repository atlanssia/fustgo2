package processor

import (
	"context"
	"testing"
	"time"

	"github.com/fustgo/fustgo2/pkg/record"
)

func TestFieldMapper(t *testing.T) {
	// 创建 FieldMapper
	mapper := NewFieldMapper()

	// 配置
	config := map[string]interface{}{
		"mappings": map[string]interface{}{
			"old_name": "new_name",
			"user_id":  "id",
		},
		"drop_unmapped": false,
	}

	if err := mapper.Init(config); err != nil {
		t.Fatalf("Failed to init mapper: %v", err)
	}

	if err := mapper.Validate(); err != nil {
		t.Fatalf("Failed to validate mapper: %v", err)
	}

	// 创建输入输出通道
	input := make(chan *record.Record, 10)
	output := make(chan *record.Record, 10)

	// 启动处理
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := mapper.(interface {
			Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
		}).Process(ctx, input, output); err != nil && err != context.Canceled {
			t.Errorf("Process failed: %v", err)
		}
	}()

	// 发送测试数据
	testRec := record.NewRecord(map[string]interface{}{
		"old_name":  "Alice",
		"user_id":   123,
		"other_field": "value",
	})
	input <- testRec
	close(input)

	// 接收结果
	result := <-output

	// 验证映射
	if name, ok := result.GetString("new_name"); !ok || name != "Alice" {
		t.Errorf("Expected new_name 'Alice', got '%s'", name)
	}

	if id, ok := result.GetInt64("id"); !ok || id != 123 {
		t.Errorf("Expected id 123, got %d", id)
	}

	// 验证未映射字段保留
	if other, ok := result.GetString("other_field"); !ok || other != "value" {
		t.Errorf("Expected other_field 'value', got '%s'", other)
	}
}

func TestFieldMapperDropUnmapped(t *testing.T) {
	// 创建 FieldMapper
	mapper := NewFieldMapper()

	// 配置 - 丢弃未映射字段
	config := map[string]interface{}{
		"mappings": map[string]interface{}{
			"name": "username",
		},
		"drop_unmapped": true,
	}

	if err := mapper.Init(config); err != nil {
		t.Fatalf("Failed to init mapper: %v", err)
	}

	// 创建输入输出通道
	input := make(chan *record.Record, 10)
	output := make(chan *record.Record, 10)

	// 启动处理
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := mapper.(interface {
			Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
		}).Process(ctx, input, output); err != nil && err != context.Canceled {
			t.Errorf("Process failed: %v", err)
		}
	}()

	// 发送测试数据
	testRec := record.NewRecord(map[string]interface{}{
		"name":  "Alice",
		"age":   25,
		"city":  "Beijing",
	})
	input <- testRec
	close(input)

	// 接收结果
	result := <-output

	// 验证只有映射的字段
	if username, ok := result.GetString("username"); !ok || username != "Alice" {
		t.Errorf("Expected username 'Alice', got '%s'", username)
	}

	// 未映射字段应该被丢弃
	if _, ok := result.GetInt64("age"); ok {
		t.Error("age field should be dropped")
	}

	if _, ok := result.GetString("city"); ok {
		t.Error("city field should be dropped")
	}
}

func TestFilter(t *testing.T) {
	// 创建 Filter
	filter := NewFilter()

	// 配置 - 包含模式，年龄 >= 18
	config := map[string]interface{}{
		"mode": "include",
		"rules": []interface{}{
			map[string]interface{}{
				"field":    "age",
				"operator": "gte",
				"value":    float64(18),
			},
		},
	}

	if err := filter.Init(config); err != nil {
		t.Fatalf("Failed to init filter: %v", err)
	}

	if err := filter.Validate(); err != nil {
		t.Fatalf("Failed to validate filter: %v", err)
	}

	// 创建输入输出通道
	input := make(chan *record.Record, 10)
	output := make(chan *record.Record, 10)

	// 启动处理
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := filter.(interface {
			Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
		}).Process(ctx, input, output); err != nil && err != context.Canceled {
			t.Errorf("Process failed: %v", err)
		}
	}()

	// 发送测试数据
	input <- record.NewRecord(map[string]interface{}{"name": "Alice", "age": 25})   // 应通过
	input <- record.NewRecord(map[string]interface{}{"name": "Bob", "age": 17})     // 应被过滤
	input <- record.NewRecord(map[string]interface{}{"name": "Charlie", "age": 30}) // 应通过
	close(input)

	// 收集结果
	var results []*record.Record
	for rec := range output {
		results = append(results, rec)
	}

	// 验证
	if len(results) != 2 {
		t.Errorf("Expected 2 records, got %d", len(results))
	}

	// 验证通过的记录
	if len(results) >= 1 {
		if name, _ := results[0].GetString("name"); name != "Alice" {
			t.Errorf("Expected first record name 'Alice', got '%s'", name)
		}
	}
	if len(results) >= 2 {
		if name, _ := results[1].GetString("name"); name != "Charlie" {
			t.Errorf("Expected second record name 'Charlie', got '%s'", name)
		}
	}
}

func TestFilterExcludeMode(t *testing.T) {
	// 创建 Filter
	filter := NewFilter()

	// 配置 - 排除模式，排除年龄 < 18 的记录
	config := map[string]interface{}{
		"mode": "exclude",
		"rules": []interface{}{
			map[string]interface{}{
				"field":    "age",
				"operator": "lt",
				"value":    float64(18),
			},
		},
	}

	if err := filter.Init(config); err != nil {
		t.Fatalf("Failed to init filter: %v", err)
	}

	// 创建输入输出通道
	input := make(chan *record.Record, 10)
	output := make(chan *record.Record, 10)

	// 启动处理
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := filter.(interface {
			Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
		}).Process(ctx, input, output); err != nil && err != context.Canceled {
			t.Errorf("Process failed: %v", err)
		}
	}()

	// 发送测试数据
	input <- record.NewRecord(map[string]interface{}{"name": "Alice", "age": 25})  // 应通过
	input <- record.NewRecord(map[string]interface{}{"name": "Bob", "age": 15})    // 应被排除
	input <- record.NewRecord(map[string]interface{}{"name": "Charlie", "age": 20}) // 应通过
	close(input)

	// 收集结果
	var results []*record.Record
	for rec := range output {
		results = append(results, rec)
	}

	// 验证
	if len(results) != 2 {
		t.Errorf("Expected 2 records, got %d", len(results))
	}
}

func TestFilterStringOperators(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		value    string
		testData map[string]interface{}
		expected bool
	}{
		{
			name:     "contains - match",
			operator: "contains",
			value:    "test",
			testData: map[string]interface{}{"text": "this is a test"},
			expected: true,
		},
		{
			name:     "contains - no match",
			operator: "contains",
			value:    "xyz",
			testData: map[string]interface{}{"text": "this is a test"},
			expected: false,
		},
		{
			name:     "startswith - match",
			operator: "startswith",
			value:    "hello",
			testData: map[string]interface{}{"text": "hello world"},
			expected: true,
		},
		{
			name:     "endswith - match",
			operator: "endswith",
			value:    "world",
			testData: map[string]interface{}{"text": "hello world"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewFilter()

			config := map[string]interface{}{
				"mode": "include",
				"rules": []interface{}{
					map[string]interface{}{
						"field":    "text",
						"operator": tt.operator,
						"value":    tt.value,
					},
				},
			}

			if err := filter.Init(config); err != nil {
				t.Fatalf("Failed to init filter: %v", err)
			}

			input := make(chan *record.Record)
			output := make(chan *record.Record)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// 启动处理协程
			done := make(chan error, 1)
			go func() {
				err := filter.(interface {
					Process(context.Context, <-chan *record.Record, chan<- *record.Record) error
				}).Process(ctx, input, output)
				done <- err
			}()

			// 发送测试数据
			input <- record.NewRecord(tt.testData)
			close(input)

			// 收集所有输出
			var received bool
			for {
				select {
				case _, ok := <-output:
					if !ok {
						// 输出通道关闭，处理完成
						goto checkResult
					}
					received = true
				case <-time.After(200 * time.Millisecond):
					t.Fatal("Timeout waiting for processing to complete")
				}
			}

			checkResult:
			// 验证结果
			if received != tt.expected {
				if tt.expected {
					t.Errorf("Expected record to pass through, but it was filtered")
				} else {
					t.Errorf("Expected record to be filtered out, but it passed through")
				}
			}

			// 等待处理完成
			select {
			case err := <-done:
				if err != nil {
					t.Errorf("Process error: %v", err)
				}
			case <-time.After(200 * time.Millisecond):
				t.Error("Process did not complete")
			}
		})
	}
}
