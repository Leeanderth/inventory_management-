import { Navigate } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'
import { getFirstAccessiblePath } from '@/router/permissions'

function HomeRedirect() {
  const { user } = useAuth()

  return <Navigate to={getFirstAccessiblePath(user)} replace />
}

export default HomeRedirect
