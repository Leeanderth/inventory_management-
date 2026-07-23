const PASSPORT_KEY = 'passport'

type PassportPayload = {
  token?: string
  expiresAt?: number
}

function readPassport(): PassportPayload | null {
  const raw = localStorage.getItem(PASSPORT_KEY)
  if (!raw) return null

  try {
    return JSON.parse(raw) as PassportPayload
  } catch {
    return { token: raw }
  }
}

export const passport = {
  isValid() {
    const data = readPassport()
    if (!data?.token) return false
    if (!data.expiresAt) return true
    return data.expiresAt > Date.now()
  },

  getToken() {
    return readPassport()?.token || ''
  },

  setToken(token: string, expiresAt?: number) {
    localStorage.setItem(PASSPORT_KEY, JSON.stringify({ token, expiresAt }))
  },

  clear() {
    localStorage.removeItem(PASSPORT_KEY)
  },
}
