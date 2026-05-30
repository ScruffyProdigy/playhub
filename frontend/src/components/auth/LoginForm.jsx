import { useState } from 'react'
import { requestMagicLink } from '../../lib/auth'
import { useAuth } from './AuthProvider'

export default function LoginForm() {
  const { refreshSession } = useAuth()
  const [email, setEmail] = useState('')
  const [status, setStatus] = useState('idle')
  const [message, setMessage] = useState('')

  async function handleSubmit(event) {
    event.preventDefault()
    setStatus('loading')
    setMessage('')

    try {
      await requestMagicLink(email.trim())
      setStatus('sent')
      setMessage('Check your email for your sign-in link.')
    } catch (error) {
      setStatus('error')
      setMessage(error.message || 'Could not send magic link')
    }
  }

  return (
    <section className="auth-card" aria-labelledby="login-heading">
      <h2 id="login-heading">Sign in</h2>
      <p className="auth-copy">Enter your email and we&apos;ll send you a magic link.</p>

      <form className="auth-form" onSubmit={handleSubmit}>
        <label htmlFor="email">Email</label>
        <input
          id="email"
          name="email"
          type="email"
          autoComplete="email"
          required
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="you@example.com"
          disabled={status === 'loading'}
        />
        <button type="submit" disabled={status === 'loading' || email.trim() === ''}>
          {status === 'loading' ? 'Sending…' : 'Send magic link'}
        </button>
      </form>

      {message ? (
        <p className={status === 'error' ? 'auth-message auth-message-error' : 'auth-message'} role="status">
          {message}
        </p>
      ) : null}

      <button type="button" className="auth-link-button" onClick={refreshSession}>
        I already signed in
      </button>
    </section>
  )
}
