import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from 'react'
import { useAuth } from '../auth/AuthProvider'
import { navigateTo } from '../../lib/usePathname'
import {
  fetchMyRoom,
  fetchRoom,
  joinRoom,
  leaveRoom as leaveRoomApi,
  subscribeToRoom,
} from '../../lib/rooms'
import { prefetchSubscriptionAuth } from '../../lib/queue'

const ActiveRoomContext = createContext(null)

function mergeMessage(messages, incoming) {
  if (!incoming?.id) {
    return messages
  }
  if (messages.some((msg) => msg.id === incoming.id)) {
    return messages
  }
  return [...messages, incoming]
}

function countUnread(messages, lastReadId, userId) {
  if (!messages?.length) {
    return 0
  }
  if (!lastReadId) {
    return messages.filter((msg) => msg.author?.id !== userId).length
  }
  let pastLastRead = false
  let count = 0
  for (const msg of messages) {
    if (msg.id === lastReadId) {
      pastLastRead = true
      continue
    }
    if (pastLastRead && msg.author?.id !== userId) {
      count += 1
    }
  }
  return count
}

export function ActiveRoomProvider({ children, pendingInviteCode = null }) {
  const { user, loading: authLoading } = useAuth()
  const [room, setRoom] = useState(null)
  const [messages, setMessages] = useState([])
  const [roomOpen, setRoomOpen] = useState(() => Boolean(pendingInviteCode))
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [lastReadMessageId, setLastReadMessageId] = useState(null)
  const unsubscribeRef = useRef(null)

  const unreadCount = useMemo(
    () => countUnread(messages, lastReadMessageId, user?.id),
    [messages, lastReadMessageId, user?.id],
  )

  const markRead = useCallback(() => {
    const latest = messages[messages.length - 1]
    if (latest?.id) {
      setLastReadMessageId(latest.id)
    }
  }, [messages])

  const applyRoom = useCallback((nextRoom, nextMessages = []) => {
    setRoom(nextRoom)
    setMessages(nextMessages)
    const latest = nextMessages[nextMessages.length - 1]
    setLastReadMessageId(latest?.id ?? null)
    setError('')
  }, [])

  const refresh = useCallback(async () => {
    if (!user) {
      setRoom(null)
      setMessages([])
      setRoomOpen(false)
      return null
    }
    setLoading(true)
    try {
      const mine = await fetchMyRoom()
      if (mine) {
        applyRoom(mine, mine.messages || [])
        return mine
      }
      setRoom(null)
      setMessages([])
      return null
    } catch {
      setRoom(null)
      setMessages([])
      return null
    } finally {
      setLoading(false)
    }
  }, [user, applyRoom])

  const openRoomByCode = useCallback(
    async (inviteCode, { openPanel = true } = {}) => {
      if (!user || !inviteCode) {
        return null
      }
      setLoading(true)
      setError('')
      try {
        await joinRoom(inviteCode)
        const full = await fetchRoom(inviteCode)
        if (!full) {
          setError('Room not found.')
          return null
        }
        applyRoom(full, full.messages || [])
        if (openPanel) {
          setRoomOpen(true)
          navigateTo(`/room/${full.inviteCode}`, { replace: true })
        }
        return full
      } catch (err) {
        setError(err.message || 'Could not join room.')
        return null
      } finally {
        setLoading(false)
      }
    },
    [user, applyRoom],
  )

  const openRoom = useCallback(() => {
    if (!room?.inviteCode) {
      return
    }
    setRoomOpen(true)
    markRead()
    navigateTo(`/room/${room.inviteCode}`)
  }, [room?.inviteCode, markRead])

  const dismissRoom = useCallback(() => {
    setRoomOpen(false)
    navigateTo('/', { replace: true })
  }, [])

  const openRoomFromCreate = useCallback(
    (createdRoom) => {
      applyRoom(createdRoom, [])
      setRoomOpen(true)
      navigateTo(`/room/${createdRoom.inviteCode}`)
    },
    [applyRoom],
  )

  const handleLeave = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      await leaveRoomApi()
      setRoom(null)
      setMessages([])
      setRoomOpen(false)
      setLastReadMessageId(null)
      navigateTo('/', { replace: true })
    } catch (err) {
      setError(err.message || 'Could not leave room.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (authLoading) {
      return
    }
    if (!user) {
      setRoom(null)
      setMessages([])
      setRoomOpen(false)
      return
    }
    void refresh()
  }, [authLoading, user, refresh])

  useEffect(() => {
    if (authLoading || !user) {
      return
    }
    if (pendingInviteCode) {
      void openRoomByCode(pendingInviteCode, { openPanel: true })
      return
    }
    setRoomOpen(false)
  }, [authLoading, user, pendingInviteCode, openRoomByCode])

  useEffect(() => {
    if (roomOpen) {
      markRead()
    }
  }, [roomOpen, markRead, messages])

  useEffect(() => {
    unsubscribeRef.current?.()
    unsubscribeRef.current = null

    if (!user || !room?.id) {
      return undefined
    }

    let cancelled = false

    void (async () => {
      try {
        await prefetchSubscriptionAuth()
        if (cancelled) {
          return
        }
        const unsubscribe = await subscribeToRoom(room.id, {
          onRoomUpdate: (updated) => {
            setRoom((prev) => (prev ? { ...prev, ...updated } : prev))
          },
          onMessage: (msg) => {
            setMessages((prev) => mergeMessage(prev, msg))
            if (!roomOpen && msg.author?.id !== user.id) {
              // unread derived from lastReadMessageId
            }
          },
        })
        if (cancelled) {
          unsubscribe()
          return
        }
        unsubscribeRef.current = unsubscribe
      } catch {
        // chat still works; refresh on focus
      }
    })()

    return () => {
      cancelled = true
      unsubscribeRef.current?.()
      unsubscribeRef.current = null
    }
  }, [user, room?.id, roomOpen])

  const value = useMemo(
    () => ({
      room,
      messages,
      setMessages,
      roomOpen,
      loading,
      error,
      setError,
      unreadCount,
      refresh,
      openRoom,
      dismissRoom,
      openRoomByCode,
      openRoomFromCreate,
      handleLeave,
      markRead,
    }),
    [
      room,
      messages,
      roomOpen,
      loading,
      error,
      unreadCount,
      refresh,
      openRoom,
      dismissRoom,
      openRoomByCode,
      openRoomFromCreate,
      handleLeave,
      markRead,
    ],
  )

  return <ActiveRoomContext.Provider value={value}>{children}</ActiveRoomContext.Provider>
}

export function useActiveRoom() {
  const context = useContext(ActiveRoomContext)
  if (!context) {
    throw new Error('useActiveRoom must be used within ActiveRoomProvider')
  }
  return context
}
