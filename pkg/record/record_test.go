package record

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	data := map[string]interface{}{
		"id":   1,
		"name": "test",
	}
	
	record := NewRecord(data)
	
	if record == nil {
		t.Fatal("NewRecord returned nil")
	}
	
	if record.Data == nil {
		t.Error("Record.Data is nil")
	}
	
	if record.Meta == nil {
		t.Error("Record.Meta is nil")
	}
	
	if record.Timestamp.IsZero() {
		t.Error("Record.Timestamp is zero")
	}
}

func TestRecordClone(t *testing.T) {
	original := NewRecord(map[string]interface{}{
		"id":   1,
		"name": "test",
	})
	original.Meta.Source = "mysql"
	original.Meta.Table = "users"
	
	cloned := original.Clone()
	
	// 验证克隆的数据
	if cloned.Data["id"] != original.Data["id"] {
		t.Error("Cloned data mismatch")
	}
	
	// 修改克隆的数据不应影响原始数据
	cloned.Data["name"] = "modified"
	if original.Data["name"] == "modified" {
		t.Error("Modifying cloned data affected original")
	}
	
	// 验证元数据
	if cloned.Meta.Source != "mysql" {
		t.Error("Cloned meta source mismatch")
	}
	if cloned.Meta.Table != "users" {
		t.Error("Cloned meta table mismatch")
	}
}

func TestRecordGetField(t *testing.T) {
	record := NewRecord(map[string]interface{}{
		"id":   int64(123),
		"name": "test",
		"age":  25,
	})
	
	// 测试获取存在的字段
	val, ok := record.GetField("name")
	if !ok {
		t.Error("GetField failed for existing field")
	}
	if val != "test" {
		t.Errorf("GetField returned wrong value: %v", val)
	}
	
	// 测试获取不存在的字段
	_, ok = record.GetField("nonexistent")
	if ok {
		t.Error("GetField should return false for nonexistent field")
	}
}

func TestRecordGetString(t *testing.T) {
	record := NewRecord(map[string]interface{}{
		"name":   "test",
		"number": 123,
	})
	
	// 测试获取字符串字段
	str, ok := record.GetString("name")
	if !ok {
		t.Error("GetString failed for string field")
	}
	if str != "test" {
		t.Errorf("GetString returned wrong value: %s", str)
	}
	
	// 测试获取非字符串字段
	_, ok = record.GetString("number")
	if ok {
		t.Error("GetString should return false for non-string field")
	}
}

func TestRecordGetInt64(t *testing.T) {
	record := NewRecord(map[string]interface{}{
		"int64":   int64(123),
		"int":     int(456),
		"float64": float64(789.0),
		"string":  "not a number",
	})
	
	tests := []struct {
		name     string
		field    string
		expected int64
		ok       bool
	}{
		{"int64 field", "int64", 123, true},
		{"int field", "int", 456, true},
		{"float64 field", "float64", 789, true},
		{"string field", "string", 0, false},
		{"nonexistent field", "nonexistent", 0, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := record.GetInt64(tt.field)
			if ok != tt.ok {
				t.Errorf("GetInt64(%s) ok = %v, want %v", tt.field, ok, tt.ok)
			}
			if val != tt.expected {
				t.Errorf("GetInt64(%s) = %v, want %v", tt.field, val, tt.expected)
			}
		})
	}
}

func TestRecordGetFloat64(t *testing.T) {
	record := NewRecord(map[string]interface{}{
		"float64": float64(123.45),
		"float32": float32(67.89),
		"int":     int(100),
		"string":  "not a number",
	})
	
	tests := []struct {
		name     string
		field    string
		expected float64
		ok       bool
	}{
		{"float64 field", "float64", 123.45, true},
		{"int field", "int", 100.0, true},
		{"string field", "string", 0, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := record.GetFloat64(tt.field)
			if ok != tt.ok {
				t.Errorf("GetFloat64(%s) ok = %v, want %v", tt.field, ok, tt.ok)
			}
			if tt.ok && val != tt.expected {
				t.Errorf("GetFloat64(%s) = %v, want %v", tt.field, val, tt.expected)
			}
		})
	}
}

func TestRecordGetBool(t *testing.T) {
	record := NewRecord(map[string]interface{}{
		"bool_true":  true,
		"bool_false": false,
		"string":     "not a bool",
	})
	
	tests := []struct {
		name     string
		field    string
		expected bool
		ok       bool
	}{
		{"true field", "bool_true", true, true},
		{"false field", "bool_false", false, true},
		{"string field", "string", false, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := record.GetBool(tt.field)
			if ok != tt.ok {
				t.Errorf("GetBool(%s) ok = %v, want %v", tt.field, ok, tt.ok)
			}
			if val != tt.expected {
				t.Errorf("GetBool(%s) = %v, want %v", tt.field, val, tt.expected)
			}
		})
	}
}

func TestRecordSetField(t *testing.T) {
	record := NewRecord(map[string]interface{}{})
	
	record.SetField("name", "test")
	record.SetField("age", 25)
	
	if record.Data["name"] != "test" {
		t.Error("SetField failed to set name")
	}
	if record.Data["age"] != 25 {
		t.Error("SetField failed to set age")
	}
}

func TestRecordToJSON(t *testing.T) {
	now := time.Now()
	record := &Record{
		Data: map[string]interface{}{
			"id":   1,
			"name": "test",
		},
		Meta: &RecordMeta{
			Source: "mysql",
			Table:  "users",
		},
		Timestamp: now,
	}
	
	jsonStr, err := record.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	
	// 验证 JSON 字符串不为空
	if jsonStr == "" {
		t.Error("ToJSON returned empty string")
	}
	
	// 验证可以解析回来
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}
}

func TestRecordFromJSON(t *testing.T) {
	jsonStr := `{
		"data": {"id": 1, "name": "test"},
		"meta": {"source": "mysql", "table": "users"},
		"timestamp": "2025-10-21T10:00:00Z"
	}`
	
	record, err := FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}
	
	if record == nil {
		t.Fatal("FromJSON returned nil")
	}
	
	// 验证数据
	if record.Data["name"] != "test" {
		t.Error("FromJSON data mismatch")
	}
	
	// 验证元数据
	if record.Meta.Source != "mysql" {
		t.Error("FromJSON meta source mismatch")
	}
	if record.Meta.Table != "users" {
		t.Error("FromJSON meta table mismatch")
	}
}

func TestRecordSize(t *testing.T) {
	record := NewRecord(map[string]interface{}{
		"id":   1,
		"name": "test",
	})
	
	size := record.Size()
	if size <= 0 {
		t.Error("Size should be greater than 0")
	}
	
	// 添加更多数据应该增加大小
	record.SetField("description", "This is a long description that should increase the size")
	newSize := record.Size()
	if newSize <= size {
		t.Error("Size should increase after adding more data")
	}
}

func BenchmarkNewRecord(b *testing.B) {
	data := map[string]interface{}{
		"id":   1,
		"name": "test",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewRecord(data)
	}
}

func BenchmarkRecordClone(b *testing.B) {
	record := NewRecord(map[string]interface{}{
		"id":   1,
		"name": "test",
		"age":  25,
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		record.Clone()
	}
}

func BenchmarkRecordToJSON(b *testing.B) {
	record := NewRecord(map[string]interface{}{
		"id":   1,
		"name": "test",
		"age":  25,
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = record.ToJSON()
	}
}

func BenchmarkRecordFromJSON(b *testing.B) {
	jsonStr := `{"data":{"id":1,"name":"test","age":25},"meta":{},"timestamp":"2025-10-21T10:00:00Z"}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromJSON(jsonStr)
	}
}
