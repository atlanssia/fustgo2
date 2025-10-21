import { useState, useCallback } from 'react'
import { Button, Card, Space, Input, Select, Form, message, Modal, List } from 'antd'
import { 
  DragDropContext, 
  Droppable, 
  Draggable, 
  DropResult 
} from '@hello-pangea/dnd'
import { 
  PlusOutlined, 
  DeleteOutlined, 
  MenuOutlined,
  DatabaseOutlined,
  ApiOutlined,
  FileOutlined,
  PlayCircleOutlined,
  SaveOutlined
} from '@ant-design/icons'

interface PipelineItem {
  id: string
  type: 'source' | 'processor' | 'sink'
  plugin: string
  name: string
  config: Record<string, any>
}

interface PluginTemplate {
  id: string
  name: string
  type: 'source' | 'processor' | 'sink'
  plugin: string
  icon: React.ReactNode
  description: string
  defaultConfig: Record<string, any>
}

const pluginTemplates: PluginTemplate[] = [
  { 
    id: 'mysql-source', 
    name: 'MySQL 数据源', 
    type: 'source',
    plugin: 'mysql',
    icon: <DatabaseOutlined />, 
    description: '从 MySQL 数据库读取数据',
    defaultConfig: { host: 'localhost', port: 3306, database: '', username: '', password: '' }
  },
  { 
    id: 'pg-source', 
    name: 'PostgreSQL 数据源', 
    type: 'source',
    plugin: 'postgresql',
    icon: <DatabaseOutlined />, 
    description: '从 PostgreSQL 数据库读取数据',
    defaultConfig: { host: 'localhost', port: 5432, database: '', username: '', password: '' }
  },
  { 
    id: 'csv-source', 
    name: 'CSV 文件源', 
    type: 'source',
    plugin: 'csv',
    icon: <FileOutlined />, 
    description: '从 CSV 文件读取数据',
    defaultConfig: { path: '', delimiter: ',' }
  },
  { 
    id: 'field-mapper', 
    name: '字段映射器', 
    type: 'processor',
    plugin: 'field_mapper',
    icon: <ApiOutlined />, 
    description: '字段名称映射和转换',
    defaultConfig: { mappings: {} }
  },
  { 
    id: 'filter', 
    name: '数据过滤器', 
    type: 'processor',
    plugin: 'filter',
    icon: <ApiOutlined />, 
    description: '根据条件过滤数据',
    defaultConfig: { condition: '' }
  },
  { 
    id: 'mysql-sink', 
    name: 'MySQL 数据目标', 
    type: 'sink',
    plugin: 'mysql',
    icon: <DatabaseOutlined />, 
    description: '将数据写入 MySQL 数据库',
    defaultConfig: { host: 'localhost', port: 3306, database: '', username: '', password: '', table: '' }
  },
  { 
    id: 'pg-sink', 
    name: 'PostgreSQL 数据目标', 
    type: 'sink',
    plugin: 'postgresql',
    icon: <DatabaseOutlined />, 
    description: '将数据写入 PostgreSQL 数据库',
    defaultConfig: { host: 'localhost', port: 5432, database: '', username: '', password: '', table: '' }
  },
  { 
    id: 'csv-sink', 
    name: 'CSV 文件目标', 
    type: 'sink',
    plugin: 'csv',
    icon: <FileOutlined />, 
    description: '将数据写入 CSV 文件',
    defaultConfig: { path: '', delimiter: ',' }
  },
]

