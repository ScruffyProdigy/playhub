import LoginForm from './LoginForm'
import UserSessionCard from './UserSessionCard'
import { useAuth } from './AuthProvider'

function SessionLoading() {
  return (
    <p className="auth-message" role="status">
      Loading session…
    </p>
  )
}

export default function AuthPanel() {
  const { user, loading, error } = useAuth()

  if (loading) {
    return <SessionLoading />
  }

  if (user) {
    return <UserSessionCard user={user} />
  }

  return (
    <>
      {error ? <p className="auth-message auth-message-error">{error}</p> : null}
      <LoginForm />
    </>
  )
}
