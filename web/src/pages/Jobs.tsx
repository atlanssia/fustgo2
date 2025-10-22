import { Button, Table, Tag, Space } from 'antd'
import { PlusOutlined, PlayCircleOutlined, PauseCircleOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'

interface Job {
  id: string
  name: string
  status: string
  enabled: boolean
  lastExecutionTime?: string
  nextExecutionTime?: string
}

const columns: ColumnsType<Job> = [
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (status: string) => {
      const color = status === 'running' ? 'green' : status === 'failed' ? 'red' : 'blue'
      return <Tag color={color}>{status}</Tag>
    },
  },
  {
    title: '启用状态',
    dataIndex: 'enabled',
    key: 'enabled',
    render: (enabled: boolean) => (
      <Tag color={enabled ? 'success' : 'default'}>{enabled ? '已启用' : '已禁用'}</Tag>
    ),
  },
  {
    title: '最后执行',
    dataIndex: 'lastExecutionTime',
    key: 'lastExecutionTime',
  },
  {
    title: '下次执行',
    dataIndex: 'nextExecutionTime',
    key: 'nextExecutionTime',
  },
  {
    title: '操作',
    key: 'action',
    render: () => (
      <Space size="small">
        <Button type="link" size="small" icon={<PlayCircleOutlined />}>
          执行
        </Button>
        <Button type="link" size="small" icon={<PauseCircleOutlined />}>
          暂停
        </Button>
        <Button type="link" size="small">
          编辑
        </Button>
      </Space>
    ),
  },
]

const mockData: Job[] = [
  {
    id: '1',
    name: 'MySQL 到 PostgreSQL 用户数据同步',
    status: 'succeeded',
    enabled: true,
    lastExecutionTime: '2025-10-21 10:00:00',
    nextExecutionTime: '2025-10-22 02:00:00',
  },
]

export default function Jobs() {
  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">任务管理</h1>
        <Button type="primary" icon={<PlusOutlined />}>
          新建任务
        </Button>
      </div>

      <Table columns={columns} dataSource={mockData} rowKey="id" />
    </div>
  )
}
