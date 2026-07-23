import { createContext, useContext } from 'react'
import type { CurrentUser } from '@/types'

type AuthContextValue = {
  user: CurrentUser | null
  setUser: (user: CurrentUser | null) => void
  logout: () => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthContext.Provider')
  return context
}
