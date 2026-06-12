import AuthPanel from '../auth/AuthPanel'
import { APP_NAME } from '../../lib/brand'
import RegisterGameForm from './RegisterGameForm'

export default function DeveloperLandingPage() {
  return (
    <main className="app-shell developer-shell">
      <header className="app-header">
        <p className="developer-back">
          <a className="auth-link" href="/">
            ← Back to {APP_NAME}
          </a>
        </p>
        <h1>Have an idea for a multiplayer game?</h1>
      </header>

      <section className="panel-card developer-intro">
        <p className="panel-copy">
          JoinQuest handles the boring stuff — player accounts, rooms, tables, matchmaking, finding
          players — so you can focus on building the game.
        </p>
        <p className="panel-copy">
          You&apos;ll need some basic web dev experience and a <strong>public URL</strong> where your
          game can run. That&apos;s it.
        </p>
        <p className="panel-copy">
          Register below to get started. You&apos;re not committing to anything, and your game
          won&apos;t show up for other players until <em>you</em> say it&apos;s ready (and we&apos;ve
          had a quick look).
        </p>
      </section>

      <AuthPanel />

      <section className="panel-card" aria-labelledby="register-game-heading">
        <h2 id="register-game-heading">Register your game</h2>
        <RegisterGameForm />
      </section>
    </main>
  )
}
