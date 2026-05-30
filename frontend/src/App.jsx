import CompleteMagicPage from './components/auth/CompleteMagicPage'
import AuthPanel from './components/auth/AuthPanel'
import { AuthProvider } from './components/auth/AuthProvider'
import './App.css'

function HomePage() {
  return (
    <main className="app-shell">
      <header className="app-header">
        <h1>PlayHub</h1>
        <p className="tagline">Your Gaming Hub - Queue, Play, Trade</p>
      </header>

      <AuthPanel />
    </main>
  )
}

function App() {
  return (
    <AuthProvider>
      {window.location.pathname.startsWith('/auth/complete') ? (
        <CompleteMagicPage />
      ) : (
        <HomePage />
      )}
    </AuthProvider>
  )
}

export default App
