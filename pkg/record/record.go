package record

import (
	"encoding/json"
	"time"
)

// Operation 数据操作类型
type Operation string

const (
	OperationInsert Operation = "INSERT"
	OperationUpdate Operation = "UPDATE"
	OperationDelete Operation = "DELETE"
	OperationRead   Operation = "READ"
	OperationUnknown Operation = "UNKNOWN"
)

// Record 数据记录 - ETL 管道中传输的最小单元
type Record struct {
	// Data 记录的数据字段 (key-value 对)
	Data map[string]interface{} `json:"data"`
	
	// Meta 记录元数据
	Meta *RecordMeta `json:"meta"`
	
	// Timestamp 记录时间戳
	Timestamp time.Time `json:"timestamp"`
}

// RecordMeta 记录元数据
type RecordMeta struct {
	// Source 数据源标识
	Source string `json:"source,omitempty"`
	
	// Schema 数据库模式/数据库名
	Schema string `json:"schema,omitempty"`
	
	// Table 表名/集合名
	Table string `json:"table,omitempty"`
	
	// Operation 操作类型
	Operation Operation `json:"operation,omitempty"`
	
	// Offset 数据偏移量(用于断点续传)
	Offset int64 `json:"offset,omitempty"`
	
	// PartitionKey 分区键(用于分布式处理)
	PartitionKey string `json:"partition_key,omitempty"`
	
	// Extra 额外的元数据信息
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// NewRecord 创建新的数据记录
func NewRecord(data map[string]interface{}) *Record {
	return &Record{
		Data:      data,
		Meta:      &RecordMeta{},
		Timestamp: time.Now(),
	}
}

// NewRecordWithMeta 创建带元数据的记录
func NewRecordWithMeta(data map[string]interface{}, meta *RecordMeta) *Record {
	return &Record{
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}

// Clone 克隆记录(深拷贝)
func (r *Record) Clone() *Record {
	// 深拷贝 Data
	data := make(map[string]interface{})
	for k, v := range r.Data {
		data[k] = v
	}
	
	// 深拷贝 Meta
	meta := &RecordMeta{
		Source:       r.Meta.Source,
		Schema:       r.Meta.Schema,
		Table:        r.Meta.Table,
		Operation:    r.Meta.Operation,
		Offset:       r.Meta.Offset,
		PartitionKey: r.Meta.PartitionKey,
	}
	
	if r.Meta.Extra != nil {
		meta.Extra = make(map[string]interface{})
		for k, v := range r.Meta.Extra {
			meta.Extra[k] = v
		}
	}
	
	return &Record{
		Data:      data,
		Meta:      meta,
		Timestamp: r.Timestamp,
	}
}

// GetField 获取字段值
func (r *Record) GetField(name string) (interface{}, bool) {
	val, ok := r.Data[name]
	return val, ok
}

// SetField 设置字段值
func (r *Record) SetField(name string, value interface{}) {
	r.Data[name] = value
}

// GetString 获取字符串字段
func (r *Record) GetString(name string) (string, bool) {
	val, ok := r.Data[name]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// GetInt64 获取整数字段
func (r *Record) GetInt64(name string) (int64, bool) {
	val, ok := r.Data[name]
	if !ok {
		return 0, false
	}
	
	switch v := val.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

// GetFloat64 获取浮点数字段
func (r *Record) GetFloat64(name string) (float64, bool) {
	val, ok := r.Data[name]
	if !ok {
		return 0, false
	}
	
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

// GetBool 获取布尔字段
func (r *Record) GetBool(name string) (bool, bool) {
	val, ok := r.Data[name]
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

// ToJSON 转换为 JSON 字符串
func (r *Record) ToJSON() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从 JSON 字符串创建记录
func FromJSON(jsonStr string) (*Record, error) {
	var record Record
	err := json.Unmarshal([]byte(jsonStr), &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Size 计算记录的近似大小(字节)
func (r *Record) Size() int {
	// 简单估算: JSON 序列化后的大小
	bytes, err := json.Marshal(r)
	if err != nil {
		return 0
	}
	return len(bytes)
}
