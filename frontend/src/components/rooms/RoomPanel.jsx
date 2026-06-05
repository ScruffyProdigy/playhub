import { useEffect, useRef, useState } from 'react'
import LoginForm from '../auth/LoginForm'
import { useAuth } from '../auth/AuthProvider'
import { useActiveRoom } from './ActiveRoomProvider'
import { sendRoomMessage } from '../../lib/rooms'
import { IconClose } from '../icons/ShareIcons'
import RoomShareToolbar from './RoomShareToolbar'

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

export default function RoomPanel({ compact = false }) {
  const { user, loading: authLoading } = useAuth()
  const {
    room,
    messages,
    setMessages,
    loading,
    error,
    setError,
    handleLeave,
  } = useActiveRoom()
  const [draft, setDraft] = useState('')
  const [busy, setBusy] = useState(false)
  const chatEndRef = useRef(null)

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

  if (authLoading) {
    return <p className="status-message">Loading session…</p>
  }

  if (!user) {
    return (
      <div className="room-panel room-panel--guest">
        <h2 className="room-panel__title">Join room</h2>
        <p className="panel-copy">Sign in to enter the chat.</p>
        <LoginForm />
      </div>
    )
  }

  if (!room && loading) {
    return <p className="status-message">{busy ? 'Joining room…' : 'Loading room…'}</p>
  }

  if (!room) {
    return null
  }

  const memberCount = room.members?.length ?? 0
  const membersList = (
    <ul className="room-members__list">
      {room.members.map((member) => (
        <li key={member.id}>
          {displayName(member)}
          {member.id === room.host?.id ? ' (host)' : ''}
        </li>
      ))}
    </ul>
  )

  return (
    <div className={`room-panel ${compact ? 'room-panel--compact' : ''}`}>
      <div className="room-panel__body">
        <header className="room-panel__header">
          <div>
            <h2 className="room-panel__title">Room {room.inviteCode}</h2>
            <p className="room-panel__meta">
              {memberCount} {memberCount === 1 ? 'member' : 'members'}
            </p>
          </div>
        </header>

        {error ? <p className="status-message status-message-error">{error}</p> : null}

        {compact ? (
          <details className="room-panel__section room-panel__collapsible room-members">
            <summary className="room-panel__section-title">
              Members ({memberCount})
            </summary>
            {membersList}
          </details>
        ) : (
          <section className="room-panel__section room-members">
            <h3 className="room-panel__section-title">Members</h3>
            {membersList}
          </section>
        )}

        <section className="room-panel__section room-panel__share">
          <RoomShareToolbar joinUrl={room.joinUrl} inviteCode={room.inviteCode} />
        </section>

        <section className="room-panel__section room-panel__chat room-chat">
          <h3 className="room-panel__section-title">Chat</h3>
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
        </section>
      </div>

      <form className="room-panel__composer room-chat__composer auth-form" onSubmit={handleSend}>
        <label htmlFor="room-message">Message</label>
        <input
          id="room-message"
          type="text"
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          maxLength={2000}
          disabled={busy}
          autoComplete="off"
          placeholder="Message"
        />
        <button type="submit" className="game-list-button" disabled={busy || !draft.trim()}>
          Send
        </button>
      </form>

      <div className="room-panel__leave">
        <button
          type="button"
          className="room-panel__leave-btn"
          onClick={handleLeave}
          disabled={busy || loading}
        >
          <IconClose aria-hidden="true" />
          <span>Leave room</span>
          <IconClose aria-hidden="true" />
        </button>
      </div>
    </div>
  )
}
