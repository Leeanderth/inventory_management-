import { Navigate } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'
import { getFirstAccessiblePath } from '@/router/permissions'

type RequirePermissionRouteProps = {
  permission: string
  children: React.ReactNode
}

function RequirePermissionRoute({ permission, children }: RequirePermissionRouteProps) {
  const { user } = useAuth()

  if (!user) return <Navigate to="/login" replace />
  if (!user.permissions.includes(permission)) {
    return <Navigate to={getFirstAccessiblePath(user)} replace />
  }

  return children
}

export default RequirePermissionRoute
