import { useCallback, useEffect, useState } from 'react'
import CompleteMagicPage from './components/CompleteMagicPage'
import LoginForm from './components/LoginForm'
import { fetchCurrentUser } from './lib/auth'
import './App.css'

function HomePage() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadUser = useCallback(async () => {
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

  useEffect(() => {
    loadUser()
  }, [loadUser])

  return (
    <main className="app-shell">
      <header className="app-header">
        <h1>PlayHub</h1>
        <p className="tagline">Your Gaming Hub - Queue, Play, Trade</p>
      </header>

      {loading ? (
        <p className="auth-message" role="status">
          Loading session…
        </p>
      ) : user ? (
        <section className="auth-card user-card" aria-labelledby="welcome-heading">
          <h2 id="welcome-heading">Welcome back</h2>
          <p className="user-email">{user.email}</p>
          {user.displayName ? <p className="user-name">{user.displayName}</p> : null}
          <button type="button" onClick={loadUser}>
            Refresh session
          </button>
        </section>
      ) : (
        <>
          {error ? <p className="auth-message auth-message-error">{error}</p> : null}
          <LoginForm onSignedIn={loadUser} />
        </>
      )}
    </main>
  )
}

function App() {
  if (window.location.pathname.startsWith('/auth/complete')) {
    return <CompleteMagicPage />
  }

  return <HomePage />
}

export default App
