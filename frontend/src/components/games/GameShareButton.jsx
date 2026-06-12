import { useState } from 'react'
import { gamePageShareUrl } from '../../lib/gameCard'
import { IconShare } from '../icons/ShareIcons'

export default function GameShareButton({ game, className = '' }) {
  const [status, setStatus] = useState('')
  const shareUrl = gamePageShareUrl(game)

  async function handleShare() {
    if (!shareUrl) {
      return
    }

    const title = game?.name?.trim() || 'JoinQuest game'
    const text = `Check out ${title} on JoinQuest`

    if (navigator.share) {
      try {
        await navigator.share({ title, text, url: shareUrl })
        return
      } catch {
        // user cancelled or share failed
      }
    }

    try {
      await navigator.clipboard.writeText(shareUrl)
      setStatus('Link copied!')
    } catch {
      setStatus('Could not copy link')
    }
    window.setTimeout(() => setStatus(''), 2000)
  }

  if (!shareUrl) {
    return null
  }

  return (
    <div className={`game-share ${className}`.trim()}>
      <button type="button" className="game-detail__toolbar-btn" onClick={handleShare}>
        <IconShare />
        Share
      </button>
      {status ? (
        <p className="game-share__status" role="status">
          {status}
        </p>
      ) : null}
    </div>
  )
}
