import { useEffect, useState } from 'react'
import { completeMagicLoginOnce } from '../../lib/auth'

function getTokenFromLocation() {
  const params = new URLSearchParams(window.location.search)
  return params.get('token')?.trim() || ''
}

export default function CompleteMagicPage() {
  const [status, setStatus] = useState('loading')
  const [message, setMessage] = useState('Completing sign-in…')

  useEffect(() => {
    const token = getTokenFromLocation()
    if (!token) {
      setStatus('error')
      setMessage('Missing sign-in token. Request a new magic link.')
      return
    }

    let cancelled = false

    completeMagicLoginOnce(token)
      .then(() => {
        if (cancelled) {
          return
        }
        setStatus('success')
        setMessage('Signed in. Redirecting…')
        window.history.replaceState({}, '', '/')
        window.location.assign('/')
      })
      .catch((error) => {
        if (cancelled) {
          return
        }
        setStatus('error')
        setMessage(error.message || 'Could not complete sign-in')
      })

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <main className="auth-page">
      <h1>PlayHub</h1>
      <section className="auth-card" aria-live="polite">
        <h2>Signing you in</h2>
        <p className={status === 'error' ? 'auth-message auth-message-error' : 'auth-message'}>{message}</p>
        {status === 'error' ? (
          <a className="auth-link" href="/">
            Back to home
          </a>
        ) : null}
      </section>
    </main>
  )
}
