export default function RoomSheet({ children, onDismiss }) {
  return (
    <div className="room-sheet" role="dialog" aria-modal="true" aria-label="Room chat">
      <button type="button" className="room-sheet__backdrop" onClick={onDismiss} aria-label="Dismiss room" />
      <div className="room-sheet__panel panel-card">{children}</div>
    </div>
  )
}
