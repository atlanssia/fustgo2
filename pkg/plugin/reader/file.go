package reader

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// FileReader 文件读取器插件
// 支持 CSV, JSON, 纯文本格式
type FileReader struct {
	*plugin.BasePlugin
	progress *plugin.BaseProgress
	metrics  *plugin.BaseMetrics

	// 配置
	filePath string
	format   string // csv, json, text
	encoding string // utf-8, gbk, etc.

	// CSV 配置
	delimiter rune
	hasHeader bool
	headers   []string
}

// NewFileReader 创建文件读取器
func NewFileReader() plugin.Plugin {
	return &FileReader{
		BasePlugin: plugin.NewBasePlugin("file", plugin.TypeReader),
		progress:   plugin.NewBaseProgress(),
		metrics:    plugin.NewBaseMetrics(),
		delimiter:  ',',
		encoding:   "utf-8",
	}
}

// Init 初始化插件
func (r *FileReader) Init(config map[string]interface{}) error {
	if err := r.BasePlugin.Init(config); err != nil {
		return err
	}

	// 获取配置
	var ok bool
	r.filePath, ok = r.GetConfigString("file_path")
	if !ok {
		r.filePath, ok = r.GetConfigString("path")
	}
	if !ok || r.filePath == "" {
		return fmt.Errorf("file_path is required")
	}

	// 格式（默认根据文件扩展名判断）
	r.format, _ = r.GetConfigString("format")
	if r.format == "" {
		ext := strings.ToLower(filepath.Ext(r.filePath))
		switch ext {
		case ".csv":
			r.format = "csv"
		case ".json", ".jsonl":
			r.format = "json"
		default:
			r.format = "text"
		}
	}

	// 编码
	if encoding, ok := r.GetConfigString("encoding"); ok {
		r.encoding = encoding
	}

	// CSV 特定配置
	if r.format == "csv" {
		if delim, ok := r.GetConfigString("delimiter"); ok && len(delim) > 0 {
			r.delimiter = rune(delim[0])
		}
		// has_header 默认为 true
		if hasHeader, ok := r.GetConfigBool("has_header"); ok {
			r.hasHeader = hasHeader
		} else {
			r.hasHeader = true // 默认 CSV 有表头
		}
	}

	return nil
}

// Validate 验证配置
func (r *FileReader) Validate() error {
	// 检查文件是否存在
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", r.filePath)
	}

	// 验证格式
	validFormats := []string{"csv", "json", "text"}
	valid := false
	for _, f := range validFormats {
		if r.format == f {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unsupported format: %s", r.format)
	}

	return nil
}

// Read 读取文件
func (r *FileReader) Read(ctx context.Context, output chan<- *record.Record) error {
	// 打开文件
	file, err := os.Open(r.filePath)
	if err != nil {
		close(output)
		return plugin.NewError(r.Name(), "open_file", err, false)
	}
	defer file.Close()
	defer close(output)

	// 获取文件大小
	stat, _ := file.Stat()
	fileSize := stat.Size()
	r.progress.SetMessage(fmt.Sprintf("Reading file: %s (%.2f MB)",
		filepath.Base(r.filePath), float64(fileSize)/(1024*1024)))

	// 根据格式读取
	switch r.format {
	case "csv":
		return r.readCSV(ctx, file, output)
	case "json":
		return r.readJSON(ctx, file, output)
	case "text":
		return r.readText(ctx, file, output)
	default:
		return fmt.Errorf("unsupported format: %s", r.format)
	}
}

// readCSV 读取 CSV 文件
func (r *FileReader) readCSV(ctx context.Context, file *os.File, output chan<- *record.Record) error {
	reader := csv.NewReader(file)
	reader.Comma = r.delimiter
	reader.LazyQuotes = true

	lineNum := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.metrics.IncrErrorCount(1)
			continue
		}

		lineNum++

		// 第一行作为 header
		if lineNum == 1 && r.hasHeader {
			r.headers = row
			continue
		}

		// 创建记录
		data := make(map[string]interface{})
		for i, value := range row {
			var key string
			if len(r.headers) > i {
				key = r.headers[i]
			} else {
				key = fmt.Sprintf("col_%d", i)
			}
			data[key] = value
		}

		// 创建 Record
		rec := record.NewRecord(data)
		rec.Meta.Source = r.filePath
		rec.Meta.Table = filepath.Base(r.filePath)
		rec.Meta.Offset = int64(lineNum)

		// 发送记录
		select {
		case output <- rec:
			r.metrics.IncrRecordsProcessed(1)
			r.metrics.IncrSuccessCount(1)
			r.progress.IncrProcessedRecords(1)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// readJSON 读取 JSON 文件（每行一个 JSON 对象）
func (r *FileReader) readJSON(ctx context.Context, file *os.File, output chan<- *record.Record) error {
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 解析 JSON
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			r.metrics.IncrErrorCount(1)
			continue
		}

		// 创建 Record
		rec := record.NewRecord(data)
		rec.Meta.Source = r.filePath
		rec.Meta.Table = filepath.Base(r.filePath)
		rec.Meta.Offset = int64(lineNum)

		// 发送记录
		select {
		case output <- rec:
			r.metrics.IncrRecordsProcessed(1)
			r.metrics.IncrSuccessCount(1)
			r.progress.IncrProcessedRecords(1)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return plugin.NewError(r.Name(), "scan_file", err, false)
	}

	return nil
}

// readText 读取纯文本文件（每行一个记录）
func (r *FileReader) readText(ctx context.Context, file *os.File, output chan<- *record.Record) error {
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lineNum++
		line := scanner.Text()

		// 创建 Record
		data := map[string]interface{}{
			"line":    lineNum,
			"content": line,
		}
		rec := record.NewRecord(data)
		rec.Meta.Source = r.filePath
		rec.Meta.Table = filepath.Base(r.filePath)
		rec.Meta.Offset = int64(lineNum)

		// 发送记录
		select {
		case output <- rec:
			r.metrics.IncrRecordsProcessed(1)
			r.metrics.IncrSuccessCount(1)
			r.progress.IncrProcessedRecords(1)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return plugin.NewError(r.Name(), "scan_file", err, false)
	}

	return nil
}

// GetProgress 获取进度
func (r *FileReader) GetProgress() *plugin.Progress {
	return r.progress.GetProgress()
}

// Close 关闭插件
func (r *FileReader) Close() error {
	return nil
}

// 注册插件
func init() {
	info := &plugin.Info{
		Name:        "file",
		Type:        plugin.TypeReader,
		Version:     "1.0.0",
		Author:      "FustGo Team",
		Description: "Read data from files (CSV, JSON, Text)",
		ConfigSchema: map[string]interface{}{
			"file_path": "string (required) - Path to the file",
			"format":    "string (optional) - File format: csv, json, text",
			"encoding":  "string (optional) - File encoding, default: utf-8",
			"delimiter": "string (optional) - CSV delimiter, default: ,",
			"has_header": "bool (optional) - CSV has header row, default: false",
		},
	}

	plugin.Register("file", NewFileReader, info)
}
