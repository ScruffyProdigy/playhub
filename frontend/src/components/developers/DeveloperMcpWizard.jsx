import { useEffect, useMemo, useState } from 'react'
import {
  buildClaudeMcpAddCommand,
  buildInstallDevCommand,
  buildInstallDevInspectCommand,
  buildMcpServerConfig,
  createDeveloperApiKey,
  fetchMyDeveloperApiKeys,
  INSTALL_DEV_SCRIPT_GITHUB,
} from '../../lib/developers'

const CLIENTS = [
  {
    id: 'cursor',
    label: 'Cursor',
    installClient: 'cursor',
    configPath: '~/.cursor/mcp.json or .cursor/mcp.json',
  },
  {
    id: 'claude-code',
    label: 'Claude Code',
    installClient: 'claude',
    configPath: '.mcp.json (created automatically by claude mcp add)',
  },
  {
    id: 'claude-desktop',
    label: 'Claude Desktop',
    installClient: 'claude-desktop',
    configPath: 'Settings → Developer → Edit Config (see manual setup below)',
  },
]

const AGENT_VERIFY_PROMPT = 'List my JoinQuest games using the MCP tools.'
const AGENT_START_PROMPT = 'I want to create a game on JoinQuest'

const RESTART_STEPS = {
  cursor: [
    'Quit Cursor completely (Cmd+Q on Mac — closing the window isn’t enough).',
    'Reopen your game project folder.',
    'Start a fresh Agent chat.',
    `Ask: “${AGENT_VERIFY_PROMPT}”`,
    `Then: “${AGENT_START_PROMPT}”`,
  ],
  'claude-desktop': [
    'Quit Claude Desktop, then open it again.',
    'Start a new conversation.',
    'Check that joinquest-integration appears under MCP servers.',
    `Ask: “${AGENT_VERIFY_PROMPT}”`,
    `Then: “${AGENT_START_PROMPT}”`,
  ],
  'claude-code': [
    'Start a new Claude Code session in your game repo.',
    'Approve joinquest-integration when prompted.',
    `Ask: “${AGENT_VERIFY_PROMPT}”`,
    `Then: “${AGENT_START_PROMPT}”`,
  ],
}

const INSTALL_GLOBAL_MCP =
  'curl -fsSL https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-mcp.sh | sh'

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

function ClientTabs({ clientId, onChange }) {
  return (
    <div className="developer-mcp__tabs" role="tablist" aria-label="Which editor do you use?">
      {CLIENTS.map((item) => (
        <button
          key={item.id}
          type="button"
          role="tab"
          aria-selected={clientId === item.id}
          className={`developer-mcp__tab${clientId === item.id ? ' developer-mcp__tab--active' : ''}`}
          onClick={() => onChange(item.id)}
        >
          {item.label}
        </button>
      ))}
    </div>
  )
}

