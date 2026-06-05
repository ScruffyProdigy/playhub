import { useEffect, useState } from 'react'
import QRCode from 'qrcode'
import { roomShareText } from '../../lib/rooms'
import { IconCopy, IconQr, IconShare } from '../icons/ShareIcons'

function IconSms(props) {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true" {...props}>
      <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
    </svg>
  )
}

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
      <div className="room-share__actions" role="group" aria-label="Share room">
        <button type="button" className="room-share__icon-btn" onClick={copyLink} aria-label="Copy invite">
          <IconCopy />
          <span className="room-share__icon-label">Copy</span>
        </button>
        <button type="button" className="room-share__icon-btn" onClick={() => setQrOpen(true)} aria-label="Show QR code">
          <IconQr />
          <span className="room-share__icon-label">QR</span>
        </button>
        <button type="button" className="room-share__icon-btn" onClick={nativeShare} aria-label="Share">
          <IconShare />
          <span className="room-share__icon-label">Share</span>
        </button>
        <button type="button" className="room-share__icon-btn" onClick={openTextMessage} aria-label="Text invite">
          <IconSms />
          <span className="room-share__icon-label">Text</span>
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
