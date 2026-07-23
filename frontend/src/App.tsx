import { useEffect, useMemo, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { message } from 'antd'
import './App.css'
import { AuthContext } from '@/contexts/AuthContext'
import { ajax } from '@/utils/request'
import { passport } from '@/utils/passport'
import type { CurrentUser } from '@/types'
import { getFirstAccessiblePath } from '@/router/permissions'

function App() {
  const [user, setUser] = useState<CurrentUser | null>(null)
  const [checking, setChecking] = useState(true)
  const navigate = useNavigate()
  const location = useLocation()

  useEffect(() => {
    if (!passport.isValid()) {
      setChecking(false)
      if (location.pathname !== '/login') navigate('/login', { replace: true })
      return
    }

    ajax({}, '/auth/me')
      .then((res: { user: CurrentUser }) => {
        setUser(res.user)
        if (location.pathname === '/login' || location.pathname === '/') {
          navigate(getFirstAccessiblePath(res.user), { replace: true })
        }
      })
      .catch(() => {
        passport.clear()
        if (location.pathname !== '/login') navigate('/login', { replace: true })
      })
      .finally(() => setChecking(false))
  }, [location.pathname, navigate])

  const contextValue = useMemo(
    () => ({
      user,
      setUser,
      logout: () => {
        passport.clear()
        setUser(null)
        message.success('已退出登录')
        navigate('/login', { replace: true })
      },
    }),
    [navigate, user],
  )

  if (checking) return <div className="route-loading">加载中...</div>

  return (
    <AuthContext.Provider value={contextValue}>
      <Outlet />
    </AuthContext.Provider>
  )
}

export default App
