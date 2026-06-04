import { useCallback, useEffect, useRef, useState } from 'react'
import {
  fetchMyQueueStatus,
  joinQueue,
  leaveQueue,
  prefetchSubscriptionAuth,
  subscribeToQueue,
} from '../../lib/queue'

/** Apply join mutation result (waiting or matched for the player who triggered the match). */
function applyJoinResponse(result, handlers) {
  if (result?.message) {
    handlers.setNotice?.(result.message)
  }
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
    if (update.message) {
      setError(update.message)
    } else {
      setError('')
    }
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

/**
 * @param {string | undefined} queueId
 * @param {{ skipSubscription?: boolean }} [options] — set when the sticky banner already subscribes to this queue
 */
export function useGameQueue(queueId, { skipSubscription = false } = {}) {
  const [queueState, setQueueState] = useState('idle')
  const [joinUrl, setJoinUrl] = useState('')
  const [queuedCount, setQueuedCount] = useState(0)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
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
    if (!queueId || unsubscribeRef.current) {
      return
    }

    const generation = ++subscribeGenerationRef.current

    try {
      await prefetchSubscriptionAuth()
    } catch (err) {
      if (generation === subscribeGenerationRef.current) {
        setError(err.message || 'Sign in required for live updates')
      }
      return
    }

    if (generation !== subscribeGenerationRef.current) {
      return
    }

    try {
      const unsubscribe = await subscribeToQueue(queueId, {
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
        setError(err.message || 'Live updates unavailable')
      }
    }
  }, [queueId, applyUpdate])

  useEffect(() => {
    if (!queueId) {
      return undefined
    }

    let cancelled = false
    const syncGen = ++syncGenerationRef.current

    async function syncFromServer() {
      try {
        await prefetchSubscriptionAuth()
        if (cancelled || syncGen !== syncGenerationRef.current) {
          return
        }
        if (!skipSubscription) {
          await ensureSubscribed()
          if (cancelled || syncGen !== syncGenerationRef.current) {
            return
          }
        }
        const status = await fetchMyQueueStatus(queueId)
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
        void fetchMyQueueStatus(queueId)
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
  }, [queueId, skipSubscription, ensureSubscribed, clearSubscription])

  async function handleJoin() {
    if (!queueId) {
      setError('This game is not available for group matchmaking yet')
      return
    }
    setBusy(true)
    setError('')
    setNotice('')
    syncGenerationRef.current += 1
    try {
      if (!skipSubscription) {
        await ensureSubscribed()
      }
      const result = await joinQueue(queueId)
      applyJoinResponse(result, {
        setQueueState,
        setJoinUrl,
        setQueuedCount,
        setError,
        setNotice,
      })
    } catch (err) {
      setError(err.message || 'Could not start looking for a group')
    } finally {
      setBusy(false)
    }
  }

  async function handleLeave() {
    if (!queueId) {
      return
    }
    setBusy(true)
    setError('')
    syncGenerationRef.current += 1
    try {
      await leaveQueue(queueId)
      setQueueState('idle')
      setJoinUrl('')
      setQueuedCount(0)
    } catch (err) {
      setError(err.message || 'Could not stop looking')
    } finally {
      setBusy(false)
    }
  }

  return {
    queueState,
    joinUrl,
    queuedCount,
    error,
    notice,
    busy,
    handleJoin,
    handleLeave,
  }
}
