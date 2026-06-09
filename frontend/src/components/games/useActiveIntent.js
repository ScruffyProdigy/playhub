import { useCallback, useEffect, useRef, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import { fetchMyActiveIntent, leaveActiveGame, resolveIntentLaunchUrl } from '../../lib/intent'
import { fetchMyQueueStatus, leaveQueue, prefetchSubscriptionAuth, subscribeToQueue } from '../../lib/queue'
import { subscribeToMyTableSeat, TABLE_UPDATED_EVENT, fetchMyTableSeat } from '../../lib/tables'
import { useActiveTableSeat } from './useActiveTableSeat'

function mergeMatchedJoinUrl(intent, joinUrl) {
  if (!intent || intent.status !== 'MATCHED' || !joinUrl || intent.joinUrl === joinUrl) {
    return intent
  }
  return { ...intent, joinUrl }
}

export function useActiveIntent() {
  const { user, loading: authLoading } = useAuth()
  const [activeIntent, setActiveIntent] = useState(null)
  const [loading, setLoading] = useState(false)
  const [busy, setBusy] = useState(false)
  const queueUnsubRef = useRef(null)
  const seatUnsubRef = useRef(null)

  const {
    activeTableSeat,
    loading: tableLoading,
    busy: tableBusy,
    refresh: refreshTable,
    handleLeave: leaveTableSeat,
  } = useActiveTableSeat()

  const enrichMatchedJoinUrl = useCallback(async (intent, tableSeat) => {
    if (intent?.status !== 'MATCHED' || intent?.joinUrl) {
      return intent
    }
    if (tableSeat?.joinUrl) {
      return { ...intent, joinUrl: tableSeat.joinUrl }
    }
    if (!intent?.queueId) {
      return intent
    }
    try {
      const status = await fetchMyQueueStatus(intent.queueId)
      if (status?.joinUrl) {
        return { ...intent, joinUrl: status.joinUrl }
      }
    } catch {
      // Provision may still be in flight; banner can retry on the next refresh.
    }
    return intent
  }, [])

  const refreshIntent = useCallback(async () => {
    if (!user) {
      setActiveIntent(null)
      return
    }
    setLoading(true)
    try {
      const [result, tableSeat] = await Promise.all([fetchMyActiveIntent(), fetchMyTableSeat()])
      setActiveIntent(await enrichMatchedJoinUrl(result, tableSeat))
    } catch {
      setActiveIntent(null)
    } finally {
      setLoading(false)
    }
  }, [enrichMatchedJoinUrl, user])

  const refresh = useCallback(async () => {
    await Promise.all([refreshIntent(), refreshTable()])
  }, [refreshIntent, refreshTable])

  useEffect(() => {
    if (authLoading) {
      return
    }
    void refreshIntent()
  }, [authLoading, refreshIntent])

  useEffect(() => {
    const handler = () => {
      void refresh()
    }
    window.addEventListener(TABLE_UPDATED_EVENT, handler)
    return () => window.removeEventListener(TABLE_UPDATED_EVENT, handler)
  }, [refresh])

  useEffect(() => {
    queueUnsubRef.current?.()
    queueUnsubRef.current = null

    const queueId = activeIntent?.queueId
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
          onUpdate: (update) => {
            if (update?.status === 'MATCHED' && update?.joinUrl) {
              setActiveIntent((prev) => mergeMatchedJoinUrl(prev, update.joinUrl))
            }
            void refreshIntent()
          },
        })
        if (cancelled) {
          unsubscribe()
          return
        }
        queueUnsubRef.current = unsubscribe
      } catch {
        // Banner still works via refresh.
      }
    })()

    return () => {
      cancelled = true
      queueUnsubRef.current?.()
      queueUnsubRef.current = null
    }
  }, [user, activeIntent?.queueId, refreshIntent])

  useEffect(() => {
    seatUnsubRef.current?.()
    seatUnsubRef.current = null

    if (authLoading || !user) {
      return undefined
    }

    let cancelled = false

    void (async () => {
      try {
        await prefetchSubscriptionAuth()
        if (cancelled) {
          return
        }
        const unsubscribe = await subscribeToMyTableSeat({
          onUpdate: (seat) => {
            if (seat?.joinUrl) {
              setActiveIntent((prev) => mergeMatchedJoinUrl(prev, seat.joinUrl))
            }
            void refresh()
          },
        })
        if (cancelled) {
          unsubscribe()
          return
        }
        seatUnsubRef.current = unsubscribe
      } catch {
        // Banner still works via refresh.
      }
    })()

    return () => {
      cancelled = true
      seatUnsubRef.current?.()
      seatUnsubRef.current = null
    }
  }, [authLoading, user, refresh])

  useEffect(() => {
    if (activeIntent?.status !== 'MATCHED' || resolveIntentLaunchUrl(activeIntent, activeTableSeat)) {
      return undefined
    }
    const timer = window.setInterval(() => {
      void refreshIntent()
    }, 3000)
    return () => window.clearInterval(timer)
  }, [activeIntent, activeTableSeat, refreshIntent])

  const leaveQueueIntent = useCallback(async () => {
    if (!activeIntent?.queueId) {
      return
    }
    setBusy(true)
    try {
      await leaveQueue(activeIntent.queueId)
      await refreshIntent()
    } finally {
      setBusy(false)
    }
  }, [activeIntent?.queueId, refreshIntent])

  const handleLeave = useCallback(async () => {
    if (activeIntent?.status === 'MATCHED' || activeTableSeat?.status === 'started') {
      setBusy(true)
      try {
        await leaveActiveGame()
        await refresh()
      } finally {
        setBusy(false)
      }
      return
    }
    if (activeIntent?.status === 'WAITING') {
      await leaveQueueIntent()
      return
    }
    await leaveTableSeat()
  }, [activeIntent?.status, activeTableSeat?.status, leaveQueueIntent, leaveTableSeat, refresh])

  return {
    activeIntent,
    activeTableSeat,
    loading: loading || tableLoading,
    busy: busy || tableBusy,
    refresh,
    handleLeave,
  }
}
