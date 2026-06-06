import { useCallback, useEffect, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import { fetchMyTableSeat, leaveTable } from '../../lib/tables'

export function useActiveTableSeat() {
  const { user, loading: authLoading } = useAuth()
  const [activeTableSeat, setActiveTableSeat] = useState(null)
  const [loading, setLoading] = useState(false)
  const [busy, setBusy] = useState(false)

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
      setActiveTableSeat(null)
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
