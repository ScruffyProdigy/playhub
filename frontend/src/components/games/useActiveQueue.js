import { useCallback, useEffect, useRef, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import { fetchMyActiveQueue } from '../../lib/activeQueue'
import { leaveQueue, prefetchSubscriptionAuth, subscribeToQueue } from '../../lib/queue'

export function useActiveQueue() {
  const { user, loading: authLoading } = useAuth()
  const [activeQueue, setActiveQueue] = useState(null)
  const [loading, setLoading] = useState(false)
  const [busy, setBusy] = useState(false)
  const unsubscribeRef = useRef(null)

  const refresh = useCallback(async () => {
    if (!user) {
      setActiveQueue(null)
      return
    }
    setLoading(true)
    try {
      const result = await fetchMyActiveQueue()
      setActiveQueue(result)
    } catch {
      setActiveQueue(null)
    } finally {
      setLoading(false)
    }
  }, [user])

  useEffect(() => {
    if (authLoading) {
      return
    }
    void refresh()
  }, [authLoading, refresh])

  useEffect(() => {
    unsubscribeRef.current?.()
    unsubscribeRef.current = null

    const queueId = activeQueue?.queueId
    if (!user || !queueId) {
      return undefined
    }

    let cancelled = false

    void (async () => {
      try {
        await prefetchSubscriptionAuth()
        if (cancelled) {
          return
        }
        const unsubscribe = await subscribeToQueue(queueId, {
          onUpdate: () => {
            void refresh()
          },
        })
        if (cancelled) {
          unsubscribe()
          return
        }
        unsubscribeRef.current = unsubscribe
      } catch {
        // Banner still works; refresh on tab focus via GameLobby items.
      }
    })()

    return () => {
      cancelled = true
      unsubscribeRef.current?.()
      unsubscribeRef.current = null
    }
  }, [user, activeQueue?.queueId, refresh])

  const handleLeave = useCallback(async () => {
    if (!activeQueue?.queueId) {
      return
    }
    setBusy(true)
    try {
      await leaveQueue(activeQueue.queueId)
      await refresh()
    } finally {
      setBusy(false)
    }
  }, [activeQueue?.queueId, refresh])

  return { activeQueue, loading, busy, refresh, handleLeave }
}
