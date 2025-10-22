package types

import (
	"time"
)

// PluginType 插件类型
type PluginType string

const (
	PluginTypeReader    PluginType = "reader"
	PluginTypeWriter    PluginType = "writer"
	PluginTypeProcessor PluginType = "processor"
)

// Plugin 插件信息
type Plugin struct {
	// ID 插件唯一标识
	ID string `json:"id" gorm:"primaryKey"`
	
	// Name 插件名称
	Name string `json:"name" gorm:"uniqueIndex;not null"`
	
	// Type 插件类型
	Type PluginType `json:"type" gorm:"index;not null"`
	
	// Version 版本
	Version string `json:"version"`
	
	// Description 描述
	Description string `json:"description"`
	
	// Author 作者
	Author string `json:"author"`
	
	// ConfigSchema 配置模式(JSON Schema)
	ConfigSchema string `json:"config_schema" gorm:"type:text"`
	
	// Enabled 是否启用
	Enabled bool `json:"enabled" gorm:"default:true"`
	
	// BuiltIn 是否内置插件
	BuiltIn bool `json:"built_in" gorm:"default:false"`
	
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Plugin) TableName() string {
	return "plugins"
}

// PluginInstance 插件实例配置
type PluginInstance struct {
	// ID 实例唯一标识
	ID string `json:"id" gorm:"primaryKey"`
	
	// Name 实例名称
	Name string `json:"name" gorm:"uniqueIndex;not null"`
	
	// PluginID 关联的插件ID
	PluginID string `json:"plugin_id" gorm:"index;not null"`
	
	// Config 配置信息(JSON格式)
	Config string `json:"config" gorm:"type:text"`
	
	// Description 描述
	Description string `json:"description"`
	
	// Tags 标签
	Tags string `json:"tags" gorm:"type:text"`
	
	// Enabled 是否启用
	Enabled bool `json:"enabled" gorm:"default:true"`
	
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// CreatedBy 创建人
	CreatedBy string `json:"created_by"`
}

// TableName 指定表名
func (PluginInstance) TableName() string {
	return "plugin_instances"
}