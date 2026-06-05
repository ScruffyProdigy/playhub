import { useCallback, useEffect, useRef, useState } from 'react'
import LoginForm from '../auth/LoginForm'
import useWaitForSignIn from '../auth/useWaitForSignIn'
import { useAuth } from '../auth/AuthProvider'
import {
  fetchRoom,
  joinRoom,
  leaveRoom,
  parseRoomInviteCode,
  sendRoomMessage,
  subscribeToRoom,
} from '../../lib/rooms'
import RoomShareToolbar from './RoomShareToolbar'
import { APP_NAME } from '../../lib/brand'

function displayName(user) {
  return user?.displayName?.trim() || 'Player'
}

function mergeMessage(messages, incoming) {
  if (!incoming?.id) {
    return messages
  }
  if (messages.some((msg) => msg.id === incoming.id)) {
    return messages
  }
  return [...messages, incoming]
}

export default function RoomPage() {
  const inviteCode = parseRoomInviteCode()
  const { user, loading: authLoading, refreshSession } = useAuth()
  const [room, setRoom] = useState(null)
  const [messages, setMessages] = useState([])
  const [draft, setDraft] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [liveError, setLiveError] = useState('')
  const chatEndRef = useRef(null)

  const loadRoom = useCallback(async () => {
    if (!inviteCode || !user) {
      return
    }
    setBusy(true)
    setError('')
    try {
      await joinRoom(inviteCode)
      const full = await fetchRoom(inviteCode)
      if (!full) {
        setError('Room not found.')
        return
      }
      setRoom(full)
      setMessages(full.messages || [])
    } catch (err) {
      setError(err.message || 'Could not join room.')
    } finally {
      setBusy(false)
    }
  }, [inviteCode, user])

  useWaitForSignIn({
    enabled: Boolean(inviteCode) && !user && !authLoading,
    onSignedIn: refreshSession,
  })

  useEffect(() => {
    if (user && inviteCode) {
      loadRoom()
    }
  }, [user, inviteCode, loadRoom])

  useEffect(() => {
    if (!room?.id) {
      return undefined
    }
    let unsubscribe = () => {}
    subscribeToRoom(room.id, {
      onRoomUpdate: (updated) => {
        setRoom((prev) => ({ ...prev, ...updated, messages: prev?.messages }))
      },
      onMessage: (msg) => {
        setMessages((prev) => mergeMessage(prev, msg))
      },
      onError: (msg) => setLiveError(msg),
    })
      .then((fn) => {
        unsubscribe = fn
      })
      .catch((err) => setLiveError(err.message || 'Live updates unavailable'))

    return () => unsubscribe()
  }, [room?.id])

  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  async function handleSend(event) {
    event.preventDefault()
    const body = draft.trim()
    if (!body || !room?.id || busy) {
      return
    }
    setBusy(true)
    setError('')
    try {
      const msg = await sendRoomMessage(room.id, body)
      setMessages((prev) => mergeMessage(prev, msg))
      setDraft('')
    } catch (err) {
      setError(err.message || 'Could not send message.')
    } finally {
      setBusy(false)
    }
  }

  async function handleLeave() {
    setBusy(true)
    setError('')
    try {
      await leaveRoom()
      window.location.href = '/'
    } catch (err) {
      setError(err.message || 'Could not leave room.')
      setBusy(false)
    }
  }

  if (!inviteCode) {
    return (
      <main className="app-shell auth-page">
        <h1>Invalid room link</h1>
        <p className="panel-copy">Use a link like /room/ABC123.</p>
        <a className="auth-link" href="/">
          Back to {APP_NAME}
        </a>
      </main>
    )
  }

  if (authLoading) {
    return (
      <main className="app-shell auth-page">
        <p className="status-message">Loading session…</p>
      </main>
    )
  }

  if (!user) {
    return (
      <main className="app-shell auth-page">
        <header className="app-header">
          <h1>Join room {inviteCode}</h1>
          <p className="tagline">Sign in to enter the chat.</p>
        </header>
        <section className="panel-card">
          <LoginForm />
        </section>
        <a className="auth-link" href="/">
          Back to {APP_NAME}
        </a>
      </main>
    )
  }

  return (
    <main className="app-shell room-page">
      <header className="app-header">
        <h1>Room {room?.inviteCode || inviteCode}</h1>
        <p className="tagline">
          {room?.members?.length ?? 0} {room?.members?.length === 1 ? 'member' : 'members'}
        </p>
      </header>

      {error ? <p className="status-message status-message-error">{error}</p> : null}
      {liveError ? <p className="status-message status-message-error">{liveError}</p> : null}

      {room ? (
        <>
          <section className="panel-card room-members">
            <h2>Members</h2>
            <ul className="room-members__list">
              {room.members.map((member) => (
                <li key={member.id}>
                  {displayName(member)}
                  {member.id === room.host?.id ? ' (host)' : ''}
                </li>
              ))}
            </ul>
          </section>

          <section className="panel-card">
            <RoomShareToolbar joinUrl={room.joinUrl} inviteCode={room.inviteCode} />
          </section>

          <section className="panel-card room-chat">
            <h2>Chat</h2>
            <div className="room-chat__log" aria-live="polite">
              {messages.length === 0 ? (
                <p className="panel-copy">Say hello — messages appear here for everyone in the room.</p>
              ) : (
                messages.map((msg) => (
                  <article key={msg.id} className="room-chat__message">
                    <p className="room-chat__author">{displayName(msg.author)}</p>
                    <p className="room-chat__body">{msg.body}</p>
                  </article>
                ))
              )}
              <div ref={chatEndRef} />
            </div>
            <form className="room-chat__composer auth-form" onSubmit={handleSend}>
              <label htmlFor="room-message">Message</label>
              <input
                id="room-message"
                type="text"
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                maxLength={2000}
                disabled={busy}
                autoComplete="off"
              />
              <button type="submit" className="game-list-button" disabled={busy || !draft.trim()}>
                Send
              </button>
            </form>
          </section>

          <div className="room-page__actions">
            <a className="game-list-button game-list-button-secondary" href="/">
              Catalog
            </a>
            <button type="button" className="game-list-button game-list-button-secondary" onClick={handleLeave} disabled={busy}>
              Leave room
            </button>
          </div>
        </>
      ) : (
        <p className="status-message">{busy ? 'Joining room…' : 'Loading room…'}</p>
      )}
    </main>
  )
}
