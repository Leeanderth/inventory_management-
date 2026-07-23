import { PlusOutlined } from '@ant-design/icons'
import { Button, Card, Checkbox, Form, Input, Modal, Popconfirm, Space, Table, Tag, Typography, message } from 'antd'
import { useEffect, useMemo, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { ajax } from '@/utils/request'
import type { PermissionItem, RoleFormValues, RoleItem } from '@/types'

const moduleNames: Record<string, string> = {
  role: '角色权限',
  stock: '库存',
  user: '用户',
}

function RolePermissionView() {
  const { user } = useAuth()
  const [roles, setRoles] = useState<RoleItem[]>([])
  const [permissions, setPermissions] = useState<PermissionItem[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingRole, setEditingRole] = useState<RoleItem | null>(null)
  const [form] = Form.useForm<RoleFormValues>()

  const canCreate = Boolean(user?.permissions.includes('role:create'))
  const canUpdate = Boolean(user?.permissions.includes('role:update'))
  const canDelete = Boolean(user?.permissions.includes('role:delete'))

  const permissionGroups = useMemo(() => {
    return permissions.reduce<Record<string, PermissionItem[]>>((groups, permission) => {
      if (!groups[permission.module]) groups[permission.module] = []
      groups[permission.module].push(permission)
      return groups
    }, {})
  }, [permissions])

  const loadData = async () => {
    setLoading(true)
    try {
      const [roleRes, permissionRes] = await Promise.all([
        ajax({}, '/roles') as Promise<{ items: RoleItem[] }>,
        ajax({}, '/permissions') as Promise<{ items: PermissionItem[] }>,
      ])
      setRoles(roleRes.items)
      setPermissions(permissionRes.items)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const openCreateModal = () => {
    setEditingRole(null)
    form.resetFields()
    setModalOpen(true)
  }

  const openEditModal = (role: RoleItem) => {
    setEditingRole(role)
    form.setFieldsValue({
      code: role.code,
      name: role.name,
      description: role.description,
      permission_codes: role.permission_codes || [],
    })
    setModalOpen(true)
  }

  const handleSave = async () => {
    const values = await form.validateFields()
    setSaving(true)
    try {
      const payload = {
        code: values.code,
        name: values.name,
        description: values.description || '',
        permission_codes: values.permission_codes || [],
      }
      if (editingRole) {
        await ajax(payload, `/roles/${editingRole.id}`, 'put')
        message.success('角色已更新')
      } else {
        await ajax(payload, '/roles', 'post')
        message.success('角色已创建')
      }
      setModalOpen(false)
      await loadData()
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (role: RoleItem) => {
    await ajax({}, `/roles/${role.id}`, 'delete')
    message.success('角色已删除')
    await loadData()
  }

  return (
    <section className="page-section">
      <div className="page-heading">
        <div>
          <Typography.Title level={2}>角色权限管理</Typography.Title>
          <Typography.Paragraph>管理角色基础信息，并为角色分配系统权限。</Typography.Paragraph>
        </div>
        {canCreate && (
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreateModal}>
            新建角色
          </Button>
        )}
      </div>

      <Card className="table-card" variant="borderless">
        <Table
          rowKey="id"
          loading={loading}
          dataSource={roles}
          columns={[
            {
              title: '角色名称',
              dataIndex: 'name',
              width: 180,
              render: (_, role) => (
                <Space direction="vertical" size={2}>
                  <Typography.Text strong>{role.name}</Typography.Text>
                  <Typography.Text type="secondary">{role.code}</Typography.Text>
                </Space>
              ),
            },
            { title: '类型', dataIndex: 'is_system', width: 120, render: value => (value ? <Tag color="blue">系统内置</Tag> : <Tag>自定义</Tag>) },
            { title: '说明', dataIndex: 'description', render: value => value || '-' },
            {
              title: '权限',
              dataIndex: 'permission_codes',
              render: (codes: string[] | null) => {
                const safeCodes = codes || []
                return (
                <div className="permission-tags">
                  {safeCodes.length ? safeCodes.map(code => <Tag key={code}>{code}</Tag>) : <Typography.Text type="secondary">暂无</Typography.Text>}
                </div>
                )
              },
            },
            {
              title: '操作',
              width: 160,
              render: (_, role) => (
                <Space>
                  {canUpdate && (
                    <Button type="link" onClick={() => openEditModal(role)}>
                      编辑
                    </Button>
                  )}
                  {canDelete && !role.is_system && (
                    <Popconfirm title="确认删除该角色？" onConfirm={() => handleDelete(role)}>
                      <Button type="link" danger>
                        删除
                      </Button>
                    </Popconfirm>
                  )}
                </Space>
              ),
            },
          ]}
        />
      </Card>

      <Modal
        title={editingRole ? '编辑角色' : '新建角色'}
        open={modalOpen}
        onOk={handleSave}
        onCancel={() => setModalOpen(false)}
        confirmLoading={saving}
        destroyOnHidden
        width={720}
      >
        <Form form={form} layout="vertical" requiredMark={false}>
          <Form.Item name="code" label="角色标识" rules={[{ required: true, message: '请输入角色标识' }]}>
            <Input placeholder="stock_operator" disabled={editingRole?.is_system} />
          </Form.Item>
          <Form.Item name="name" label="角色名称" rules={[{ required: true, message: '请输入角色名称' }]}>
            <Input placeholder="库存操作员" />
          </Form.Item>
          <Form.Item name="description" label="角色说明">
            <Input.TextArea rows={3} placeholder="描述该角色的职责范围" />
          </Form.Item>
          <Form.Item name="permission_codes" label="权限配置">
            <Checkbox.Group className="permission-checkbox-group">
              {Object.entries(permissionGroups).map(([module, items]) => (
                <div className="permission-module" key={module}>
                  <Typography.Text strong>{moduleNames[module] || module}</Typography.Text>
                  <div className="permission-options">
                    {items.map(permission => (
                      <Checkbox key={permission.code} value={permission.code}>
                        {permission.name}
                        <Typography.Text type="secondary"> {permission.code}</Typography.Text>
                      </Checkbox>
                    ))}
                  </div>
                </div>
              ))}
            </Checkbox.Group>
          </Form.Item>
        </Form>
      </Modal>
    </section>
  )
}

export default RolePermissionView
