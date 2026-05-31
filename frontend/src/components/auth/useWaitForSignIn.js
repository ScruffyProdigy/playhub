import { useEffect } from 'react'
import { onAuthComplete } from '../../lib/authBroadcast'

const SESSION_POLL_MS = 2500

export default function useWaitForSignIn({ enabled, onSignedIn }) {
  useEffect(() => {
    if (!enabled) {
      return undefined
    }

    return onAuthComplete(onSignedIn)
  }, [enabled, onSignedIn])

  useEffect(() => {
    if (!enabled) {
      return undefined
    }

    const intervalId = window.setInterval(onSignedIn, SESSION_POLL_MS)
    return () => window.clearInterval(intervalId)
  }, [enabled, onSignedIn])
}
