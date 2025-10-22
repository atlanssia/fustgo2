import { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, theme } from 'antd'
import {
  DashboardOutlined,
  DatabaseOutlined,
  HistoryOutlined,
  LinkOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SettingOutlined,
  ApiOutlined,
  PartitionOutlined,
} from '@ant-design/icons'

const { Header, Sider, Content } = Layout

export default function MainLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const {
    token: { colorBgContainer },
  } = theme.useToken()

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: '/pipelines',
      icon: <PartitionOutlined />,
      label: '数据管道',
    },
    {
      key: '/jobs',
      icon: <DatabaseOutlined />,
      label: '任务管理',
    },
    {
      key: '/executions',
      icon: <HistoryOutlined />,
      label: '执行历史',
    },
    {
      key: '/connections',
      icon: <LinkOutlined />,
      label: '连接配置',
    },
    {
      key: '/plugins',
      icon: <ApiOutlined />,
      label: '插件管理',
    },
    {
      key: '/config',
      icon: <SettingOutlined />,
      label: '拖拽配置',
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <div className="h-16 flex items-center justify-center text-white text-xl font-bold">
          {collapsed ? 'FG' : 'FustGo'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ padding: 0, background: colorBgContainer }}>
          <div className="flex items-center justify-between px-4">
            <button
              onClick={() => setCollapsed(!collapsed)}
              className="text-lg w-16 h-16 hover:bg-gray-100"
            >
              {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            </button>
            <div className="text-gray-600">现代化数据同步平台</div>
          </div>
        </Header>
        <Content className="m-6 p-6 min-h-[280px]" style={{ background: colorBgContainer }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
