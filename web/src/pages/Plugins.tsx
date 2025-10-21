import { useState, useEffect } from 'react'
import { Button, Table, Tag, Space, Modal, Form, Input, Select, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'

interface Plugin {
  id: string
  name: string
  type: string
  version: string
  description: string
  enabled: boolean
  builtIn: boolean
  createdAt: string
}

const mockPlugins: Plugin[] = [
  {
    id: '1',
    name: 'MySQL Reader',
    type: 'reader',
    version: '1.0.0',
    description: 'MySQL数据库读取插件',
    enabled: true,
    builtIn: true,
    createdAt: '2025-10-21 10:00:00',
  },
  {
    id: '2',
    name: 'PostgreSQL Writer',
    type: 'writer',
    version: '1.0.0',
    description: 'PostgreSQL数据库写入插件',
    enabled: true,
    builtIn: true,
    createdAt: '2025-10-21 10:00:00',
  },
  {
    id: '3',
    name: 'Field Mapper',
    type: 'processor',
    version: '1.0.0',
    description: '字段映射处理器',
    enabled: true,
    builtIn: true,
    createdAt: '2025-10-21 10:00:00',
  },
]

export default function Plugins() {
  const [plugins, setPlugins] = useState<Plugin[]>(mockPlugins)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingPlugin, setEditingPlugin] = useState<Plugin | null>(null)
  const [form] = Form.useForm()

  const columns: ColumnsType<Plugin> = [
    {
      title: '插件名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        let color = 'blue'
        if (type === 'writer') color = 'green'
        if (type === 'processor') color = 'purple'
        return <Tag color={color}>{type}</Tag>
      },
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'success' : 'default'}>
          {enabled ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '内置',
      dataIndex: 'builtIn',
      key: 'builtIn',
      render: (builtIn: boolean) => (
        <Tag color={builtIn ? 'processing' : 'default'}>
          {builtIn ? '是' : '否'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const handleEdit = (plugin: Plugin) => {
    setEditingPlugin(plugin)
    form.setFieldsValue(plugin)
    setIsModalVisible(true)
  }

  const handleDelete = (plugin: Plugin) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除插件 "${plugin.name}" 吗？`,
      onOk: () => {
        setPlugins(plugins.filter(p => p.id !== plugin.id))
        message.success('删除成功')
      },
    })
  }

  const handleModalOk = () => {
    form.validateFields().then(values => {
      if (editingPlugin) {
        // 更新插件
        setPlugins(plugins.map(p => 
          p.id === editingPlugin.id ? { ...p, ...values } : p
        ))
        message.success('更新成功')
      } else {
        // 创建插件
        const newPlugin: Plugin = {
          id: (plugins.length + 1).toString(),
          ...values,
          createdAt: new Date().toISOString().slice(0, 19).replace('T', ' '),
        }
        setPlugins([...plugins, newPlugin])
        message.success('创建成功')
      }
      setIsModalVisible(false)
      form.resetFields()
      setEditingPlugin(null)
    })
  }

  const handleModalCancel = () => {
    setIsModalVisible(false)
    form.resetFields()
    setEditingPlugin(null)
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">插件管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalVisible(true)}>
          添加插件
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={plugins} 
        rowKey="id" 
        pagination={{
          pageSize: 10,
        }}
      />

      <Modal
        title={editingPlugin ? "编辑插件" : "添加插件"}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="插件名称"
            rules={[{ required: true, message: '请输入插件名称' }]}
          >
            <Input placeholder="请输入插件名称" />
          </Form.Item>
          
          <Form.Item
            name="type"
            label="插件类型"
            rules={[{ required: true, message: '请选择插件类型' }]}
          >
            <Select placeholder="请选择插件类型">
              <Select.Option value="reader">Reader</Select.Option>
              <Select.Option value="writer">Writer</Select.Option>
              <Select.Option value="processor">Processor</Select.Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="version"
            label="版本"
            rules={[{ required: true, message: '请输入版本号' }]}
          >
            <Input placeholder="请输入版本号" />
          </Form.Item>
          
          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea placeholder="请输入插件描述" rows={3} />
          </Form.Item>
          
          <Form.Item
            name="enabled"
            label="是否启用"
            initialValue={true}
          >
            <Select>
              <Select.Option value={true}>启用</Select.Option>
              <Select.Option value={false}>禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}