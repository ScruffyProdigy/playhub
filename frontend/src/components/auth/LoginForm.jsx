import { useCallback, useLayoutEffect, useRef, useState } from 'react'
import { flushSync } from 'react-dom'
import { completeSignInWithCode, requestSignIn } from '../../lib/auth'
import { useAuth } from './AuthProvider'
import useWaitForSignIn from './useWaitForSignIn'
import { focusCodeInput } from './focusCodeInput'

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
  const codeInputRef = useRef(null)
  const emailRequestStarted = useRef(false)
  const onVerifyStep = step === 'verify'
  const isSigningIn = status === 'loading' && normalizeCode(code).length === 6

  const handleSignedInElsewhere = useCallback(() => {
    void refreshSession({ silent: true })
  }, [refreshSession])

  useWaitForSignIn({
    enabled: onVerifyStep && !user,
    onSignedIn: handleSignedInElsewhere,
  })

  function showVerifyStep() {
    flushSync(() => {
      setStep('verify')
      setCode('')
      setStatus('loading')
      setMessage('Sending sign-in email…')
    })
  }

  function focusCodeField() {
    focusCodeInput(codeInputRef.current)
  }

  /** Show the code step and focus the visible 123456 field (enabled so iOS autofill can attach). */
  function prepareCodeEntry() {
    if (step === 'verify') return
    showVerifyStep()
    focusCodeField()
  }

  // Refocus after the sign-in email is sent so iOS can offer one-time-code autofill.
  useLayoutEffect(() => {
    if (onVerifyStep && status === 'idle') {
      focusCodeField()
    }
  }, [onVerifyStep, status])

  async function sendSignInEmail() {
    const trimmedEmail = email.trim()
    if (!trimmedEmail) return

    try {
      await requestSignIn(trimmedEmail)
      setStatus('idle')
      setMessage('We sent a 6-digit code and sign-in link to your email.')
      focusCodeField()
    } catch (error) {
      emailRequestStarted.current = false
      flushSync(() => {
        setStep('email')
        setStatus('error')
        setMessage(error.message || 'Could not send sign-in email')
      })
    }
  }

  async function handleEmailContinue(event) {
    event.preventDefault()
    if (email.trim() === '' || emailRequestStarted.current) return

    prepareCodeEntry()

    emailRequestStarted.current = true
    await sendSignInEmail().finally(() => {
      emailRequestStarted.current = false
    })
  }

  async function handleFormSubmit(event) {
    event.preventDefault()

    if (!onVerifyStep) {
      await handleEmailContinue(event)
      return
    }

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
      focusCodeField()
    }
  }

  function handleUseDifferentEmail() {
    emailRequestStarted.current = false
    setStep('email')
    setCode('')
    setStatus('idle')
    setMessage('')
  }

  const emailContinueDisabled = (status === 'loading' && !onVerifyStep) || email.trim() === ''

  if (onVerifyStep) {
    return (
      <section className="panel-card" aria-labelledby="verify-heading">
        <h2 id="verify-heading">Enter your code</h2>
        <p className="panel-copy">
          Check <strong>{email}</strong> for a 6-digit code. Tap the field below to use autofill from Mail or Messages,
          or paste your code.
        </p>

        {message ? (
          <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'} role="status">
            {message}
          </p>
        ) : null}

        <form className="auth-form" onSubmit={handleFormSubmit}>
          <label htmlFor="login-code">Sign-in code</label>
          <input
            ref={codeInputRef}
            id="login-code"
            name="code"
            type="text"
            inputMode="numeric"
            autoComplete="one-time-code"
            autoCapitalize="off"
            autoCorrect="off"
            spellCheck={false}
            enterKeyHint="done"
            pattern="\d{6}"
            maxLength={6}
            required
            value={code}
            onChange={(event) => setCode(normalizeCode(event.target.value))}
            placeholder="123456"
            disabled={isSigningIn}
            className="auth-code-input"
          />
          <button type="submit" disabled={isSigningIn || normalizeCode(code).length !== 6}>
            {isSigningIn ? 'Signing in…' : 'Continue'}
          </button>
        </form>

        <button
          type="button"
          className="auth-link-button"
          onClick={handleUseDifferentEmail}
          disabled={isSigningIn}
        >
          Use a different email
        </button>
      </section>
    )
  }

  return (
    <section className="panel-card" aria-labelledby="login-heading">
      <h2 id="login-heading">Sign in</h2>
      <p className="panel-copy">Enter your email and we&apos;ll send a 6-digit code and sign-in link.</p>

      {message ? (
        <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'} role="status">
          {message}
        </p>
      ) : null}

      <form className="auth-form" onSubmit={handleFormSubmit}>
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
        <button
          type="submit"
          disabled={emailContinueDisabled}
          onTouchStart={() => {
            if (email.trim() === '') return
            prepareCodeEntry()
          }}
        >
          {status === 'loading' ? 'Sending…' : 'Continue'}
        </button>
      </form>
    </section>
  )
}
