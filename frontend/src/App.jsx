import CompleteSignInPage from './components/auth/CompleteSignInPage'
import AuthPanel from './components/auth/AuthPanel'
import ActiveQueueBanner from './components/games/ActiveQueueBanner'
import GameLobby from './components/games/GameLobby'
import { AuthProvider, useAuth } from './components/auth/AuthProvider'
import { useActiveQueue } from './components/games/useActiveQueue'
import { APP_NAME, APP_TAGLINE } from './lib/brand'
import './App.css'

function HomePage() {
  const { user, loading: authLoading } = useAuth()
  const { activeQueue, busy, refresh, handleLeave } = useActiveQueue()

  return (
    <main className="app-shell">
      {!authLoading && user ? (
        <ActiveQueueBanner activeQueue={activeQueue} busy={busy} onLeave={handleLeave} />
      ) : null}

      <header className="app-header">
        <h1>{APP_NAME}</h1>
        <p className="tagline">{APP_TAGLINE}</p>
      </header>

      <AuthPanel />
      <GameLobby activeQueue={activeQueue} onQueueChange={refresh} />
    </main>
  )
}

function App() {
  return (
    <AuthProvider>
      {window.location.pathname.startsWith('/auth/complete') ? (
        <CompleteSignInPage />
      ) : (
        <HomePage />
      )}
    </AuthProvider>
  )
}

export default App
