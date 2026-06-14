import { useEffect, useState } from 'react'
import { completeLinkEmailWithLinkOnce, isMergeConfirmationRequired } from '../../lib/auth'
import { notifyAuthComplete } from '../../lib/authBroadcast'
import { APP_NAME } from '../../lib/brand'
import { formatMergeWarning, MERGE_CANCEL, MERGE_CONFIRM } from '../../lib/playerCopy'
import { useAuth } from './AuthProvider'

function getTokenFromLocation() {
  const params = new URLSearchParams(window.location.search)
  return params.get('token')?.trim() || ''
}

export default function LinkEmailPage() {
  const { user, acceptSessionUser, refreshSession } = useAuth()
  const [status, setStatus] = useState('loading')
  const [message, setMessage] = useState('Verifying your email…')
  const [token, setToken] = useState('')

  useEffect(() => {
    const linkToken = getTokenFromLocation()
    if (!linkToken) {
      setStatus('error')
      setMessage('Missing verification token. Request a new code from Account settings.')
      return undefined
    }

    setToken(linkToken)
    let cancelled = false

    async function attemptLink(confirmMerge) {
      setStatus('loading')
      setMessage(confirmMerge ? 'Merging accounts…' : 'Verifying your email…')
      try {
        const updatedUser = await completeLinkEmailWithLinkOnce(linkToken, confirmMerge)
        if (cancelled) {
          return
        }
        acceptSessionUser(updatedUser)
        notifyAuthComplete()
        void refreshSession({ silent: true })
        setStatus('success')
        setMessage('Email linked. Redirecting…')
        window.history.replaceState({}, '', '/account')
        window.location.assign('/account')
      } catch (error) {
        if (cancelled) {
          return
        }
        if (isMergeConfirmationRequired(error.message)) {
          setStatus('merge')
          setMessage(error.message)
          return
        }
        setStatus('error')
        setMessage(error.message || 'Could not verify this email')
      }
    }

    void attemptLink(false)
    return () => {
      cancelled = true
    }
  }, [acceptSessionUser, refreshSession])

  async function handleConfirmMerge() {
    if (!token) {
      return
    }
    setStatus('loading')
    setMessage('Merging accounts…')
    try {
      const updatedUser = await completeLinkEmailWithLinkOnce(token, true)
      acceptSessionUser(updatedUser)
      notifyAuthComplete()
      void refreshSession({ silent: true })
      setStatus('success')
      setMessage('Email linked. Redirecting…')
      window.history.replaceState({}, '', '/account')
      window.location.assign('/account')
    } catch (error) {
      setStatus('error')
      setMessage(error.message || 'Could not verify this email')
    }
  }

  return (
    <main className="auth-page">
      <h1>{APP_NAME}</h1>
      <section className="panel-card" aria-live="polite">
        <h2>Link your email</h2>
        {status === 'merge' ? (
          <div className="account-merge-warning" role="alert">
            <p className="panel-copy">{formatMergeWarning(null, user?.displayName)}</p>
            <div className="account-merge-warning__actions">
              <button type="button" onClick={() => void handleConfirmMerge()}>
                {MERGE_CONFIRM}
              </button>
              <a className="button-secondary auth-link" href="/account">
                {MERGE_CANCEL}
              </a>
            </div>
          </div>
        ) : (
          <p className={status === 'error' ? 'status-message status-message-error' : 'status-message'}>{message}</p>
        )}
        {status === 'error' ? (
          <a className="auth-link" href="/account">
            Back to account settings
          </a>
        ) : null}
      </section>
    </main>
  )
}
