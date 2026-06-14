/** Split wordmark: Join (light) + Quest (indigo). */
export default function AppTitle({ className = '' }) {
  const classes = ['app-title', className].filter(Boolean).join(' ')
  return (
    <h1 className={classes} aria-label="JoinQuest">
      <span className="app-title__join" aria-hidden="true">
        Join
      </span>
      <span className="app-title__quest" aria-hidden="true">
        Quest
      </span>
    </h1>
  )
}
