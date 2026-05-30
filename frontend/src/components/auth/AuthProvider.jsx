import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { fetchCurrentUser } from '../../lib/auth'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const refreshSession = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const currentUser = await fetchCurrentUser()
      setUser(currentUser)
    } catch (err) {
      setError(err.message || 'Could not load session')
      setUser(null)
    } finally {
      setLoading(false)
    }
  }, [])

  const clearSession = useCallback(() => {
    setUser(null)
    setError('')
  }, [])

  useEffect(() => {
    refreshSession()
  }, [refreshSession])

  const value = useMemo(
    () => ({
      user,
      loading,
      error,
      refreshSession,
      clearSession,
    }),
    [user, loading, error, refreshSession, clearSession],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
