import { useCallback, useEffect, useRef, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import {
  fetchMyTableSeat,
  leaveTable,
  subscribeToMyTableSeat,
  TABLE_UPDATED_EVENT,
} from '../../lib/tables'
import { prefetchSubscriptionAuth } from '../../lib/queue'

export function useActiveTableSeat() {
  const { user, loading: authLoading } = useAuth()
  const [activeTableSeat, setActiveTableSeat] = useState(null)
  const [loading, setLoading] = useState(false)
  const [busy, setBusy] = useState(false)
  const unsubscribeRef = useRef(null)

  const refresh = useCallback(async () => {
    if (!user) {
      setActiveTableSeat(null)
      return
    }
    setLoading(true)
    try {
      const result = await fetchMyTableSeat()
      setActiveTableSeat(result)
    } catch {
      setActiveTableSeat((prev) => (prev?.status === 'started' ? prev : null))
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
    const handler = () => {
      void refresh()
    }
    window.addEventListener(TABLE_UPDATED_EVENT, handler)
    return () => window.removeEventListener(TABLE_UPDATED_EVENT, handler)
  }, [refresh])

  useEffect(() => {
    unsubscribeRef.current?.()
    unsubscribeRef.current = null

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
            setActiveTableSeat(seat)
          },
        })
        if (cancelled) {
          unsubscribe()
          return
        }
        unsubscribeRef.current = unsubscribe
      } catch {
        // Banner still works via refresh and TABLE_UPDATED_EVENT.
      }
    })()

    return () => {
      cancelled = true
      unsubscribeRef.current?.()
      unsubscribeRef.current = null
    }
  }, [authLoading, user, refresh])

  useEffect(() => {
    if (!activeTableSeat?.tableId || activeTableSeat.status === 'started') {
      return undefined
    }
    const timer = window.setInterval(() => {
      void refresh()
    }, 3000)
    return () => window.clearInterval(timer)
  }, [activeTableSeat?.tableId, activeTableSeat?.status, refresh])

  const handleLeave = useCallback(async () => {
    if (!activeTableSeat?.tableId) {
      return
    }
    setBusy(true)
    try {
      await leaveTable(activeTableSeat.tableId)
      await refresh()
    } finally {
      setBusy(false)
    }
  }, [activeTableSeat?.tableId, refresh])

  return { activeTableSeat, loading, busy, refresh, handleLeave }
}
