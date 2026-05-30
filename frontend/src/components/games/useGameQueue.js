import { useCallback, useEffect, useRef, useState } from 'react'
import {
  fetchMyQueueStatus,
  joinGame,
  leaveQueue,
  prefetchSubscriptionAuth,
  subscribeToQueue,
} from '../../lib/queue'

/** Apply joinGame mutation result (waiting or matched for the player who triggered the match). */
function applyJoinResponse(result, handlers) {
  if (result?.joinUrl) {
    return applyQueueUpdate(
      { status: 'MATCHED', joinUrl: result.joinUrl, queuedCount: 0 },
      handlers,
    )
  }
  if (result?.queued) {
    return applyQueueUpdate(
      { status: 'WAITING', queuedCount: result.queuedCount ?? 1 },
      handlers,
    )
  }
  return false
}

function applyQueueUpdate(update, { setQueueState, setJoinUrl, setQueuedCount, setError }) {
  if (update.status === 'MATCHED' && update.joinUrl) {
    setQueueState('matched')
    setJoinUrl(update.joinUrl)
    setError('')
    return true
  }
  if (update.status === 'WAITING') {
    setQueueState('waiting')
    setQueuedCount(update.queuedCount ?? 0)
    setError('')
    return false
  }
  if (update.status === 'LEFT') {
    setQueueState('idle')
    setJoinUrl('')
    setQueuedCount(0)
    setError('')
    return false
  }
  return false
}

/** Map myQueueStatus query result into local queue UI state. */
function applyMyQueueStatus(result, handlers) {
  if (!result?.queued) {
    if (result?.joinUrl) {
      return applyQueueUpdate(
        { status: 'MATCHED', joinUrl: result.joinUrl, queuedCount: 0 },
        handlers,
      )
    }
    applyQueueUpdate({ status: 'LEFT', queuedCount: 0 }, handlers)
    return false
  }
  if (result.joinUrl) {
    return applyQueueUpdate(
      { status: 'MATCHED', joinUrl: result.joinUrl, queuedCount: result.queuedCount ?? 0 },
      handlers,
    )
  }
  return applyQueueUpdate(
    { status: 'WAITING', queuedCount: result.queuedCount ?? 0 },
    handlers,
  )
}

export function useGameQueue(gameId) {
  const [queueState, setQueueState] = useState('idle')
  const [joinUrl, setJoinUrl] = useState('')
  const [queuedCount, setQueuedCount] = useState(0)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const unsubscribeRef = useRef(null)
  const subscribeGenerationRef = useRef(0)
  const syncGenerationRef = useRef(0)

  const clearSubscription = useCallback(() => {
    unsubscribeRef.current?.()
    unsubscribeRef.current = null
  }, [])

  const applyUpdate = useCallback(
    (update) =>
      applyQueueUpdate(update, {
        setQueueState,
        setJoinUrl,
        setQueuedCount,
        setError,
      }),
    [],
  )

  const ensureSubscribed = useCallback(async () => {
    if (unsubscribeRef.current) {
      return
    }

    const generation = ++subscribeGenerationRef.current

    try {
      await prefetchSubscriptionAuth()
    } catch (err) {
      if (generation === subscribeGenerationRef.current) {
        setError(err.message || 'Sign in required for live queue updates')
      }
      return
    }

    if (generation !== subscribeGenerationRef.current) {
      return
    }

    try {
      const unsubscribe = await subscribeToQueue(gameId, {
        onUpdate: (update) => {
          applyUpdate(update)
        },
        onError: (message) => {
          if (generation === subscribeGenerationRef.current) {
            setError(message)
          }
        },
      })
      if (generation === subscribeGenerationRef.current) {
        unsubscribeRef.current = unsubscribe
      } else {
        unsubscribe()
      }
    } catch (err) {
      if (generation === subscribeGenerationRef.current) {
        setError(err.message || 'Queue updates unavailable')
      }
    }
  }, [gameId, applyUpdate])

  // Cross-tab sync: subscribe on mount and hydrate from server state.
  useEffect(() => {
    let cancelled = false

    const syncGen = ++syncGenerationRef.current

    async function syncFromServer() {
      try {
        await prefetchSubscriptionAuth()
        if (cancelled || syncGen !== syncGenerationRef.current) {
          return
        }
        await ensureSubscribed()
        if (cancelled || syncGen !== syncGenerationRef.current) {
          return
        }
        const status = await fetchMyQueueStatus(gameId)
        if (cancelled || syncGen !== syncGenerationRef.current) {
          return
        }
        applyMyQueueStatus(status, {
          setQueueState,
          setJoinUrl,
          setQueuedCount,
          setError,
        })
      } catch {
        // Signed-out or transient errors: keep idle until user joins.
      }
    }

    void syncFromServer()

    const onVisible = () => {
      if (document.visibilityState === 'visible') {
        const visGen = syncGenerationRef.current
        void fetchMyQueueStatus(gameId)
          .then((result) => {
            if (cancelled || visGen !== syncGenerationRef.current) {
              return
            }
            applyMyQueueStatus(result, {
              setQueueState,
              setJoinUrl,
              setQueuedCount,
              setError,
            })
          })
          .catch(() => {})
      }
    }
    document.addEventListener('visibilitychange', onVisible)

    return () => {
      cancelled = true
      subscribeGenerationRef.current += 1
      document.removeEventListener('visibilitychange', onVisible)
      clearSubscription()
    }
  }, [gameId, ensureSubscribed, clearSubscription])

  async function handleJoin() {
    setBusy(true)
    setError('')
    syncGenerationRef.current += 1
    try {
      await ensureSubscribed()
      const result = await joinGame(gameId)
      applyJoinResponse(result, {
        setQueueState,
        setJoinUrl,
        setQueuedCount,
        setError,
      })
    } catch (err) {
      setError(err.message || 'Could not join queue')
    } finally {
      setBusy(false)
    }
  }

  async function handleLeave() {
    setBusy(true)
    setError('')
    syncGenerationRef.current += 1
    try {
      await leaveQueue(gameId)
      setQueueState('idle')
      setJoinUrl('')
      setQueuedCount(0)
    } catch (err) {
      setError(err.message || 'Could not leave queue')
    } finally {
      setBusy(false)
    }
  }

  return {
    queueState,
    joinUrl,
    queuedCount,
    error,
    busy,
    handleJoin,
    handleLeave,
  }
}
