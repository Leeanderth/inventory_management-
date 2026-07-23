import { PlusOutlined } from '@ant-design/icons'
import { Button, Card, Drawer, Form, Input, InputNumber, Modal, Popconfirm, Space, Table, Tag, Typography, message } from 'antd'
import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { ajax } from '@/utils/request'
import type { ProductFormValues, ProductItem, StockMovementItem } from '@/types'

function ProductView() {
  const { user } = useAuth()
  const [products, setProducts] = useState<ProductItem[]>([])
  const [movements, setMovements] = useState<StockMovementItem[]>([])
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [editingProduct, setEditingProduct] = useState<ProductItem | null>(null)
  const [form] = Form.useForm<ProductFormValues>()

  const canCreate = Boolean(user?.permissions.includes('stock:create'))
  const canUpdate = Boolean(user?.permissions.includes('stock:update'))
  const canDelete = Boolean(user?.permissions.includes('stock:delete'))

  const loadProducts = async (nextKeyword = keyword) => {
    setLoading(true)
    try {
      const res = (await ajax({ keyword: nextKeyword }, '/products')) as { items: ProductItem[] }
      setProducts(res.items)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProducts('')
  }, [])

  const openCreateModal = () => {
    setEditingProduct(null)
    form.resetFields()
    form.setFieldsValue({ quantity: 0 })
    setModalOpen(true)
  }

  const openEditModal = (product: ProductItem) => {
    setEditingProduct(product)
    form.setFieldsValue(product)
    setModalOpen(true)
  }

  const handleSave = async () => {
    const values = await form.validateFields()
    setSaving(true)
    try {
      const payload = {
        name: values.name,
        sku: values.sku,
        category: values.category || '',
        quantity: values.quantity || 0,
        remark: values.remark || '',
      }
      if (editingProduct) {
        await ajax(payload, `/products/${editingProduct.id}`, 'put')
        message.success('商品已更新')
      } else {
        await ajax(payload, '/products', 'post')
        message.success('商品已创建')
      }
      setModalOpen(false)
      await loadProducts()
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (product: ProductItem) => {
    await ajax({}, `/products/${product.id}`, 'delete')
    message.success('商品已删除')
    await loadProducts()
  }

  const openMovements = async (product: ProductItem) => {
    setEditingProduct(product)
    setDrawerOpen(true)
    const res = (await ajax({}, `/products/${product.id}/movements`)) as { items: StockMovementItem[] }
    setMovements(res.items)
  }

  return (
    <section className="page-section">
      <div className="page-heading">
        <div>
          <Typography.Title level={2}>库存管理</Typography.Title>
          <Typography.Paragraph>维护商品基础信息，库存数量变化会自动记录变动。</Typography.Paragraph>
        </div>
        {canCreate && (
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreateModal}>
            新增商品
          </Button>
        )}
      </div>

      <Card className="table-card" variant="borderless">
        <div className="toolbar">
          <Input.Search
            allowClear
            placeholder="搜索商品名称或 SKU"
            value={keyword}
            onChange={event => setKeyword(event.target.value)}
            onSearch={value => loadProducts(value)}
            style={{ maxWidth: 320 }}
          />
        </div>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={products}
          columns={[
            { title: '商品名称', dataIndex: 'name' },
            { title: 'SKU', dataIndex: 'sku', width: 150 },
            { title: '分类', dataIndex: 'category', width: 140, render: value => value || '-' },
            { title: '库存数量', dataIndex: 'quantity', width: 120, render: value => <Tag color={value > 0 ? 'green' : 'default'}>{value}</Tag> },
            { title: '备注', dataIndex: 'remark', render: value => value || '-' },
            {
              title: '操作',
              width: 220,
              render: (_, product) => (
                <Space>
                  <Button type="link" onClick={() => openMovements(product)}>
                    变动记录
                  </Button>
                  {canUpdate && (
                    <Button type="link" onClick={() => openEditModal(product)}>
                      编辑
                    </Button>
                  )}
                  {canDelete && (
                    <Popconfirm title="确认删除该商品？" onConfirm={() => handleDelete(product)}>
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
        title={editingProduct ? '编辑商品' : '新增商品'}
        open={modalOpen}
        onOk={handleSave}
        onCancel={() => setModalOpen(false)}
        confirmLoading={saving}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" requiredMark={false}>
          <Form.Item name="name" label="商品名称" rules={[{ required: true, message: '请输入商品名称' }]}>
            <Input placeholder="商品名称" />
          </Form.Item>
          <Form.Item name="sku" label="SKU" rules={[{ required: true, message: '请输入 SKU' }]}>
            <Input placeholder="唯一商品编码" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Input placeholder="分类" />
          </Form.Item>
          <Form.Item name="quantity" label="库存数量" rules={[{ required: true, message: '请输入库存数量' }]}>
            <InputNumber min={0} precision={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={3} placeholder="备注" />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer title={`${editingProduct?.name || ''} 库存变动记录`} open={drawerOpen} onClose={() => setDrawerOpen(false)} width={620}>
        <Table
          rowKey="id"
          dataSource={movements}
          pagination={false}
          columns={[
            { title: '变动前', dataIndex: 'before_quantity' },
            { title: '变动后', dataIndex: 'after_quantity' },
            { title: '变化', dataIndex: 'change_quantity', render: value => <Tag color={value >= 0 ? 'green' : 'red'}>{value}</Tag> },
            { title: '操作人', dataIndex: 'operator', render: operator => operator?.display_name || operator?.username || '-' },
            { title: '备注', dataIndex: 'remark', render: value => value || '-' },
          ]}
        />
      </Drawer>
    </section>
  )
}

export default ProductView
