import { useState } from 'react'
import { logout } from '../../lib/auth'
import { useAuth } from './AuthProvider'

export default function UserSessionCard({ user }) {
  const { clearSession } = useAuth()
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')

  async function handleLogout() {
    setStatus('loading')
    setError('')

    try {
      await logout()
      clearSession()
    } catch (err) {
      setStatus('error')
      setError(err.message || 'Could not log out')
    } finally {
      setStatus((current) => (current === 'loading' ? 'idle' : current))
    }
  }

  return (
    <section className="panel-card user-card" aria-labelledby="welcome-heading">
      <h2 id="welcome-heading">Welcome back</h2>
      <p className="user-email">{user.email}</p>
      {user.displayName ? <p className="user-name">{user.displayName}</p> : null}
      <button type="button" onClick={handleLogout} disabled={status === 'loading'}>
        {status === 'loading' ? 'Logging out…' : 'Log out'}
      </button>
      {error ? (
        <p className="status-message status-message-error" role="status">
          {error}
        </p>
      ) : null}
    </section>
  )
}
