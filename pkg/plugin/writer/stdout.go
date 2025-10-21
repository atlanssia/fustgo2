package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// StdoutWriter 标准输出写入器
// 将数据输出到标准输出（控制台）
type StdoutWriter struct {
	*plugin.BasePlugin
	metrics *plugin.BaseMetrics

	// 配置
	format string // json, text, table
	pretty bool   // 是否格式化 JSON
	mu     sync.Mutex
}

// NewStdoutWriter 创建标准输出写入器
func NewStdoutWriter() plugin.Plugin {
	return &StdoutWriter{
		BasePlugin: plugin.NewBasePlugin("stdout", plugin.TypeWriter),
		metrics:    plugin.NewBaseMetrics(),
		format:     "json",
		pretty:     false,
	}
}

// Init 初始化插件
func (w *StdoutWriter) Init(config map[string]interface{}) error {
	if err := w.BasePlugin.Init(config); err != nil {
		return err
	}

	// 获取配置
	if format, ok := w.GetConfigString("format"); ok {
		w.format = format
	}

	if pretty, ok := w.GetConfigBool("pretty"); ok {
		w.pretty = pretty
	}

	// 验证配置
	return w.Validate()
}

// Validate 验证配置
func (w *StdoutWriter) Validate() error {
	// 验证格式
	validFormats := []string{"json", "text", "table"}
	valid := false
	for _, f := range validFormats {
		if w.format == f {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unsupported format: %s (valid: json, text, table)", w.format)
	}

	return nil
}

// Write 写入数据
func (w *StdoutWriter) Write(ctx context.Context, input <-chan *record.Record) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rec, ok := <-input:
			if !ok {
				// 通道已关闭
				return nil
			}

			if err := w.writeRecord(rec); err != nil {
				w.metrics.IncrErrorCount(1)
				// 继续处理，不中断
				fmt.Fprintf(os.Stderr, "Error writing record: %v\n", err)
				continue
			}

			w.metrics.IncrRecordsProcessed(1)
			w.metrics.IncrSuccessCount(1)
		}
	}
}

// writeRecord 写入单条记录
func (w *StdoutWriter) writeRecord(rec *record.Record) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	switch w.format {
	case "json":
		return w.writeJSON(rec)
	case "text":
		return w.writeText(rec)
	case "table":
		return w.writeTable(rec)
	default:
		return fmt.Errorf("unsupported format: %s", w.format)
	}
}

// writeJSON 以 JSON 格式输出
func (w *StdoutWriter) writeJSON(rec *record.Record) error {
	var data []byte
	var err error

	if w.pretty {
		data, err = json.MarshalIndent(rec, "", "  ")
	} else {
		data, err = json.Marshal(rec)
	}

	if err != nil {
		return plugin.NewError(w.Name(), "marshal_json", err, false)
	}

	fmt.Println(string(data))
	return nil
}

// writeText 以文本格式输出
func (w *StdoutWriter) writeText(rec *record.Record) error {
	// 简单的 key=value 格式
	parts := make([]string, 0, len(rec.Data))
	for key, value := range rec.Data {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}
	fmt.Println(strings.Join(parts, " "))
	return nil
}

// writeTable 以表格格式输出（简化版）
func (w *StdoutWriter) writeTable(rec *record.Record) error {
	// 打印分隔线
	fmt.Println("---")
	
	// 打印每个字段
	for key, value := range rec.Data {
		fmt.Printf("%-20s: %v\n", key, value)
	}
	
	return nil
}

// Flush 刷新缓冲区
func (w *StdoutWriter) Flush() error {
	// 标准输出无需手动刷新
	return nil
}

// GetMetrics 获取指标
func (w *StdoutWriter) GetMetrics() *plugin.Metrics {
	return w.metrics.GetMetrics()
}

// Close 关闭插件
func (w *StdoutWriter) Close() error {
	return nil
}

// 注册插件
func init() {
	info := &plugin.Info{
		Name:        "stdout",
		Type:        plugin.TypeWriter,
		Version:     "1.0.0",
		Author:      "FustGo Team",
		Description: "Write data to standard output (console)",
		ConfigSchema: map[string]interface{}{
			"format": "string (optional) - Output format: json, text, table. Default: json",
			"pretty": "bool (optional) - Pretty print JSON. Default: false",
		},
	}

	plugin.Register("stdout", NewStdoutWriter, info)
}
