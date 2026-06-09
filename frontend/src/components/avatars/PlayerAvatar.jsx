import { avatarInitial, resolveUserAvatarUrl } from '../../lib/avatars'
import { KING_LABEL } from '../../lib/playerCopy'
import { displayName } from '../../lib/tables'

export default function PlayerAvatar({ user, size = 'md', className = '', title, ring }) {
  const baseLabel = title ?? (user ? displayName(user) : 'Player')
  const label = ring === 'king' ? `${baseLabel} (${KING_LABEL})` : baseLabel
  const sizeClass = size === 'sm' ? 'player-avatar--sm' : 'player-avatar--md'
  const avatarUrl = resolveUserAvatarUrl(user)

  const avatar = avatarUrl ? (
    <img
      src={avatarUrl}
      alt=""
      aria-hidden="true"
      className={`player-avatar ${sizeClass}${className ? ` ${className}` : ''}`}
      title={label}
    />
  ) : (
    <span
      className={`player-avatar player-avatar--fallback ${sizeClass}${className ? ` ${className}` : ''}`}
      title={label}
      aria-hidden="true"
    >
      {avatarInitial(user)}
    </span>
  )

  if (ring === 'king') {
    return (
      <span
        className={`player-avatar-frame player-avatar-frame--king player-avatar-frame--${size}`}
        title={label}
      >
        {avatar}
      </span>
    )
  }

  return avatar
}
