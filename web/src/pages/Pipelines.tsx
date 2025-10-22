import { useState } from 'react'
import { Button, Table, Tag, Space, Modal, Form, Input, message, Switch } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, PlayCircleOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useNavigate } from 'react-router-dom'

interface Pipeline {
  id: string
  name: string
  description: string
  status: string
  enabled: boolean
  createdAt: string
  lastExecutionTime?: string
}

const mockPipelines: Pipeline[] = [
  {
    id: '1',
    name: '用户数据同步',
    description: '从MySQL同步用户数据到PostgreSQL',
    status: 'active',
    enabled: true,
    createdAt: '2025-10-21 10:00:00',
    lastExecutionTime: '2025-10-21 14:30:00',
  },
  {
    id: '2',
    name: '订单数据同步',
    description: '从MongoDB同步订单数据到Elasticsearch',
    status: 'draft',
    enabled: false,
    createdAt: '2025-10-21 11:00:00',
  },
]

export default function Pipelines() {
  const [pipelines, setPipelines] = useState<Pipeline[]>(mockPipelines)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingPipeline, setEditingPipeline] = useState<Pipeline | null>(null)
  const [form] = Form.useForm()
  const navigate = useNavigate()

  const columns: ColumnsType<Pipeline> = [
    {
      title: '管道名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        let color = 'default'
        if (status === 'active') color = 'success'
        if (status === 'draft') color = 'warning'
        return <Tag color={color}>{status}</Tag>
      },
    },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'success' : 'default'}>
          {enabled ? '是' : '否'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '最后执行',
      dataIndex: 'lastExecutionTime',
      key: 'lastExecutionTime',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="small">
          <Button type="link" size="small" onClick={() => navigate(`/config?id=${record.id}`)}>
            配置
          </Button>
          <Button type="link" size="small" icon={<PlayCircleOutlined />} onClick={() => executePipeline(record)}>
            执行
          </Button>
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

  const executePipeline = (pipeline: Pipeline) => {
    message.success(`管道 "${pipeline.name}" 开始执行`)
  }

  const handleEdit = (pipeline: Pipeline) => {
    setEditingPipeline(pipeline)
    form.setFieldsValue(pipeline)
    setIsModalVisible(true)
  }

  const handleDelete = (pipeline: Pipeline) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除管道 "${pipeline.name}" 吗？`,
      onOk: () => {
        setPipelines(pipelines.filter(p => p.id !== pipeline.id))
        message.success('删除成功')
      },
    })
  }

  const handleModalOk = () => {
    form.validateFields().then(values => {
      if (editingPipeline) {
        // 更新管道
        setPipelines(pipelines.map(p => 
          p.id === editingPipeline.id ? { ...p, ...values } : p
        ))
        message.success('更新成功')
      } else {
        // 创建管道
        const newPipeline: Pipeline = {
          id: (pipelines.length + 1).toString(),
          ...values,
          status: 'draft',
          createdAt: new Date().toISOString().slice(0, 19).replace('T', ' '),
        }
        setPipelines([...pipelines, newPipeline])
        message.success('创建成功')
      }
      setIsModalVisible(false)
      form.resetFields()
      setEditingPipeline(null)
    })
  }

  const handleModalCancel = () => {
    setIsModalVisible(false)
    form.resetFields()
    setEditingPipeline(null)
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">数据管道</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalVisible(true)}>
          创建管道
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={pipelines} 
        rowKey="id" 
        pagination={{
          pageSize: 10,
        }}
      />

      <Modal
        title={editingPipeline ? "编辑管道" : "创建管道"}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="管道名称"
            rules={[{ required: true, message: '请输入管道名称' }]}
          >
            <Input placeholder="请输入管道名称" />
          </Form.Item>
          
          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea placeholder="请输入管道描述" rows={3} />
          </Form.Item>
          
          <Form.Item
            name="enabled"
            label="是否启用"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}