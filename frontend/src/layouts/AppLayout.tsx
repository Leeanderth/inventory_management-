import { InboxOutlined, LogoutOutlined, SafetyCertificateOutlined, TeamOutlined } from '@ant-design/icons'
import { Button, Layout, Menu, Typography } from 'antd'
import { Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'

function AppLayout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  if (!user) return <Navigate to="/login" replace />

  const hasPermission = (code: string) => user.permissions.includes(code)
  const menuItems = [
    ...(hasPermission('stock:view') ? [{ key: '/stock', icon: <InboxOutlined />, label: '库存管理' }] : []),
    ...(hasPermission('user:view') ? [{ key: '/users', icon: <TeamOutlined />, label: '用户管理' }] : []),
    ...(hasPermission('role:view') ? [{ key: '/roles', icon: <SafetyCertificateOutlined />, label: '角色权限管理' }] : []),
  ]

  if (menuItems.length === 0) return <div className="route-loading">当前用户暂无可访问页面</div>

  return (
    <Layout className="dashboard-shell">
      <Layout.Sider width={232} className="dashboard-sider">
        <div className="brand-block">
          <span className="brand-mark">IM</span>
          <div>
            <strong>库存管理</strong>
            <small>Inventory</small>
          </div>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={item => navigate(item.key)}
        />
      </Layout.Sider>

      <Layout>
        <Layout.Header className="dashboard-header">
          <div>
            <Typography.Text strong>{user.display_name || user.username}</Typography.Text>
            <Typography.Text type="secondary"> · {user.role.name}</Typography.Text>
          </div>
          <Button icon={<LogoutOutlined />} onClick={logout}>
            退出
          </Button>
        </Layout.Header>

        <Layout.Content className="dashboard-content">
          <Outlet />
        </Layout.Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
