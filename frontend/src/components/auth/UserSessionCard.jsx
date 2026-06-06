import { useEffect, useState } from 'react'
import { logout } from '../../lib/auth'
import { needsProfileSetup } from '../../lib/avatars'
import { useAuth } from './AuthProvider'
import PlayerProfileEditor from '../avatars/PlayerProfileEditor'
import PlayerAvatar from '../avatars/PlayerAvatar'

export default function UserSessionCard({ user }) {
  const { clearSession } = useAuth()
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')
  const setupRequired = needsProfileSetup(user)
  const [editorOpen, setEditorOpen] = useState(setupRequired)

  useEffect(() => {
    if (needsProfileSetup(user)) {
      setEditorOpen(true)
    }
  }, [user])

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

  function handleSaved(updated) {
    if (!needsProfileSetup(updated)) {
      setEditorOpen(false)
    }
  }

  return (
    <section className="panel-card user-card" aria-labelledby="welcome-heading">
      <div className="user-card__header">
        <PlayerAvatar user={user} size="md" />
        <div>
          <h2 id="welcome-heading">{setupRequired ? 'Set up your display' : 'Welcome back'}</h2>
          <p className="user-email">{user.email}</p>
          {!setupRequired && user.displayName ? <p className="user-name">{user.displayName}</p> : null}
        </div>
      </div>

      {editorOpen ? (
        <PlayerProfileEditor
          user={user}
          required={setupRequired}
          onSaved={handleSaved}
          onCancel={setupRequired ? undefined : () => setEditorOpen(false)}
        />
      ) : (
        <button
          type="button"
          className="game-list-button game-list-button-secondary"
          onClick={() => setEditorOpen(true)}
        >
          Change display
        </button>
      )}

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
