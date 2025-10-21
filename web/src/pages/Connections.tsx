import { Button, Table, Tag, Space } from 'antd'
import { PlusOutlined, CheckCircleOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'

interface Connection {
  id: string
  name: string
  type: string
  status: string
  lastTested?: string
}

const columns: ColumnsType<Connection> = [
  {
    title: '连接名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    render: (type: string) => <Tag color="blue">{type}</Tag>,
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (status: string) => {
      const color = status === 'success' ? 'success' : 'default'
      const text = status === 'success' ? '正常' : '未测试'
      return <Tag color={color}>{text}</Tag>
    },
  },
  {
    title: '最后测试',
    dataIndex: 'lastTested',
    key: 'lastTested',
  },
  {
    title: '操作',
    key: 'action',
    render: () => (
      <Space size="small">
        <Button type="link" size="small" icon={<CheckCircleOutlined />}>
          测试连接
        </Button>
        <Button type="link" size="small">
          编辑
        </Button>
        <Button type="link" size="small" danger>
          删除
        </Button>
      </Space>
    ),
  },
]

const mockData: Connection[] = [
  {
    id: '1',
    name: '生产环境 MySQL',
    type: 'MySQL',
    status: 'success',
    lastTested: '2025-10-21 09:00:00',
  },
  {
    id: '2',
    name: '分析库 PostgreSQL',
    type: 'PostgreSQL',
    status: 'success',
    lastTested: '2025-10-21 09:00:00',
  },
]

export default function Connections() {
  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">连接配置</h1>
        <Button type="primary" icon={<PlusOutlined />}>
          新建连接
        </Button>
      </div>

      <Table columns={columns} dataSource={mockData} rowKey="id" />
    </div>
  )
}
