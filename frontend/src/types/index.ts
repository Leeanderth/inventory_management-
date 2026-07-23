export type LoginValues = {
  username: string
  password: string
}

export type RoleInfo = {
  id: number
  code: string
  name: string
}

export type CurrentUser = {
  id: number
  username: string
  display_name: string
  role: RoleInfo
  permissions: string[]
}

export type LoginResponse = {
  token: string
  expires_at: number
  user: CurrentUser
}

export type RoleItem = {
  id: number
  code: string
  name: string
  description: string
  is_system: boolean
  permission_codes: string[]
}

export type PermissionItem = {
  id: number
  code: string
  name: string
  module: string
  description: string
}

export type RoleFormValues = {
  code: string
  name: string
  description?: string
  permission_codes?: string[]
}

export type UserItem = {
  id: number
  username: string
  display_name: string
  status: 'enabled' | 'disabled'
  role_id: number
  role: RoleInfo
}

export type UserFormValues = {
  username: string
  password?: string
  display_name?: string
  status: 'enabled' | 'disabled'
  role_id: number
}

export type ProductItem = {
  id: number
  name: string
  sku: string
  category: string
  quantity: number
  remark: string
  created_at: string
  updated_at: string
}

export type ProductFormValues = {
  name: string
  sku: string
  category?: string
  quantity: number
  remark?: string
}

export type StockMovementItem = {
  id: number
  before_quantity: number
  after_quantity: number
  change_quantity: number
  remark: string
  created_at: string
  operator?: {
    username: string
    display_name: string
  }
}
