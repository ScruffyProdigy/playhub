import { useState } from 'react'
import { createGuestSession } from '../../lib/auth'
import { notifyAuthComplete } from '../../lib/authBroadcast'
import {
  JUMP_IN,
  JUMP_IN_HINT,
  SIGN_IN_DIVIDER,
  SIGN_IN_DIVIDER_LABEL,
  SIGN_IN_HEADING,
} from '../../lib/playerCopy'
import { useAuth } from './AuthProvider'
import EmailSignInForm from './EmailSignInForm'
import SocialSignInRow from './SocialSignInRow'

export default function SignInPanel() {
  const { acceptSessionUser, refreshSession } = useAuth()
  const [status, setStatus] = useState('idle')
  const [message, setMessage] = useState('')

  async function handleJumpIn() {
    if (status === 'loading') {
      return
    }
    setStatus('loading')
    setMessage('')
    try {
      const guestUser = await createGuestSession()
      acceptSessionUser(guestUser)
      notifyAuthComplete()
      void refreshSession({ silent: true })
    } catch (error) {
      setStatus('error')
      setMessage(error.message || 'Could not start a guest session')
    } finally {
      setStatus('idle')
    }
  }

  return (
    <section className="panel-card auth-sign-in" aria-labelledby="sign-in-heading">
      <h2 id="sign-in-heading">{SIGN_IN_HEADING}</h2>

      <div className="auth-sign-in__guest">
        <button type="button" className="auth-sign-in__primary" onClick={() => void handleJumpIn()} disabled={status === 'loading'}>
          {status === 'loading' ? 'Starting…' : JUMP_IN}
        </button>
        <p className="panel-copy auth-sign-in__hint">{JUMP_IN_HINT}</p>
      </div>

      {message ? (
        <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'} role="status">
          {message}
        </p>
      ) : null}

      <div className="auth-sign-in__divider" role="separator" aria-label={SIGN_IN_DIVIDER_LABEL}>
        <span>{SIGN_IN_DIVIDER}</span>
      </div>

      <div className="auth-sign-in__email">
        <EmailSignInForm />
      </div>

      <SocialSignInRow />
    </section>
  )
}
