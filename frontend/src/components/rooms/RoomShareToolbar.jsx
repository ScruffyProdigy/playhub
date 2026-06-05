import { useEffect, useState } from 'react'
import QRCode from 'qrcode'
import { roomShareText } from '../../lib/rooms'

export default function RoomShareToolbar({ joinUrl, inviteCode }) {
  const [qrOpen, setQrOpen] = useState(false)
  const [qrDataUrl, setQrDataUrl] = useState('')
  const [copyStatus, setCopyStatus] = useState('')

  useEffect(() => {
    if (!qrOpen || !joinUrl) {
      return undefined
    }
    let cancelled = false
    QRCode.toDataURL(joinUrl, { width: 240, margin: 2 })
      .then((url) => {
        if (!cancelled) {
          setQrDataUrl(url)
        }
      })
      .catch(() => {
        if (!cancelled) {
          setQrDataUrl('')
        }
      })
    return () => {
      cancelled = true
    }
  }, [qrOpen, joinUrl])

  async function copyLink() {
    const text = roomShareText(joinUrl)
    try {
      await navigator.clipboard.writeText(text)
      setCopyStatus('Copied!')
    } catch {
      setCopyStatus('Could not copy')
    }
    window.setTimeout(() => setCopyStatus(''), 2000)
  }

  async function nativeShare() {
    const text = roomShareText(joinUrl)
    if (navigator.share) {
      try {
        await navigator.share({ title: 'JoinQuest room', text, url: joinUrl })
      } catch {
        // user cancelled
      }
      return
    }
    await copyLink()
  }

  function openTextMessage() {
    const body = encodeURIComponent(roomShareText(joinUrl))
    window.location.href = `sms:?body=${body}`
  }

  return (
    <div className="room-share">
      <p className="room-share__code">
        Room code: <strong>{inviteCode}</strong>
      </p>
      <div className="room-share__actions">
        <button type="button" className="game-list-button" onClick={copyLink}>
          Copy invite
        </button>
        <button type="button" className="game-list-button" onClick={() => setQrOpen(true)}>
          QR code
        </button>
        <button type="button" className="game-list-button" onClick={nativeShare}>
          Share…
        </button>
        <button type="button" className="game-list-button game-list-button-secondary" onClick={openTextMessage}>
          Text
        </button>
      </div>
      {copyStatus ? <p className="room-share__status">{copyStatus}</p> : null}

      {qrOpen ? (
        <div className="room-qr-modal" role="dialog" aria-label="Room QR code">
          <div className="room-qr-modal__backdrop" onClick={() => setQrOpen(false)} aria-hidden="true" />
          <div className="room-qr-modal__panel panel-card">
            <h2>Scan to join</h2>
            <p className="panel-copy">Friends can open this link to join your room.</p>
            {qrDataUrl ? (
              <img className="room-qr-modal__image" src={qrDataUrl} alt={`QR code for ${joinUrl}`} />
            ) : (
              <p className="status-message">Generating QR…</p>
            )}
            <p className="room-qr-modal__url">{joinUrl}</p>
            <button type="button" className="game-list-button" onClick={() => setQrOpen(false)}>
              Close
            </button>
          </div>
        </div>
      ) : null}
    </div>
  )
}
