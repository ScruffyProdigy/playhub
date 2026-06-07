import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { fetchCurrentUser } from '../../lib/auth'
import { isTransientServerError } from '../../lib/graphql'
import { clearSubscriptionAuthCache, prefetchSubscriptionAuth } from '../../lib/queue'

const SESSION_RETRY_ATTEMPTS = 4
const SESSION_RETRY_DELAY_MS = 750

async function fetchCurrentUserWithRetries() {
  let lastErr
  for (let attempt = 0; attempt < SESSION_RETRY_ATTEMPTS; attempt += 1) {
    try {
      return await fetchCurrentUser()
    } catch (err) {
      lastErr = err
      if (!isTransientServerError(err.message) || attempt === SESSION_RETRY_ATTEMPTS - 1) {
        throw err
      }
      await new Promise((resolve) => setTimeout(resolve, SESSION_RETRY_DELAY_MS * (attempt + 1)))
    }
  }
  throw lastErr
}

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [sessionUnavailable, setSessionUnavailable] = useState(false)

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
    setSessionUnavailable(false)
    try {
      const currentUser = await fetchCurrentUserWithRetries()
      if (currentUser) {
        setUser(currentUser)
        void prefetchSubscriptionAuth().catch(() => {})
      } else if (!silent) {
        setUser(null)
      }
    } catch (err) {
      if (!silent) {
        const message = err.message || 'Could not load session'
        if (isTransientServerError(message)) {
          setSessionUnavailable(true)
          setError('Server briefly unavailable — your session may still be active. Try again in a moment.')
        } else {
          setError(message)
          setUser(null)
        }
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
    setSessionUnavailable(false)
  }, [])

  useEffect(() => {
    refreshSession()
  }, [refreshSession])

  const value = useMemo(
    () => ({
      user,
      loading,
      error,
      sessionUnavailable,
      refreshSession,
      acceptSessionUser,
      clearSession,
    }),
    [user, loading, error, sessionUnavailable, refreshSession, acceptSessionUser, clearSession],
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