export default function DeveloperMcpWizard({ defaultExpanded = false }) {
  const [expanded, setExpanded] = useState(defaultExpanded)
  const [clientId, setClientId] = useState('cursor')
  const [apiKeys, setApiKeys] = useState([])
  const [apiKeyValue, setApiKeyValue] = useState('')
  const [showApiKey, setShowApiKey] = useState(false)
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

  const activeKey = apiKeyValue.trim() || null
  const client = CLIENTS.find((item) => item.id === clientId) ?? CLIENTS[0]
  const restartSteps = RESTART_STEPS[clientId] ?? RESTART_STEPS.cursor

  const configText = useMemo(() => {
    const base = buildMcpServerConfig({
      apiKey: activeKey ?? '<paste-api-key-here>',
      clientId,
    })
    const payload = clientId === 'claude-code' ? buildClaudeCodeConfig(base) : base
    return JSON.stringify(payload, null, 2)
  }, [activeKey, clientId])

  const claudeAddCommand = useMemo(
    () =>
      buildClaudeMcpAddCommand({
        apiKey: activeKey ?? '<paste-api-key-here>',
      }),
    [activeKey],
  )

  const installDevCommand = useMemo(
    () =>
      buildInstallDevCommand({
        apiKey: activeKey ?? undefined,
        client: client.installClient,
      }),
    [activeKey, client.installClient],
  )

  const installInspectCommand = useMemo(
    () =>
      buildInstallDevInspectCommand({
        apiKey: activeKey ?? undefined,
        client: client.installClient,
      }),
    [activeKey, client.installClient],
  )

  async function handleGenerateKey() {
    setBusy(true)
    setError('')
    try {
      const result = await createDeveloperApiKey('Integration agent')
      setApiKeyValue(result.secret)
      setShowApiKey(true)
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
        <h2 id="mcp-heading">Connect an AI assistant</h2>
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
          Using Cursor or Claude? Hook up JoinQuest so your agent can run checks and save metadata
          — no browser cookies required.
        </p>
      ) : (
        <>
          <div className="developer-mcp__client-picker">
            <ClientTabs clientId={clientId} onChange={setClientId} />
          </div>

          <p className="panel-copy">
            Three quick steps for <strong>{client.label}</strong>. Run the install command from your
            game repo — we&apos;ll add the agent skill and connect MCP.
          </p>

          <div className="developer-mcp__keys">
            <h3>1. Get an API key</h3>
            <p className="panel-copy">
              This tells JoinQuest it&apos;s your agent calling. Paste a saved key below, or
              generate a new one — either way, step 2 fills in the install command for you.
            </p>
            {apiKeys.length > 0 ? (
              <p className="panel-copy">
                You already have a key (<code>{apiKeys[0].keyPrefix}…</code>
                {apiKeys[0].lastUsedAt ? ', used recently' : ', not used yet'}). Paste it if you
                saved it, or generate a new one.
              </p>
            ) : null}
            <div className="developer-form__field developer-mcp__key-field">
              <label htmlFor="developer-mcp-api-key">Your API key</label>
              <div className="developer-mcp__key-input-row">
                <input
                  id="developer-mcp-api-key"
                  type={showApiKey ? 'text' : 'password'}
                  autoComplete="off"
                  spellCheck={false}
                  placeholder="lq_dev_…"
                  value={apiKeyValue}
                  onChange={(event) => setApiKeyValue(event.target.value)}
                />
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => setShowApiKey((value) => !value)}
                  disabled={!apiKeyValue}
                >
                  {showApiKey ? 'Hide' : 'Show'}
                </button>
              </div>
              <p className="developer-form__hint">
                We only use this in your browser to build the copy-paste commands below — nothing
                gets sent until you run them.
              </p>
            </div>
            <div className="developer-mcp__key-actions">
              <button
                type="button"
                className="button-primary"
                disabled={busy}
                onClick={() => void handleGenerateKey()}
              >
                {busy ? 'Generating…' : 'Generate new API key'}
              </button>
              {activeKey ? (
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => void handleCopy(activeKey, 'key')}
                >
                  {copied === 'key' ? 'Copied' : 'Copy API key'}
                </button>
              ) : null}
            </div>
          </div>

          <div className="developer-mcp__install">
            <h3>2. Install the skill and MCP</h3>
            <p className="panel-copy">
              From your <strong>game repo root</strong>, run the command below.
              {activeKey ? (
                <> Your key is already in the first line — copy and run as-is.</>
              ) : (
                <>
                  {' '}
                  Paste or generate your key in step 1 first — we&apos;ll fill in{' '}
                  <code>JOINQUEST_API_KEY</code> for you.
                </>
              )}
            </p>
            <pre className="developer-mcp__config">{installDevCommand}</pre>
            <button
              type="button"
              className="button-secondary"
              onClick={() => void handleCopy(installDevCommand, 'install-main')}
            >
              {copied === 'install-main' ? 'Copied' : 'Copy install command'}
            </button>
            <p className="panel-copy developer-mcp__what-it-does">What the script does:</p>
            <ul className="developer-mcp__explainer">
              <li>
                Adds <code>.agents/skills/joinquest-integration/</code> — step-by-step instructions
                for your agent.
              </li>
              {clientId === 'cursor' ? (
                <li>
                  With <code>--cursor</code>, wires MCP into <code>.cursor/mcp.json</code>.
                </li>
              ) : null}
              {clientId === 'claude-code' ? (
                <li>
                  With <code>--claude</code>, runs <code>claude mcp add</code> for this project
                  (needs Claude Code installed).
                </li>
              ) : null}
              {clientId === 'claude-desktop' ? (
                <li>
                  With <code>--claude-desktop</code>, merges MCP into Claude Desktop&apos;s config
                  (macOS, Windows, or Linux).
                </li>
              ) : null}
              <li>
                MCP runs via <code>npx @joinquest/mcp-integration</code> when your agent needs it.
              </li>
            </ul>
            <p className="panel-copy">
              <a className="auth-link" href={INSTALL_DEV_SCRIPT_GITHUB} target="_blank" rel="noreferrer">
                Read the script on GitHub
              </a>{' '}
              before you run it, if you like.
            </p>
            <details className="developer-mcp__details">
              <summary>Want to download and review first?</summary>
              <pre className="developer-mcp__config">{installInspectCommand}</pre>
              <button
                type="button"
                className="button-secondary"
                onClick={() => void handleCopy(installInspectCommand, 'inspect')}
              >
                {copied === 'inspect' ? 'Copied' : 'Copy review-first commands'}
              </button>
            </details>
            <details className="developer-mcp__details">
              <summary>Offline or pinned install</summary>
              <p className="panel-copy">Optional — if you don&apos;t want to use npx each time:</p>
              <pre className="developer-mcp__config">{INSTALL_GLOBAL_MCP}</pre>
            </details>
          </div>

          <div className="developer-mcp__client">
            <h3>3. Restart and try it</h3>
            <p className="panel-copy">
              {client.label} needs a full restart to pick up MCP. Then check that your agent can
              reach JoinQuest:
            </p>
            <ol className="developer-mcp__steps">
              {restartSteps.map((step) => (
                <li key={step}>{step}</li>
              ))}
            </ol>
            <div className="developer-mcp__quick-copy-actions">
              <button
                type="button"
                className="button-primary"
                onClick={() => void handleCopy(AGENT_START_PROMPT, 'start-prompt')}
              >
                {copied === 'start-prompt' ? 'Copied' : 'Copy start prompt'}
              </button>
            </div>
            <p className="panel-copy">
              From there, the agent skill walks you through discovery, integration checks, and
              release. You approve anything that goes live.
            </p>
          </div>

          <details className="developer-mcp__details developer-mcp__manual">
            <summary>Manual setup (skip the script)</summary>
            <p className="panel-copy">
              Prefer to paste config yourself?{' '}
              {clientId === 'claude-desktop' ? (
                <>
                  In Claude Desktop: <strong>Settings → Developer → Edit Config</strong>. Or edit
                  the file directly (macOS:{' '}
                  <code>~/Library/Application Support/Claude/claude_desktop_config.json</code>
                  ).
                </>
              ) : (
                <>File: <code>{client.configPath}</code></>
              )}
            </p>
            {clientId === 'claude-code' ? (
              <>
                <pre className="developer-mcp__config">{claudeAddCommand}</pre>
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => void handleCopy(claudeAddCommand, 'claude-add')}
                >
                  {copied === 'claude-add' ? 'Copied' : 'Copy Claude command'}
                </button>
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
                    : clientId === 'claude-desktop'
                      ? 'Copy MCP config'
                      : 'Copy Cursor JSON'}
                </button>
              </>
            )}
          </details>

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
