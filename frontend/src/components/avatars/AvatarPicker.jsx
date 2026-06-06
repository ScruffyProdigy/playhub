import { useEffect, useState } from 'react'
import { fetchStarterAvatars, selectStarterAvatar } from '../../lib/avatars'
import { useAuth } from '../auth/AuthProvider'

export default function AvatarPicker({ currentKey, onSelected }) {
  const { acceptSessionUser } = useAuth()
  const [options, setOptions] = useState([])
  const [busyKey, setBusyKey] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    void fetchStarterAvatars().then((items) => {
      if (!cancelled) {
        setOptions(items)
      }
    })
    return () => {
      cancelled = true
    }
  }, [])

  async function handleSelect(key) {
    if (busyKey || key === currentKey) {
      return
    }
    setBusyKey(key)
    setError('')
    try {
      const user = await selectStarterAvatar(key)
      acceptSessionUser(user)
      onSelected?.(user)
    } catch (err) {
      setError(err.message || 'Could not save avatar.')
    } finally {
      setBusyKey('')
    }
  }

  return (
    <div className="avatar-picker">
      <p className="avatar-picker__lead">Choose your journey icon</p>
      <ul className="avatar-picker__grid" role="list">
        {options.map((option) => {
          const selected = option.key === currentKey
          return (
            <li key={option.key}>
              <button
                type="button"
                className={`avatar-picker__option${selected ? ' avatar-picker__option--selected' : ''}`}
                disabled={Boolean(busyKey)}
                aria-pressed={selected}
                aria-label={`${option.name}${selected ? ' (selected)' : ''}`}
                onClick={() => handleSelect(option.key)}
              >
                <img src={option.imageUrl} alt="" className="avatar-picker__image" />
                <span className="avatar-picker__name">{option.name}</span>
              </button>
            </li>
          )
        })}
      </ul>
      {error ? <p className="status-message status-message-error">{error}</p> : null}
    </div>
  )
}
