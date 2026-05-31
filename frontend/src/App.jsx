import CompleteSignInPage from './components/auth/CompleteSignInPage'
import AuthPanel from './components/auth/AuthPanel'
import GameLobby from './components/games/GameLobby'
import { AuthProvider } from './components/auth/AuthProvider'
import { APP_NAME, APP_TAGLINE } from './lib/brand'
import './App.css'

function HomePage() {
  return (
    <main className="app-shell">
      <header className="app-header">
        <h1>{APP_NAME}</h1>
        <p className="tagline">{APP_TAGLINE}</p>
      </header>

      <AuthPanel />
      <GameLobby />
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
