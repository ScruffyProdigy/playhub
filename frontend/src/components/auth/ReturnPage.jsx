import { useEffect, useState } from 'react'
import { fetchReturnDestination } from '../../lib/return'
import { APP_NAME } from '../../lib/brand'

function matchIdFromLocation() {
  const params = new URLSearchParams(window.location.search)
  return params.get('match')?.trim() || ''
}

export default function ReturnPage() {
  const [status, setStatus] = useState('loading')
  const [message, setMessage] = useState('Taking you back…')

  useEffect(() => {
    let cancelled = false
    const matchId = matchIdFromLocation()

    fetchReturnDestination(matchId || null)
      .then((dest) => {
        if (cancelled) return
        const path = dest?.path?.trim() || '/'
        const safePath = path.startsWith('/') && !path.startsWith('//') ? path : '/'
        setStatus('success')
        setMessage('Redirecting…')
        window.history.replaceState({}, '', safePath)
        window.location.assign(safePath)
      })
      .catch((error) => {
        if (cancelled) return
        setStatus('error')
        setMessage(error.message || 'Could not resolve return destination')
      })

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <main className="auth-page">
      <h1>{APP_NAME}</h1>
      <section className="panel-card" aria-live="polite">
        <h2>Welcome back</h2>
        <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'}>
          {message}
        </p>
        {status === 'error' ? (
          <a className="auth-link" href="/">
            Continue to JoinQuest
          </a>
        ) : null}
      </section>
    </main>
  )
}
