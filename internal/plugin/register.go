package plugin

import (
	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/plugin/processor/filter"
	"github.com/fustgo/fustgo2/pkg/plugin/processor/json"
	"github.com/fustgo/fustgo2/pkg/plugin/reader/mongodb"
	"github.com/fustgo/fustgo2/pkg/plugin/reader/mysql"
	"github.com/fustgo/fustgo2/pkg/plugin/writer/elasticsearch"
	"github.com/fustgo/fustgo2/pkg/plugin/writer/postgresql"
)

// RegisterPlugins 注册所有内置插件
func RegisterPlugins() error {
	// 注册MySQL读取器
	if err := plugin.Register("mysql-reader", func() plugin.Plugin {
		return mysql.NewMySQLReader()
	}, mysql.NewMySQLReader().GetInfo()); err != nil {
		return err
	}

	// 注册MongoDB读取器
	if err := plugin.Register("mongodb-reader", func() plugin.Plugin {
		return mongodb.NewMongoDBReader()
	}, mongodb.NewMongoDBReader().GetInfo()); err != nil {
		return err
	}

	// 注册PostgreSQL写入器
	if err := plugin.Register("postgresql-writer", func() plugin.Plugin {
		return postgresql.NewPostgreSQLWriter()
	}, postgresql.NewPostgreSQLWriter().GetInfo()); err != nil {
		return err
	}

	// 注册Elasticsearch写入器
	if err := plugin.Register("elasticsearch-writer", func() plugin.Plugin {
		return elasticsearch.NewElasticsearchWriter()
	}, elasticsearch.NewElasticsearchWriter().GetInfo()); err != nil {
		return err
	}

	// 注册过滤处理器
	if err := plugin.Register("filter-processor", func() plugin.Plugin {
		return filter.NewFilterProcessor()
	}, filter.NewFilterProcessor().GetInfo()); err != nil {
		return err
	}

	// 注册JSON处理器
	if err := plugin.Register("json-processor", func() plugin.Plugin {
		return json.NewJSONProcessor()
	}, json.NewJSONProcessor().GetInfo()); err != nil {
		return err
	}

	return nil
}