export default function PipelineDesigner() {
  const [items, setItems] = useState<PipelineItem[]>([
    {
      id: '1',
      type: 'source',
      plugin: 'mysql',
      name: 'MySQL 源',
      config: { host: 'localhost', port: 3306, database: 'testdb', username: 'user' }
    },
    {
      id: '2',
      type: 'processor',
      plugin: 'field_mapper',
      name: '字段映射',
      config: { mappings: { id: 'user_id', name: 'user_name' } }
    },
    {
      id: '3',
      type: 'sink',
      plugin: 'postgresql',
      name: 'PostgreSQL 目标',
      config: { host: 'localhost', port: 5432, database: 'analytics', username: 'user', table: 'users' }
    }
  ])
  
  const [selectedItem, setSelectedItem] = useState<PipelineItem | null>(null)
  const [isTemplateModalVisible, setIsTemplateModalVisible] = useState(false)
  const [pipelineName, setPipelineName] = useState('我的数据管道')
  const [form] = Form.useForm()

  const onDragEnd = (result: DropResult) => {
    if (!result.destination) return
    
    const newItems = Array.from(items)
    const [removed] = newItems.splice(result.source.index, 1)
    newItems.splice(result.destination.index, 0, removed)
    
    setItems(newItems)
  }

  const addItemFromTemplate = (template: PluginTemplate) => {
    const newItem: PipelineItem = {
      id: Date.now().toString(),
      type: template.type,
      plugin: template.plugin,
      name: template.name,
      config: { ...template.defaultConfig }
    }
    setItems([...items, newItem])
    setIsTemplateModalVisible(false)
  }

  const removeItem = (id: string) => {
    setItems(items.filter(item => item.id !== id))
    if (selectedItem?.id === id) {
      setSelectedItem(null)
    }
  }

  const updateItem = (id: string, updates: Partial<PipelineItem>) => {
    setItems(items.map(item => 
      item.id === id ? { ...item, ...updates } : item
    ))
  }

  const selectItem = (item: PipelineItem) => {
    setSelectedItem(item)
    form.setFieldsValue(item.config)
  }

  const saveConfig = () => {
    message.success('管道配置保存成功')
    console.log('Saved pipeline:', { name: pipelineName, items })
  }

  const runPipeline = () => {
    message.success('管道开始执行')
    console.log('Run pipeline:', { name: pipelineName, items })
  }

  const handleConfigChange = useCallback((changedValues: any) => {
    if (selectedItem) {
      updateItem(selectedItem.id, { 
        config: { ...selectedItem.config, ...changedValues } 
      })
    }
  }, [selectedItem])

  const getTypeTitle = (type: string) => {
    switch (type) {
      case 'source': return '数据源'
      case 'processor': return '数据处理'
      case 'sink': return '数据目标'
      default: return '未知'
    }
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'source': return 'blue'
      case 'processor': return 'green'
      case 'sink': return 'purple'
      default: return 'gray'
    }
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <Input 
          value={pipelineName} 
          onChange={e => setPipelineName(e.target.value)}
          className="text-2xl font-bold max-w-md"
          bordered={false}
        />
        <Space>
          <Button icon={<PlayCircleOutlined />} onClick={runPipeline}>
            执行管道
          </Button>
          <Button type="primary" icon={<SaveOutlined />} onClick={saveConfig}>
            保存配置
          </Button>
        </Space>
      </div>

      <div className="grid grid-cols-12 gap-6">
        {/* 左侧工具栏 */}
        <div className="col-span-2">
          <Card title="添加组件" size="small" className="mb-4">
            <Space direction="vertical" className="w-full">
              <Button 
                block 
                icon={<DatabaseOutlined />} 
                onClick={() => setIsTemplateModalVisible(true)}
              >
                添加组件
              </Button>
            </Space>
          </Card>
          
          <Card title="管道结构" size="small">
            <List
              size="small"
              dataSource={items}
              renderItem={(item, index) => (
                <List.Item
                  onClick={() => selectItem(item)}
                  className={`cursor-pointer ${selectedItem?.id === item.id ? 'bg-blue-50' : ''}`}
                >
                  <List.Item.Meta
                    avatar={<MenuOutlined />}
                    title={
                      <div>
                        <span className="font-medium">{item.name}</span>
                        <div className="text-xs text-gray-500">
                          {getTypeTitle(item.type)}
                        </div>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </div>

        {/* 中间拖拽区域 */}
        <div className="col-span-7">
          <Card title="管道流程" size="small" className="h-full">
            <DragDropContext onDragEnd={onDragEnd}>
              <Droppable droppableId="pipeline-items">
                {(provided) => (
                  <div 
                    {...provided.droppableProps} 
                    ref={provided.innerRef}
                    className="min-h-[500px] p-4"
                  >
                    {items.map((item, index) => (
                      <Draggable key={item.id} draggableId={item.id} index={index}>
                        {(provided) => (
                          <div
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            className="mb-4"
                          >
                            <Card 
                              size="small" 
                              className={`cursor-pointer ${selectedItem?.id === item.id ? 'border-blue-500 border-2' : ''}`}
                              onClick={() => selectItem(item)}
                            >
                              <div className="flex items-center justify-between">
                                <div className="flex items-center">
                                  <div {...provided.dragHandleProps} className="mr-2">
                                    <MenuOutlined className="cursor-move" />
                                  </div>
                                  <div>
                                    <div className="font-medium">{item.name}</div>
                                    <div className="text-sm text-gray-500">
                                      {item.plugin} 插件
                                    </div>
                                  </div>
                                </div>
                                <Button 
                                  type="text" 
                                  icon={<DeleteOutlined />} 
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    removeItem(item.id)
                                  }}
                                />
                              </div>
                            </Card>
                          </div>
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}
                    
                    {items.length === 0 && (
                      <div className="text-center text-gray-400 py-20">
                        <DatabaseOutlined className="text-4xl mb-2" />
                        <p>拖拽组件到此处构建数据管道</p>
                        <Button 
                          type="primary" 
                          className="mt-4"
                          onClick={() => setIsTemplateModalVisible(true)}
                        >
                          添加第一个组件
                        </Button>
                      </div>
                    )}
                  </div>
                )}
              </Droppable>
            </DragDropContext>
          </Card>
        </div>

        {/* 右侧配置面板 */}
        <div className="col-span-3">
          {selectedItem ? (
            <Card title={`${selectedItem.name} 配置`} size="small">
              <Form
                form={form}
                layout="vertical"
                initialValues={selectedItem.config}
                onValuesChange={handleConfigChange}
              >
                {selectedItem.plugin === 'mysql' && selectedItem.type === 'source' && (
                  <>
                    <Form.Item label="主机地址" name="host">
                      <Input placeholder="localhost" />
                    </Form.Item>
                    <Form.Item label="端口" name="port">
                      <Input type="number" placeholder="3306" />
                    </Form.Item>
                    <Form.Item label="数据库名" name="database">
                      <Input placeholder="database_name" />
                    </Form.Item>
                    <Form.Item label="用户名" name="username">
                      <Input placeholder="username" />
                    </Form.Item>
                    <Form.Item label="密码" name="password">
                      <Input.Password placeholder="password" />
                    </Form.Item>
                    <Form.Item label="表名" name="table">
                      <Input placeholder="table_name" />
                    </Form.Item>
                  </>
                )}
                
                {selectedItem.plugin === 'postgresql' && selectedItem.type === 'sink' && (
                  <>
                    <Form.Item label="主机地址" name="host">
                      <Input placeholder="localhost" />
                    </Form.Item>
                    <Form.Item label="端口" name="port">
                      <Input type="number" placeholder="5432" />
                    </Form.Item>
                    <Form.Item label="数据库名" name="database">
                      <Input placeholder="database_name" />
                    </Form.Item>
                    <Form.Item label="用户名" name="username">
                      <Input placeholder="username" />
                    </Form.Item>
                    <Form.Item label="密码" name="password">
                      <Input.Password placeholder="password" />
                    </Form.Item>
                    <Form.Item label="表名" name="table">
                      <Input placeholder="table_name" />
                    </Form.Item>
                  </>
                )}
                
                {selectedItem.plugin === 'field_mapper' && (
                  <>
                    <Form.Item label="字段映射" name="mappings">
                      <Input.TextArea 
                        placeholder='{"source_field": "target_field"}' 
                        rows={4} 
                      />
                    </Form.Item>
                  </>
                )}
                
                {selectedItem.plugin === 'csv' && (
                  <>
                    <Form.Item label="文件路径" name="path">
                      <Input placeholder="/path/to/file.csv" />
                    </Form.Item>
                    <Form.Item label="分隔符" name="delimiter">
                      <Input placeholder="," />
                    </Form.Item>
                  </>
                )}
              </Form>
            </Card>
          ) : (
            <Card title="配置面板" size="small">
              <div className="text-center text-gray-400 py-10">
                <ApiOutlined className="text-4xl mb-2" />
                <p>选择一个组件进行配置</p>
              </div>
            </Card>
          )}
        </div>
      </div>

      {/* 添加组件模态框 */}
      <Modal
        title="选择组件"
        open={isTemplateModalVisible}
        onCancel={() => setIsTemplateModalVisible(false)}
        footer={null}
        width={800}
      >
        <div className="grid grid-cols-3 gap-4">
          {pluginTemplates.map(template => (
            <Card
              key={template.id}
              className="cursor-pointer hover:shadow-md transition-shadow"
              onClick={() => addItemFromTemplate(template)}
            >
              <div className="text-center">
                <div className="text-2xl mb-2">{template.icon}</div>
                <div className="font-medium">{template.name}</div>
                <div className="text-sm text-gray-500 mt-1">{template.description}</div>
              </div>
            </Card>
          ))}
        </div>
      </Modal>
    </div>
  )
}