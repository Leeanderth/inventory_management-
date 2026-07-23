import { PlusOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, Modal, Select, Space, Table, Tag, Typography, message } from 'antd'
import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { ajax } from '@/utils/request'
import type { RoleItem, UserFormValues, UserItem } from '@/types'

function UserManagementView() {
  const { user } = useAuth()
  const [users, setUsers] = useState<UserItem[]>([])
  const [roles, setRoles] = useState<RoleItem[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<UserItem | null>(null)
  const [form] = Form.useForm<UserFormValues>()

  const canCreate = Boolean(user?.permissions.includes('user:create'))
  const canUpdate = Boolean(user?.permissions.includes('user:update'))
  const canDisable = Boolean(user?.permissions.includes('user:disable'))
  const canViewRoles = Boolean(user?.permissions.includes('role:view'))

  const loadData = async () => {
    setLoading(true)
    try {
      const userRes = (await ajax({}, '/users')) as { items: UserItem[] }
      setUsers(userRes.items)

      if (canViewRoles && (canCreate || canUpdate)) {
        const roleRes = (await ajax({}, '/roles')) as { items: RoleItem[] }
        setRoles(roleRes.items)
      } else {
        setRoles([])
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const openCreateModal = () => {
    if (!canViewRoles) {
      message.warning('需要角色查看权限后才能分配角色')
      return
    }
    setEditingUser(null)
    form.resetFields()
    form.setFieldsValue({ status: 'enabled' })
    setModalOpen(true)
  }

  const openEditModal = (item: UserItem) => {
    if (!canViewRoles) {
      message.warning('需要角色查看权限后才能分配角色')
      return
    }
    setEditingUser(item)
    form.setFieldsValue({
      username: item.username,
      display_name: item.display_name,
      status: item.status,
      role_id: item.role_id,
    })
    setModalOpen(true)
  }

  const handleSave = async () => {
    const values = await form.validateFields()
    setSaving(true)
    try {
      const payload = {
        username: values.username,
        password: values.password || '',
        display_name: values.display_name || '',
        status: values.status,
        role_id: values.role_id,
      }
      if (editingUser) {
        await ajax(payload, `/users/${editingUser.id}`, 'put')
        message.success('用户已更新')
      } else {
        await ajax(payload, '/users', 'post')
        message.success('用户已创建')
      }
      setModalOpen(false)
      await loadData()
    } finally {
      setSaving(false)
    }
  }

  const toggleStatus = async (item: UserItem) => {
    const status = item.status === 'enabled' ? 'disabled' : 'enabled'
    await ajax({ status }, `/users/${item.id}/status`, 'patch')
    message.success(status === 'enabled' ? '用户已启用' : '用户已禁用')
    await loadData()
  }

  return (
    <section className="page-section">
      <div className="page-heading">
        <div>
          <Typography.Title level={2}>用户管理</Typography.Title>
          <Typography.Paragraph>维护系统账号，并为用户分配角色。</Typography.Paragraph>
        </div>
        {canCreate && (
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreateModal}>
            新建用户
          </Button>
        )}
      </div>

      <Card className="table-card" variant="borderless">
        <Table
          rowKey="id"
          loading={loading}
          dataSource={users}
          columns={[
            {
              title: '用户',
              dataIndex: 'username',
              render: (_, item) => (
                <Space direction="vertical" size={2}>
                  <Typography.Text strong>{item.display_name || item.username}</Typography.Text>
                  <Typography.Text type="secondary">{item.username}</Typography.Text>
                </Space>
              ),
            },
            { title: '角色', dataIndex: 'role', render: role => role?.name || '-' },
            { title: '状态', dataIndex: 'status', render: value => <Tag color={value === 'enabled' ? 'green' : 'red'}>{value === 'enabled' ? '启用' : '禁用'}</Tag> },
            {
              title: '操作',
              width: 180,
              render: (_, item) => (
                <Space>
                  {canUpdate && (
                    <Button type="link" onClick={() => openEditModal(item)}>
                      编辑
                    </Button>
                  )}
                  {canDisable && (
                    <Button type="link" danger={item.status === 'enabled'} onClick={() => toggleStatus(item)}>
                      {item.status === 'enabled' ? '禁用' : '启用'}
                    </Button>
                  )}
                </Space>
              ),
            },
          ]}
        />
      </Card>

      <Modal title={editingUser ? '编辑用户' : '新建用户'} open={modalOpen} onOk={handleSave} onCancel={() => setModalOpen(false)} confirmLoading={saving} destroyOnHidden>
        <Form form={form} layout="vertical" requiredMark={false}>
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input placeholder="username" />
          </Form.Item>
          <Form.Item name="password" label={editingUser ? '新密码' : '密码'} rules={editingUser ? [] : [{ required: true, message: '请输入密码' }]}>
            <Input.Password placeholder={editingUser ? '不填写则不修改' : '请输入密码'} />
          </Form.Item>
          <Form.Item name="display_name" label="显示名称">
            <Input placeholder="显示名称" />
          </Form.Item>
          <Form.Item name="role_id" label="角色" rules={[{ required: true, message: '请选择角色' }]}>
            <Select options={roles.map(role => ({ label: role.name, value: role.id }))} placeholder="请选择角色" />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true, message: '请选择状态' }]}>
            <Select
              options={[
                { label: '启用', value: 'enabled' },
                { label: '禁用', value: 'disabled' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </section>
  )
}

export default UserManagementView
