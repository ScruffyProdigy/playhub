import { useEffect, useMemo, useState } from 'react'
import {
  buildClaudeMcpAddCommand,
  buildMcpServerConfig,
  createDeveloperApiKey,
  fetchMyDeveloperApiKeys,
} from '../../lib/developers'

const CLIENTS = [
  {
    id: 'cursor',
    label: 'Cursor',
    steps: [
      'Generate an API key below (shown once).',
      'Copy the MCP config below into ~/.cursor/mcp.json or .cursor/mcp.json (uses npx — Node.js 20+).',
      'Fully quit Cursor (Cmd+Q), reopen the project, and start a new Agent chat.',
      'Verify: ask the agent to call joinquest_integration_list_my_games.',
    ],
    configPath: '~/.cursor/mcp.json or .cursor/mcp.json',
  },
  {
    id: 'claude-desktop',
    label: 'Claude Desktop',
    steps: [
      'Generate an API key below (shown once).',
      'Open Claude Desktop → Settings → Developer → Edit Config.',
      'Paste the config and fully quit + reopen Claude Desktop.',
    ],
    configPath: '~/Library/Application Support/Claude/claude_desktop_config.json',
  },
  {
    id: 'claude-code',
    label: 'Claude Code',
    steps: [
      'Generate an API key below (shown once).',
      'Run the claude mcp add command below in your game repo (uses npx — Node.js 20+).',
      'Start a new claude session and approve joinquest-integration when prompted.',
    ],
    configPath: '.mcp.json (created automatically by claude mcp add)',
  },
]

const SKILL_GITHUB =
  'https://github.com/scruffyprodigy/playhub/tree/main/.agents/skills/joinquest-integration'
const INSTALL_SKILL = `mkdir -p .agents/skills
cp -r /path/to/lobby/.agents/skills/joinquest-integration .agents/skills/`
const INSTALL_CLI = 'curl -fsSL https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-mcp.sh | sh'

function buildClaudeCodeConfig(base) {
  const server = base.mcpServers['joinquest-integration']
  return {
    mcpServers: {
      'joinquest-integration': {
        type: 'stdio',
        ...server,
      },
    },
  }
}

