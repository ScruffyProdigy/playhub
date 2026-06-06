import { avatarInitial, resolveUserAvatarUrl } from '../../lib/avatars'
import { displayName } from '../../lib/tables'

export default function PlayerAvatar({ user, size = 'md', className = '', title }) {
  const label = title ?? (user ? displayName(user) : 'Player')
  const sizeClass = size === 'sm' ? 'player-avatar--sm' : 'player-avatar--md'
  const avatarUrl = resolveUserAvatarUrl(user)

  if (avatarUrl) {
    return (
      <img
        src={avatarUrl}
        alt=""
        aria-hidden="true"
        className={`player-avatar ${sizeClass}${className ? ` ${className}` : ''}`}
        title={label}
      />
    )
  }

  return (
    <span
      className={`player-avatar player-avatar--fallback ${sizeClass}${className ? ` ${className}` : ''}`}
      title={label}
      aria-hidden="true"
    >
      {avatarInitial(user)}
    </span>
  )
}
