import LoginForm from './LoginForm'
import UserSessionCard from './UserSessionCard'
import { useAuth } from './AuthProvider'

function SessionLoading() {
  return (
    <p className="status-message" role="status">
      Loading session…
    </p>
  )
}

export default function AuthPanel() {
  const { user, loading, error, sessionUnavailable, refreshSession } = useAuth()

  if (loading) {
    return <SessionLoading />
  }

  if (user) {
    return <UserSessionCard user={user} />
  }

  return (
    <>
      {error ? <p className="status-message status-message-error">{error}</p> : null}
      {sessionUnavailable ? (
        <button type="button" className="button-secondary" onClick={() => void refreshSession()}>
          Try again
        </button>
      ) : (
        <LoginForm />
      )}
    </>
  )
}