export default function DeveloperMcpWizard({ defaultExpanded = false }) {
  const [expanded, setExpanded] = useState(defaultExpanded)
  const [clientId, setClientId] = useState('cursor')
  const [apiKeys, setApiKeys] = useState([])
  const [newSecret, setNewSecret] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [copied, setCopied] = useState('')

  useEffect(() => {
    let cancelled = false
    fetchMyDeveloperApiKeys()
      .then((keys) => {
        if (!cancelled) {
          setApiKeys(keys)
        }
      })
      .catch(() => {
        // ignore — user may not be signed in yet
      })
    return () => {
      cancelled = true
    }
  }, [])

  const activeKey = newSecret || null
  const client = CLIENTS.find((item) => item.id === clientId) ?? CLIENTS[0]

  const configText = useMemo(() => {
    const base = buildMcpServerConfig({
      apiKey: activeKey ?? '<paste-api-key-here>',
      clientId,
    })
    const payload = clientId === 'claude-code' ? buildClaudeCodeConfig(base) : base
    return JSON.stringify(payload, null, 2)
  }, [activeKey, clientId])

  const cursorConfigText = useMemo(
    () =>
      JSON.stringify(
        buildMcpServerConfig({
          apiKey: activeKey ?? '<paste-api-key-here>',
          clientId: 'cursor',
        }),
        null,
        2,
      ),
    [activeKey],
  )

  const claudeAddCommand = useMemo(
    () =>
      buildClaudeMcpAddCommand({
        apiKey: activeKey ?? '<paste-api-key-here>',
      }),
    [activeKey],
  )

  async function handleGenerateKey() {
    setBusy(true)
    setError('')
    try {
      const result = await createDeveloperApiKey('Integration agent')
      setNewSecret(result.secret)
      setApiKeys((prev) => [result.apiKey, ...prev])
    } catch (err) {
      setError(err.message || 'Could not create API key.')
    } finally {
      setBusy(false)
    }
  }

  async function handleCopy(text, label) {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(label)
      setTimeout(() => setCopied(''), 2000)
    } catch {
      // ignore
    }
  }

  return (
    <section className="panel-card developer-mcp" aria-labelledby="mcp-heading">
      <div className="developer-mcp__header">
        <h2 id="mcp-heading">Connect an AI assistant (optional)</h2>
        <button
          type="button"
          className="button-secondary developer-mcp__toggle"
          onClick={() => setExpanded((v) => !v)}
          aria-expanded={expanded}
        >
          {expanded ? 'Hide' : 'Show setup'}
        </button>
      </div>
      {!expanded ? (
        <p className="panel-copy">
          Using an AI coding agent? Add the JoinQuest <strong>agent skill</strong> to your game repo,
          then connect MCP so it can run checks and save metadata — no DevTools cookies required.
        </p>
      ) : (
        <>
          <p className="panel-copy">
            <strong>Agent skill</strong> (Claude Code, Codex, Cursor, Copilot, and others) walks your
            agent through discovery → API → checks → release. <strong>MCP</strong> lets it call
            JoinQuest on your behalf. Use both for the best vibe-coding experience.
          </p>

          <div className="developer-mcp__skill">
            <h3>Agent skill (recommended first)</h3>
            <p className="panel-copy">
              Copy{' '}
              <code>.agents/skills/joinquest-integration/</code> from the{' '}
              <a className="auth-link" href={SKILL_GITHUB} target="_blank" rel="noreferrer">
                JoinQuest repo
              </a>{' '}
              into your <strong>game project</strong> at the same path. Your agent loads it automatically
              when you mention JoinQuest integration.
            </p>
            <pre className="developer-mcp__config">{INSTALL_SKILL}</pre>
            <button
              type="button"
              className="button-secondary"
              onClick={() => void handleCopy(INSTALL_SKILL, 'skill')}
            >
              {copied === 'skill' ? 'Copied' : 'Copy install commands'}
            </button>
          </div>

          <div className="developer-mcp__keys">
            <h3>1. API key</h3>
            {apiKeys.length > 0 && !newSecret ? (
              <p className="panel-copy">
                Active key: <code>{apiKeys[0].keyPrefix}…</code>
                {apiKeys[0].lastUsedAt ? ' (used recently)' : ' (not used yet)'}
              </p>
            ) : null}
            {newSecret ? (
              <div className="developer-mcp__secret" role="status">
                <p className="panel-copy">
                  <strong>Copy this key now</strong> — it won&apos;t be shown again.
                </p>
                <code>{newSecret}</code>
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => void handleCopy(newSecret, 'key')}
                >
                  {copied === 'key' ? 'Copied' : 'Copy API key'}
                </button>
                <div className="developer-mcp__quick-copy">
                  <p className="panel-copy">
                    <strong>One-click setup</strong> — configs below include your new key:
                  </p>
                  <div className="developer-mcp__quick-copy-actions">
                    <button
                      type="button"
                      className="button-primary"
                      onClick={() => void handleCopy(cursorConfigText, 'cursor-json')}
                    >
                      {copied === 'cursor-json' ? 'Copied' : 'Copy Cursor JSON'}
                    </button>
                    <button
                      type="button"
                      className="button-primary"
                      onClick={() => void handleCopy(claudeAddCommand, 'claude-cmd')}
                    >
                      {copied === 'claude-cmd' ? 'Copied' : 'Copy Claude command'}
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <button
                type="button"
                className="button-primary"
                disabled={busy}
                onClick={() => void handleGenerateKey()}
              >
                {busy ? 'Generating…' : 'Generate API key'}
              </button>
            )}
          </div>

          <div className="developer-mcp__install">
            <h3>2. Node.js</h3>
            <p className="panel-copy">
              Requires Node.js 20+. The config below uses <code>npx</code> to run{' '}
              <code>@joinquest/mcp-integration</code> — no separate install step.
            </p>
            <details className="developer-mcp__details">
              <summary>Optional: global install (offline / pinned version)</summary>
              <pre className="developer-mcp__config">{INSTALL_CLI}</pre>
              <button
                type="button"
                className="button-secondary"
                onClick={() => void handleCopy(INSTALL_CLI, 'install')}
              >
                {copied === 'install' ? 'Copied' : 'Copy install script'}
              </button>
            </details>
          </div>

          <div className="developer-mcp__client">
            <h3>3. Add to your editor</h3>
            <div className="developer-mcp__tabs" role="tablist" aria-label="MCP client">
              {CLIENTS.map((item) => (
                <button
                  key={item.id}
                  type="button"
                  role="tab"
                  aria-selected={clientId === item.id}
                  className={`developer-mcp__tab${clientId === item.id ? ' developer-mcp__tab--active' : ''}`}
                  onClick={() => setClientId(item.id)}
                >
                  {item.label}
                </button>
              ))}
            </div>
            <ol className="developer-mcp__steps">
              {client.steps.map((step) => (
                <li key={step}>{step}</li>
              ))}
            </ol>
            <p className="panel-copy">
              Config file: <code>{client.configPath}</code>
            </p>
            {clientId === 'claude-code' ? (
              <>
                <p className="panel-copy">
                  <strong>Recommended</strong> — from your game repo:
                </p>
                <pre className="developer-mcp__config">{claudeAddCommand}</pre>
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => void handleCopy(claudeAddCommand, 'claude-add')}
                >
                  {copied === 'claude-add' ? 'Copied' : 'Copy Claude command'}
                </button>
                <p className="panel-copy">
                  Or run the helper (install + prompt for API key):{' '}
                  <code>bash scripts/add-joinquest-mcp-claude.sh</code>
                </p>
                <details className="developer-mcp__details">
                  <summary>Manual .mcp.json instead</summary>
                  <pre className="developer-mcp__config">{configText}</pre>
                  <button
                    type="button"
                    className="button-secondary"
                    onClick={() => void handleCopy(configText, 'config-manual')}
                  >
                    {copied === 'config-manual' ? 'Copied' : 'Copy .mcp.json'}
                  </button>
                </details>
              </>
            ) : (
              <>
                <pre className="developer-mcp__config">{configText}</pre>
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => void handleCopy(configText, 'config')}
                >
                  {copied === 'config'
                    ? 'Copied'
                    : clientId === 'cursor'
                      ? 'Copy Cursor JSON'
                      : 'Copy MCP config'}
                </button>
              </>
            )}
          </div>

          {error ? (
            <p className="status-message status-message-error" role="alert">
              {error}
            </p>
          ) : null}
        </>
      )}
    </section>
  )
}
