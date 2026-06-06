import { useState } from 'react'
import { logout } from '../../lib/auth'
import { useAuth } from './AuthProvider'
import AvatarPicker from '../avatars/AvatarPicker'
import PlayerAvatar from '../avatars/PlayerAvatar'

export default function UserSessionCard({ user }) {
  const { clearSession } = useAuth()
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')
  const [pickerOpen, setPickerOpen] = useState(!user.avatarKey)

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
      <div className="user-card__header">
        <PlayerAvatar user={user} size="md" />
        <div>
          <h2 id="welcome-heading">Welcome back</h2>
          <p className="user-email">{user.email}</p>
          {user.displayName ? <p className="user-name">{user.displayName}</p> : null}
        </div>
      </div>

      {pickerOpen ? (
        <AvatarPicker currentKey={user.avatarKey} onSelected={() => setPickerOpen(false)} />
      ) : (
        <button type="button" className="game-list-button game-list-button-secondary" onClick={() => setPickerOpen(true)}>
          Change avatar
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
