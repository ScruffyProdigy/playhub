import { useEffect, useState } from 'react'
import { oauthErrorMessage } from '../../lib/oauth'

export default function OAuthCompletePage() {
  const [message, setMessage] = useState('Something went wrong during sign-in.')

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const error = params.get('error')
    if (error) {
      setMessage(oauthErrorMessage(error))
      window.history.replaceState({}, '', '/auth/oauth/complete')
    }
  }, [])

  return (
    <main className="app-shell app-shell--narrow">
      <section className="panel-card auth-sign-in">
        <h1>Sign-in</h1>
        <p className="panel-copy status-message status-message-error" role="alert">
          {message}
        </p>
        <a className="auth-link" href="/">
          Back to home
        </a>
      </section>
    </main>
  )
}
