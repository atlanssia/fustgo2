import { Table, Tag } from 'antd'
import type { ColumnsType } from 'antd/es/table'

interface Execution {
  id: string
  jobName: string
  status: string
  startTime: string
  endTime?: string
  duration?: string
  recordsRead: number
  recordsWritten: number
}

const columns: ColumnsType<Execution> = [
  {
    title: '任务名称',
    dataIndex: 'jobName',
    key: 'jobName',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (status: string) => {
      const colorMap: Record<string, string> = {
        succeeded: 'success',
        failed: 'error',
        running: 'processing',
        cancelled: 'default',
      }
      const textMap: Record<string, string> = {
        succeeded: '成功',
        failed: '失败',
        running: '运行中',
        cancelled: '已取消',
      }
      return <Tag color={colorMap[status]}>{textMap[status]}</Tag>
    },
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
    key: 'endTime',
  },
  {
    title: '耗时',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '读取记录',
    dataIndex: 'recordsRead',
    key: 'recordsRead',
  },
  {
    title: '写入记录',
    dataIndex: 'recordsWritten',
    key: 'recordsWritten',
  },
]

const mockData: Execution[] = [
  {
    id: '1',
    jobName: 'MySQL 到 PostgreSQL 用户数据同步',
    status: 'succeeded',
    startTime: '2025-10-21 10:00:00',
    endTime: '2025-10-21 10:05:32',
    duration: '5分32秒',
    recordsRead: 10000,
    recordsWritten: 9950,
  },
]

export default function Executions() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">执行历史</h1>
      <Table columns={columns} dataSource={mockData} rowKey="id" />
    </div>
  )
}
