import { useEffect, useState } from 'react'

const EXIT_ANIMATION_MS = 300

export default function RoomSheet({ open, onDismiss, children }) {
  const [render, setRender] = useState(open)
  const [exiting, setExiting] = useState(false)

  useEffect(() => {
    if (open) {
      setRender(true)
      setExiting(false)
      return undefined
    }
    if (render) {
      setExiting(true)
      const timer = window.setTimeout(() => {
        setRender(false)
        setExiting(false)
      }, EXIT_ANIMATION_MS)
      return () => window.clearTimeout(timer)
    }
    return undefined
  }, [open, render])

  if (!render) {
    return null
  }

  return (
    <div className={`room-sheet ${exiting ? 'room-sheet--exiting' : ''}`} role="dialog" aria-modal="true" aria-label="Room chat">
      <button type="button" className="room-sheet__backdrop" onClick={onDismiss} aria-label="Dismiss room" />
      <div className={`room-sheet__panel panel-card ${exiting ? 'room-sheet__panel--exiting' : 'room-sheet__panel--entering'}`}>
        {children}
      </div>
    </div>
  )
}
