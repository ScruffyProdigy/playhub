import { useEffect, useState } from 'react'
import { fetchEnabledOAuthProviders, startOAuthSignIn } from '../../lib/oauth'
import OAuthProviderIcon, { oauthProviderLabel } from './OAuthProviderIcon'

export default function SocialSignInRow() {
  const [providers, setProviders] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    async function load() {
      try {
        const enabled = await fetchEnabledOAuthProviders()
        if (!cancelled) {
          setProviders(enabled)
        }
      } catch {
        if (!cancelled) {
          setProviders([])
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }
    void load()
    return () => {
      cancelled = true
    }
  }, [])

  if (loading || providers.length === 0) {
    return null
  }

  return (
    <div className="auth-social" aria-label="Social sign-in options">
      <div className="auth-social__icons">
        {providers.map((provider) => (
          <button
            key={provider}
            type="button"
            className="auth-social__icon"
            aria-label={`Continue with ${oauthProviderLabel(provider)}`}
            title={`Continue with ${oauthProviderLabel(provider)}`}
            onClick={() => startOAuthSignIn(provider)}
          >
            <OAuthProviderIcon provider={provider} />
          </button>
        ))}
      </div>
    </div>
  )
}
