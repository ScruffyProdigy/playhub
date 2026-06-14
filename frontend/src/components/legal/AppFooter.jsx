export default function AppFooter() {
  return (
    <footer className="app-footer">
      <nav className="app-footer__nav" aria-label="Legal">
        <a className="app-footer__link" href="/terms">
          Terms of Service
        </a>
        <span className="app-footer__sep" aria-hidden="true">
          ·
        </span>
        <a className="app-footer__link" href="/privacy">
          Privacy Policy
        </a>
      </nav>
    </footer>
  )
}
