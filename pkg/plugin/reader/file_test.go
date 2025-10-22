package reader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fustgo/fustgo2/pkg/record"
)

func TestFileReaderCSV(t *testing.T) {
	// 创建测试 CSV 文件
	tmpFile, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试数据
	csvData := `name,age,city
Alice,25,Beijing
Bob,30,Shanghai
Charlie,35,Guangzhou
`
	if _, err := tmpFile.WriteString(csvData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// 创建 FileReader
	reader := NewFileReader()

	// 配置
	config := map[string]interface{}{
		"file_path":  tmpFile.Name(),
		"format":     "csv",
		"has_header": true,
	}

	if err := reader.Init(config); err != nil {
		t.Fatalf("Failed to init reader: %v", err)
	}

	if err := reader.Validate(); err != nil {
		t.Fatalf("Failed to validate reader: %v", err)
	}

	// 读取数据
	output := make(chan *record.Record, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := reader.(interface {
			Read(context.Context, chan<- *record.Record) error
		}).Read(ctx, output); err != nil && err != context.Canceled {
			t.Errorf("Read failed: %v", err)
		}
	}()

	// 收集记录
	var records []*record.Record
	for rec := range output {
		records = append(records, rec)
	}

	// 验证
	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// 验证第一条记录
	if len(records) > 0 {
		rec := records[0]
		if name, ok := rec.GetString("name"); !ok || name != "Alice" {
			t.Errorf("Expected name 'Alice', got '%s'", name)
		}
		if age, ok := rec.GetString("age"); !ok || age != "25" {
			t.Errorf("Expected age '25', got '%s'", age)
		}
		if city, ok := rec.GetString("city"); !ok || city != "Beijing" {
			t.Errorf("Expected city 'Beijing', got '%s'", city)
		}
	}
}

func TestFileReaderJSON(t *testing.T) {
	// 创建测试 JSON 文件
	tmpFile, err := os.CreateTemp("", "test-*.jsonl")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试数据
	jsonData := `{"name":"Alice","age":25,"city":"Beijing"}
{"name":"Bob","age":30,"city":"Shanghai"}
{"name":"Charlie","age":35,"city":"Guangzhou"}
`
	if _, err := tmpFile.WriteString(jsonData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// 创建 FileReader
	reader := NewFileReader()

	// 配置
	config := map[string]interface{}{
		"file_path": tmpFile.Name(),
		"format":    "json",
	}

	if err := reader.Init(config); err != nil {
		t.Fatalf("Failed to init reader: %v", err)
	}

	// 读取数据
	output := make(chan *record.Record, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := reader.(interface {
			Read(context.Context, chan<- *record.Record) error
		}).Read(ctx, output); err != nil && err != context.Canceled {
			t.Errorf("Read failed: %v", err)
		}
	}()

	// 收集记录
	var records []*record.Record
	for rec := range output {
		records = append(records, rec)
	}

	// 验证
	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// 验证第一条记录
	if len(records) > 0 {
		rec := records[0]
		if name, ok := rec.GetString("name"); !ok || name != "Alice" {
			t.Errorf("Expected name 'Alice', got '%s'", name)
		}
		// JSON 数字会被解析为 float64
		if age, ok := rec.GetFloat64("age"); !ok || age != 25 {
			t.Errorf("Expected age 25, got %f", age)
		}
	}
}

func TestFileReaderText(t *testing.T) {
	// 创建测试文本文件
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试数据
	textData := `Line 1
Line 2
Line 3
`
	if _, err := tmpFile.WriteString(textData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// 创建 FileReader
	reader := NewFileReader()

	// 配置
	config := map[string]interface{}{
		"file_path": tmpFile.Name(),
		"format":    "text",
	}

	if err := reader.Init(config); err != nil {
		t.Fatalf("Failed to init reader: %v", err)
	}

	// 读取数据
	output := make(chan *record.Record, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := reader.(interface {
			Read(context.Context, chan<- *record.Record) error
		}).Read(ctx, output); err != nil && err != context.Canceled {
			t.Errorf("Read failed: %v", err)
		}
	}()

	// 收集记录
	var records []*record.Record
	for rec := range output {
		records = append(records, rec)
	}

	// 验证
	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// 验证第一条记录
	if len(records) > 0 {
		rec := records[0]
		if content, ok := rec.GetString("content"); !ok || content != "Line 1" {
			t.Errorf("Expected content 'Line 1', got '%s'", content)
		}
		if line, ok := rec.GetInt64("line"); !ok || line != 1 {
			t.Errorf("Expected line 1, got %d", line)
		}
	}
}

func TestFileReaderFormatDetection(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.csv", "csv"},
		{"test.json", "json"},
		{"test.jsonl", "json"},
		{"test.txt", "text"},
		{"test.log", "text"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			// 创建临时文件
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}

			// 创建 reader
			reader := NewFileReader().(*FileReader)
			config := map[string]interface{}{
				"file_path": tmpFile,
			}

			if err := reader.Init(config); err != nil {
				t.Fatalf("Failed to init: %v", err)
			}

			if reader.format != tt.expected {
				t.Errorf("Expected format %s, got %s", tt.expected, reader.format)
			}
		})
	}
}

func TestFileReaderMissingFile(t *testing.T) {
	reader := NewFileReader()

	config := map[string]interface{}{
		"file_path": "/nonexistent/file.csv",
	}

	if err := reader.Init(config); err != nil {
		t.Fatalf("Failed to init: %v", err)
	}

	// 验证应该失败
	if err := reader.Validate(); err == nil {
		t.Error("Expected validation error for missing file")
	}
}
