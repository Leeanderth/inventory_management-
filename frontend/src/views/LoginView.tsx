import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, Typography, message } from 'antd'
import { Navigate, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { ajax } from '@/utils/request'
import { passport } from '@/utils/passport'
import type { LoginResponse, LoginValues } from '@/types'
import { getFirstAccessiblePath } from '@/router/permissions'

const defaultAccounts = [
  ['root', '超级管理员'],
  ['manager', '库存管理员'],
  ['viewer', '库存查看员'],
]

function LoginView() {
  const [loading, setLoading] = useState(false)
  const { user, setUser } = useAuth()
  const navigate = useNavigate()

  if (user) return <Navigate to={getFirstAccessiblePath(user)} replace />

  const handleLogin = async (values: LoginValues) => {
    setLoading(true)
    try {
      const res = (await ajax(values, '/auth/login', 'post')) as LoginResponse
      passport.setToken(res.token, res.expires_at)
      setUser(res.user)
      message.success('登录成功')
      navigate(getFirstAccessiblePath(res.user), { replace: true })
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="login-page">
      <section className="login-panel">
        <div className="login-copy">
          <Typography.Text className="eyebrow">Inventory Management</Typography.Text>
          <Typography.Title level={1}>商品库存管理</Typography.Title>
          <Typography.Paragraph>
            使用内置账号登录系统，后续库存、用户和角色权限都会围绕这套身份体系展开。
          </Typography.Paragraph>
          <div className="account-list">
            {defaultAccounts.map(([username, label]) => (
              <span key={username}>
                {username} · {label}
              </span>
            ))}
          </div>
        </div>

        <Card className="login-card" variant="borderless">
          <Typography.Title level={2}>登录</Typography.Title>
          <Form layout="vertical" requiredMark={false} onFinish={handleLogin}>
            <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}>
              <Input size="large" prefix={<UserOutlined />} placeholder="root" autoComplete="username" />
            </Form.Item>
            <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}>
              <Input.Password size="large" prefix={<LockOutlined />} placeholder="请输入密码" autoComplete="current-password" />
            </Form.Item>
            <Button type="primary" htmlType="submit" size="large" block loading={loading}>
              登录系统
            </Button>
          </Form>
        </Card>
      </section>
    </main>
  )
}

export default LoginView
