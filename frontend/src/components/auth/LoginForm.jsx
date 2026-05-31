import { useCallback, useState } from 'react'
import { completeSignInWithCode, requestSignIn } from '../../lib/auth'
import { useAuth } from './AuthProvider'
import useWaitForSignIn from './useWaitForSignIn'

function normalizeCode(value) {
  return value.replace(/\D/g, '').slice(0, 6)
}

export default function LoginForm() {
  const { refreshSession, user } = useAuth()
  const [email, setEmail] = useState('')
  const [code, setCode] = useState('')
  const [step, setStep] = useState('email')
  const [status, setStatus] = useState('idle')
  const [message, setMessage] = useState('')

  const handleSignedInElsewhere = useCallback(() => {
    void refreshSession({ silent: true })
  }, [refreshSession])

  useWaitForSignIn({
    enabled: step === 'verify' && !user,
    onSignedIn: handleSignedInElsewhere,
  })

  async function handleEmailSubmit(event) {
    event.preventDefault()
    setStatus('loading')
    setMessage('')

    try {
      await requestSignIn(email.trim())
      setStep('verify')
      setCode('')
      setStatus('idle')
      setMessage('We sent a 6-digit code and sign-in link to your email.')
    } catch (error) {
      setStatus('error')
      setMessage(error.message || 'Could not send sign-in email')
    }
  }

  async function handleCodeSubmit(event) {
    event.preventDefault()
    const normalizedCode = normalizeCode(code)
    if (normalizedCode.length !== 6) {
      setStatus('error')
      setMessage('Enter the 6-digit code from your email.')
      return
    }

    setStatus('loading')
    setMessage('')

    try {
      await completeSignInWithCode(email.trim(), normalizedCode)
      await refreshSession({ silent: true })
    } catch (error) {
      setStatus('error')
      setMessage(error.message || 'Invalid or expired code')
    }
  }

  function handleUseDifferentEmail() {
    setStep('email')
    setCode('')
    setStatus('idle')
    setMessage('')
  }

  if (step === 'verify') {
    return (
      <section className="panel-card" aria-labelledby="verify-heading">
        <h2 id="verify-heading">Enter your code</h2>
        <p className="panel-copy">
          Check <strong>{email}</strong> for a 6-digit code. On iPhone, you can paste the code from Mail. You can also tap
          the sign-in link in that email.
        </p>

        {message ? (
          <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'} role="status">
            {message}
          </p>
        ) : null}

        <form className="auth-form" onSubmit={handleCodeSubmit}>
          <label htmlFor="login-code">Sign-in code</label>
          <input
            id="login-code"
            name="code"
            type="text"
            inputMode="numeric"
            autoComplete="one-time-code"
            pattern="\d{6}"
            maxLength={6}
            required
            value={code}
            onChange={(event) => setCode(normalizeCode(event.target.value))}
            placeholder="123456"
            disabled={status === 'loading'}
            className="auth-code-input"
          />
          <button type="submit" disabled={status === 'loading' || normalizeCode(code).length !== 6}>
            {status === 'loading' ? 'Signing in…' : 'Continue'}
          </button>
        </form>

        <button type="button" className="auth-link-button" onClick={handleUseDifferentEmail} disabled={status === 'loading'}>
          Use a different email
        </button>
      </section>
    )
  }

  return (
    <section className="panel-card" aria-labelledby="login-heading">
      <h2 id="login-heading">Sign in</h2>
      <p className="panel-copy">Enter your email and we&apos;ll send a 6-digit code and sign-in link.</p>

      <form className="auth-form" onSubmit={handleEmailSubmit}>
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
          {status === 'loading' ? 'Sending…' : 'Continue'}
        </button>
      </form>

      {message ? (
        <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'} role="status">
          {message}
        </p>
      ) : null}
    </section>
  )
}
