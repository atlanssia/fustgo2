import { Card, Row, Col, Statistic } from 'antd'
import { DatabaseOutlined, CheckCircleOutlined, CloseCircleOutlined, SyncOutlined } from '@ant-design/icons'

export default function Dashboard() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">仪表板</h1>
      
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总任务数"
              value={12}
              prefix={<DatabaseOutlined />}
              valueStyle={{ color: '#1677ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="运行中"
              value={3}
              prefix={<SyncOutlined spin />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="成功"
              value={156}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="失败"
              value={8}
              prefix={<CloseCircleOutlined />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      <Card title="快速开始" className="mt-6">
        <p className="text-gray-600">
          欢迎使用 FustGo 数据同步平台！这是一个现代化的 ETL/ELT 解决方案。
        </p>
        <ul className="mt-4 space-y-2 text-gray-600">
          <li>• 点击左侧菜单开始使用</li>
          <li>• 在"任务管理"中创建数据同步任务</li>
          <li>• 在"连接配置"中管理数据源连接</li>
          <li>• 在"执行历史"中查看任务执行记录</li>
        </ul>
      </Card>
    </div>
  )
}
