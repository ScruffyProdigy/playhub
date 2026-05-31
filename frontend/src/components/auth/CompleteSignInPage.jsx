import { useEffect, useState } from 'react'
import { completeSignInWithLinkOnce } from '../../lib/auth'
import { notifyAuthComplete } from '../../lib/authBroadcast'
import { APP_NAME } from '../../lib/brand'

function getTokenFromLocation() {
  const params = new URLSearchParams(window.location.search)
  return params.get('token')?.trim() || ''
}

export default function CompleteSignInPage() {
  const [status, setStatus] = useState('loading')
  const [message, setMessage] = useState('Completing sign-in…')

  useEffect(() => {
    const token = getTokenFromLocation()
    if (!token) {
      setStatus('error')
      setMessage('Missing sign-in token. Request a new sign-in email.')
      return
    }

    let cancelled = false

    completeSignInWithLinkOnce(token)
      .then(() => {
        if (cancelled) {
          return
        }
        notifyAuthComplete()
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
      <h1>{APP_NAME}</h1>
      <section className="panel-card" aria-live="polite">
        <h2>Signing you in</h2>
        <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'}>{message}</p>
        {status === 'error' ? (
          <a className="auth-link" href="/">
            Back to home
          </a>
        ) : null}
      </section>
    </main>
  )
}
