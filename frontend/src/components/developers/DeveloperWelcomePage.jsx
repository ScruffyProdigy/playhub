import { useEffect, useState } from 'react'
import { APP_NAME } from '../../lib/brand'
import {
  developerDashboardPath,
  fetchMyGame,
  visibilityLabel,
} from '../../lib/developers'
import { navigateTo } from '../../lib/usePathname'

export default function DeveloperWelcomePage({ gameId }) {
  const [game, setGame] = useState(null)
  const [status, setStatus] = useState('loading')

  useEffect(() => {
    let cancelled = false
    fetchMyGame(gameId)
      .then((item) => {
        if (!cancelled) {
          setGame(item)
          setStatus(item ? 'ready' : 'missing')
        }
      })
      .catch(() => {
        if (!cancelled) {
          setStatus('error')
        }
      })
    return () => {
      cancelled = true
    }
  }, [gameId])

  if (status === 'loading') {
    return (
      <main className="app-shell developer-shell">
        <p className="status-message" role="status">
          Loading…
        </p>
      </main>
    )
  }

  if (status !== 'ready' || !game) {
    return (
      <main className="app-shell developer-shell">
        <h1>Game not found</h1>
        <a className="auth-link" href="/developers">
          Back to developers
        </a>
      </main>
    )
  }

  const connected = game.visibility === 'PRIVATE_TESTING'

  return (
    <main className="app-shell developer-shell">
      <header className="app-header">
        <p className="developer-back">
          <a className="auth-link" href="/">
            ← Back to {APP_NAME}
          </a>
        </p>
        <h1>Nice — your game is registered.</h1>
        <p className="tagline">{game.name} · {visibilityLabel(game.visibility)}</p>
      </header>

      <section className="panel-card">
        {connected ? (
          <p className="panel-copy">
            We connected to your API and synced your game modes. Head to your dashboard to run
            integration checks and spin up a test table.
          </p>
        ) : (
          <p className="panel-copy">
            Your game is saved as a draft. We couldn&apos;t reach your API yet — fix the URL or
            deploy, then run checks from your dashboard.
          </p>
        )}
        <p className="panel-copy">
          <strong>Using Claude, Cursor, or another AI assistant?</strong> Turn on the JoinQuest MCP
          and point your agent at our integration guide — it can run the same checks as the
          dashboard.
        </p>
        <p className="panel-copy">
          <strong>Prefer doing it yourself?</strong> The integration guide on your dashboard works
          for humans — step by step.
        </p>
        <p className="panel-copy">
          When things look good, spin up a <strong>test table</strong> and invite friends into your
          room. They can play before your game goes public.
        </p>
        <button
          type="button"
          className="button-primary"
          onClick={() => navigateTo(developerDashboardPath(game.id))}
        >
          Open developer dashboard
        </button>
      </section>
    </main>
  )
}
