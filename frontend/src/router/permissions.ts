import type { CurrentUser } from '@/types'

export const permissionRoutes = [
  { path: '/stock', permission: 'stock:view' },
  { path: '/users', permission: 'user:view' },
  { path: '/roles', permission: 'role:view' },
]

export function getFirstAccessiblePath(user: CurrentUser | null) {
  if (!user) return '/login'

  const route = permissionRoutes.find(item => user.permissions.includes(item.permission))
  return route?.path || '/login'
}
