package writer

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

func TestStdoutWriterJSON(t *testing.T) {
	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	writer := NewStdoutWriter()

	config := map[string]interface{}{
		"format": "json",
	}

	if err := writer.Init(config); err != nil {
		t.Fatalf("Failed to init writer: %v", err)
	}

	input := make(chan *record.Record)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动写入协程
	done := make(chan error, 1)
	go func() {
		err := writer.(interface {
			Write(context.Context, <-chan *record.Record) error
		}).Write(ctx, input)
		done <- err
	}()

	// 发送测试数据
	rec := record.NewRecord(map[string]interface{}{
		"id":   1,
		"name": "test",
	})
	input <- rec
	close(input)

	// 等待处理完成
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Write timeout")
	}

	// 关闭写入端并读取输出
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}

	// 验证指标
	wr := writer.(interface {
		GetMetrics() *plugin.Metrics
	})
	metrics := wr.GetMetrics()
	if metrics.RecordsProcessed != 1 {
		t.Errorf("Expected 1 record processed, got %d", metrics.RecordsProcessed)
	}
	if metrics.SuccessCount != 1 {
		t.Errorf("Expected 1 success, got %d", metrics.SuccessCount)
	}
}

func TestStdoutWriterText(t *testing.T) {
	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	writer := NewStdoutWriter()

	config := map[string]interface{}{
		"format": "text",
	}

	if err := writer.Init(config); err != nil {
		t.Fatalf("Failed to init writer: %v", err)
	}

	input := make(chan *record.Record)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动写入协程
	done := make(chan error, 1)
	go func() {
		err := writer.(interface {
			Write(context.Context, <-chan *record.Record) error
		}).Write(ctx, input)
		done <- err
	}()

	// 发送测试数据
	rec := record.NewRecord(map[string]interface{}{
		"id":   1,
		"name": "test",
	})
	input <- rec
	close(input)

	// 等待处理完成
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Write timeout")
	}

	// 关闭写入端并读取输出
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}
}

func TestStdoutWriterTable(t *testing.T) {
	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	writer := NewStdoutWriter()

	config := map[string]interface{}{
		"format": "table",
	}

	if err := writer.Init(config); err != nil {
		t.Fatalf("Failed to init writer: %v", err)
	}

	input := make(chan *record.Record)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动写入协程
	done := make(chan error, 1)
	go func() {
		err := writer.(interface {
			Write(context.Context, <-chan *record.Record) error
		}).Write(ctx, input)
		done <- err
	}()

	// 发送多条测试数据
	records := []*record.Record{
		record.NewRecord(map[string]interface{}{"id": 1, "name": "Alice"}),
		record.NewRecord(map[string]interface{}{"id": 2, "name": "Bob"}),
		record.NewRecord(map[string]interface{}{"id": 3, "name": "Charlie"}),
	}

	for _, rec := range records {
		input <- rec
	}
	close(input)

	// 等待处理完成
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Write timeout")
	}

	// 刷新输出
	wr := writer.(interface {
		Flush() error
	})
	if err := wr.Flush(); err != nil {
		t.Errorf("Flush error: %v", err)
	}

	// 关闭写入端并读取输出
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}

	// 验证指标
	wm := writer.(interface {
		GetMetrics() *plugin.Metrics
	})
	metrics := wm.GetMetrics()
	if metrics.RecordsProcessed != 3 {
		t.Errorf("Expected 3 records processed, got %d", metrics.RecordsProcessed)
	}
}

func TestStdoutWriterInvalidFormat(t *testing.T) {
	writer := NewStdoutWriter()

	config := map[string]interface{}{
		"format": "invalid",
	}

	err := writer.Init(config)
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}
}

func TestStdoutWriterCancellation(t *testing.T) {
	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	writer := NewStdoutWriter()

	config := map[string]interface{}{
		"format": "json",
	}

	if err := writer.Init(config); err != nil {
		t.Fatalf("Failed to init writer: %v", err)
	}

	input := make(chan *record.Record)

	ctx, cancel := context.WithCancel(context.Background())

	// 启动写入协程
	done := make(chan error, 1)
	go func() {
		err := writer.(interface {
			Write(context.Context, <-chan *record.Record) error
		}).Write(ctx, input)
		done <- err
	}()

	// 发送一条数据后立即取消
	rec := record.NewRecord(map[string]interface{}{"id": 1})
	input <- rec
	cancel()

	// 等待处理完成
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Write timeout")
	}

	// 清理
	w.Close()
	io.Copy(io.Discard, r)
	os.Stdout = oldStdout
}
