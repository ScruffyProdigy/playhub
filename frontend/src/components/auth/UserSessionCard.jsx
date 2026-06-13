import { useEffect, useState } from 'react'
import { logout } from '../../lib/auth'
import { needsProfileSetup } from '../../lib/avatars'
import { fetchSpiritAnimalJourneyEligibility, formatSpiritAnimalJourneyCooldown } from '../../lib/spiritAnimal'
import { useAuth } from './AuthProvider'
import PlayerProfileEditor from '../avatars/PlayerProfileEditor'
import SpiritAnimalFlow from '../avatars/SpiritAnimalFlow'
import PlayerAvatar from '../avatars/PlayerAvatar'

export default function UserSessionCard({ user, compact = false, showProfileActions = true }) {
  const { clearSession } = useAuth()
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')
  const setupRequired = needsProfileSetup(user)
  const [editorOpen, setEditorOpen] = useState(setupRequired)
  const [spiritFlowOpen, setSpiritFlowOpen] = useState(false)
  const [journeyEligibility, setJourneyEligibility] = useState(null)

  useEffect(() => {
    if (!showProfileActions) {
      return undefined
    }
    let cancelled = false
    void fetchSpiritAnimalJourneyEligibility()
      .then((eligibility) => {
        if (!cancelled) {
          setJourneyEligibility(eligibility)
        }
      })
      .catch(() => {
        if (!cancelled) {
          setJourneyEligibility(null)
        }
      })
    return () => {
      cancelled = true
    }
  }, [user?.id, showProfileActions])

  const canBeginSpiritAnimal = journeyEligibility?.canBegin === true
  const eligibilityPending = journeyEligibility === null

  useEffect(() => {
    if (needsProfileSetup(user)) {
      setEditorOpen(true)
    }
  }, [user])

  async function handleLogout() {
    setStatus('loading')
    setError('')

    try {
      await logout()
      clearSession()
    } catch (err) {
      setStatus('error')
      setError(err.message || 'Could not log out')
    } finally {
      setStatus((current) => (current === 'loading' ? 'idle' : current))
    }
  }

  function handleSaved(updated) {
    if (!needsProfileSetup(updated)) {
      setEditorOpen(false)
      setSpiritFlowOpen(false)
    }
  }

  function handleSpiritComplete(updated) {
    handleSaved(updated)
    void fetchSpiritAnimalJourneyEligibility().then(setJourneyEligibility).catch(() => {})
  }

  return (
    <section
      className={`panel-card user-card${compact ? ' user-card--compact' : ''}`}
      aria-labelledby="welcome-heading"
    >
      <div className="user-card__header">
        <PlayerAvatar user={user} size="md" />
        <div>
          <h2 id="welcome-heading">{setupRequired ? 'Set up your display' : 'Welcome back'}</h2>
          <p className="user-email">{user.email}</p>
          {!setupRequired && user.displayName ? <p className="user-name">{user.displayName}</p> : null}
        </div>
      </div>

      {spiritFlowOpen ? (
        <SpiritAnimalFlow
          onComplete={handleSpiritComplete}
          onCancel={() => setSpiritFlowOpen(false)}
        />
      ) : editorOpen ? (
        <PlayerProfileEditor
          user={user}
          required={setupRequired}
          onSaved={handleSaved}
          onCancel={setupRequired ? undefined : () => setEditorOpen(false)}
          onBeginSpiritAnimal={canBeginSpiritAnimal ? () => {
            setEditorOpen(false)
            setSpiritFlowOpen(true)
          } : undefined}
        />
      ) : showProfileActions ? (
        <>
          <button
            type="button"
            className="game-list-button game-list-button-secondary"
            onClick={() => setEditorOpen(true)}
          >
            Change display
          </button>
          {eligibilityPending ? null : canBeginSpiritAnimal ? (
            <button
              type="button"
              className="game-list-button game-list-button-secondary"
              onClick={() => setSpiritFlowOpen(true)}
            >
              Find my spirit animal
            </button>
          ) : (
            <p className="spirit-animal__hint">
              {formatSpiritAnimalJourneyCooldown(
                journeyEligibility?.daysRemaining,
                journeyEligibility?.cooldownEndsAt,
              )}
            </p>
          )}
        </>
      ) : null}

      <button type="button" onClick={handleLogout} disabled={status === 'loading'}>
        {status === 'loading' ? 'Logging out…' : 'Log out'}
      </button>
      {error ? (
        <p className="status-message status-message-error" role="status">
          {error}
        </p>
      ) : null}
    </section>
  )
}
