import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { fetchCurrentUser } from '../../lib/auth'
import { clearSubscriptionAuthCache, prefetchSubscriptionAuth } from '../../lib/queue'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const acceptSessionUser = useCallback((signedInUser) => {
    if (!signedInUser) {
      return
    }
    clearSubscriptionAuthCache()
    setUser(signedInUser)
    setError('')
    void prefetchSubscriptionAuth().catch(() => {})
  }, [])

  const refreshSession = useCallback(async (options = {}) => {
    const { silent = false } = options
    clearSubscriptionAuthCache()
    if (!silent) {
      setLoading(true)
    }
    setError('')
    try {
      const currentUser = await fetchCurrentUser()
      if (currentUser) {
        setUser(currentUser)
        void prefetchSubscriptionAuth().catch(() => {})
      } else if (!silent) {
        setUser(null)
      }
    } catch (err) {
      if (!silent) {
        setError(err.message || 'Could not load session')
        setUser(null)
      }
    } finally {
      if (!silent) {
        setLoading(false)
      }
    }
  }, [])

  const clearSession = useCallback(() => {
    clearSubscriptionAuthCache()
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
      acceptSessionUser,
      clearSession,
    }),
    [user, loading, error, refreshSession, acceptSessionUser, clearSession],
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
