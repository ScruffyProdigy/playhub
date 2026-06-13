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
    installFlag: '--cursor',
    installHint: 'writes .cursor/mcp.json',
  },
  {
    id: 'claude-code',
    label: 'Claude Code',
    installClient: 'claude',
    configPath: '.mcp.json (created automatically by claude mcp add)',
    installFlag: '--claude',
    installHint: 'runs claude mcp add',
  },
  {
    id: 'claude-desktop',
    label: 'Claude Desktop',
    installClient: 'skill-only',
    configPath: '~/Library/Application Support/Claude/claude_desktop_config.json',
    installFlag: '--skill-only',
    installHint: 'installs the skill only — paste MCP config manually below',
  },
]

const AGENT_VERIFY_PROMPT = 'List my JoinQuest games using the MCP tools.'
const AGENT_START_PROMPT = 'I want to create a game on JoinQuest'

const RESTART_STEPS = {
  cursor: [
    'Fully quit Cursor (Cmd+Q on Mac — don’t just close the window).',
    'Reopen your game project folder.',
    'Start a new Agent chat (not an old thread).',
    `Verify MCP: “${AGENT_VERIFY_PROMPT}”`,
    `Start integrating: “${AGENT_START_PROMPT}”`,
  ],
  'claude-desktop': [
    'Fully quit Claude Desktop and reopen it.',
    'Start a new conversation.',
    'Confirm joinquest-integration appears under MCP servers.',
    `Verify MCP: “${AGENT_VERIFY_PROMPT}”`,
    `Start integrating: “${AGENT_START_PROMPT}”`,
  ],
  'claude-code': [
    'Start a new claude session in your game repo.',
    'Approve joinquest-integration when Claude prompts you.',
    `Verify MCP: “${AGENT_VERIFY_PROMPT}”`,
    `Start integrating: “${AGENT_START_PROMPT}”`,
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
    <div className="developer-mcp__tabs" role="tablist" aria-label="Your AI editor">
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
            Pick your editor below — steps, install flags, and config match{' '}
            <strong>{client.label}</strong>. One script installs the <strong>agent skill</strong> and
            connects <strong>MCP</strong> (Node.js 20+, uses <code>npx</code>) from your{' '}
            <strong>game repo root</strong>.
          </p>

          <div className="developer-mcp__client-picker">
            <ClientTabs clientId={clientId} onChange={setClientId} />
          </div>

          <div className="developer-mcp__keys">
            <h3>1. Create an API key</h3>
            <p className="panel-copy">
              Your editor passes this key to the JoinQuest MCP server so it can call joinquest.cc
              on your behalf. In step 2, you&apos;ll set{' '}
              <code>JOINQUEST_API_KEY</code> to this value before running the install script
              {clientId === 'claude-desktop' ? ' (or paste it into the MCP JSON config)' : ''}.
            </p>
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
                    <strong>Step 2 is ready</strong> — this command sets{' '}
                    <code>JOINQUEST_API_KEY</code> to your new key, then runs the installer for{' '}
                    {client.label}. Copy and run from your game repo:
                  </p>
                  <pre className="developer-mcp__config">{installDevCommand}</pre>
                  <div className="developer-mcp__quick-copy-actions">
                    <button
                      type="button"
                      className="button-primary"
                      onClick={() => void handleCopy(installDevCommand, 'install')}
                    >
                      {copied === 'install' ? 'Copied' : 'Copy install command'}
                    </button>
                    {clientId === 'claude-desktop' ? (
                      <button
                        type="button"
                        className="button-secondary"
                        onClick={() => void handleCopy(configText, 'config')}
                      >
                        {copied === 'config' ? 'Copied' : 'Copy MCP config'}
                      </button>
                    ) : null}
                  </div>
                </div>
              </div>
            ) : (
              <>
                <p className="panel-copy developer-mcp__hint">
                  Generate a key above before running the install command in step 2. Without it,
                  replace <code>lq_dev_PASTE_YOUR_KEY</code> with your saved key.
                </p>
                <button
                  type="button"
                  className="button-primary"
                  disabled={busy}
                  onClick={() => void handleGenerateKey()}
                >
                  {busy ? 'Generating…' : 'Generate API key'}
                </button>
              </>
            )}
          </div>

          <div className="developer-mcp__install">
            <h3>2. Install skill + connect MCP</h3>
            <p className="panel-copy">
              {activeKey ? (
                <>
                  Run from your <strong>game repo root</strong>. Your API key from step 1 is already
                  in <code>JOINQUEST_API_KEY</code> below.
                </>
              ) : (
                <>
                  Run from your <strong>game repo root</strong> after step 1. Set{' '}
                  <code>export JOINQUEST_API_KEY=your-key</code> (or replace the placeholder) before
                  the <code>curl</code> line.
                </>
              )}{' '}
              Script source:{' '}
              <a className="auth-link" href={INSTALL_DEV_SCRIPT_GITHUB} target="_blank" rel="noreferrer">
                read on GitHub
              </a>
              . For {client.label}, <code>{client.installFlag}</code> {client.installHint}.
            </p>
            <ul className="developer-mcp__explainer">
              <li>
                Downloads <code>.agents/skills/joinquest-integration/</code> from the JoinQuest repo
                (markdown instructions for your agent).
              </li>
              {clientId === 'cursor' ? (
                <li>
                  <code>--cursor</code> merges <code>joinquest-integration</code> into{' '}
                  <code>.cursor/mcp.json</code> using your API key.
                </li>
              ) : null}
              {clientId === 'claude-code' ? (
                <li>
                  <code>--claude</code> runs <code>claude mcp add</code> for this project (needs
                  Claude Code installed).
                </li>
              ) : null}
              {clientId === 'claude-desktop' ? (
                <li>
                  <code>--skill-only</code> installs the skill; paste the MCP JSON into{' '}
                  <code>{client.configPath}</code> (manual step below).
                </li>
              ) : null}
              <li>
                MCP runs via <code>npx @joinquest/mcp-integration</code> — talks to joinquest.cc
                only when your agent uses MCP tools.
              </li>
            </ul>
            <p className="panel-copy">
              <strong>Install command</strong>
              {clientId !== 'claude-desktop' ? (
                <>
                  {' '}
                  — line 1 sets your API key; line 2 downloads and runs the script with{' '}
                  <code>{client.installFlag}</code>:
                </>
              ) : (
                <> — installs the skill; put your API key in the MCP JSON config below:</>
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
            <details className="developer-mcp__details">
              <summary>Prefer to download and review first?</summary>
              <pre className="developer-mcp__config">{installInspectCommand}</pre>
              <button
                type="button"
                className="button-secondary"
                onClick={() => void handleCopy(installInspectCommand, 'inspect')}
              >
                {copied === 'inspect' ? 'Copied' : 'Copy review-first commands'}
              </button>
            </details>
            {clientId === 'claude-desktop' ? (
              <details className="developer-mcp__details" open>
                <summary>Claude Desktop MCP config (paste after skill install)</summary>
                <p className="panel-copy">
                  Config file: <code>{client.configPath}</code>
                </p>
                <pre className="developer-mcp__config">{configText}</pre>
                <button
                  type="button"
                  className="button-secondary"
                  onClick={() => void handleCopy(configText, 'config')}
                >
                  {copied === 'config' ? 'Copied' : 'Copy MCP config'}
                </button>
              </details>
            ) : null}
            <details className="developer-mcp__details">
              <summary>Optional: global MCP install (offline / pinned version)</summary>
              <pre className="developer-mcp__config">{INSTALL_GLOBAL_MCP}</pre>
            </details>
          </div>

          <div className="developer-mcp__client">
            <h3>3. Restart {client.label} and test</h3>
            <p className="panel-copy">
              After install, restart so {client.label} picks up the MCP server. Confirm the agent
              can reach JoinQuest, then start integrating.
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
              <strong>Then what?</strong> The agent skill walks through discovery, game API
              implementation, integration checks, and release. You approve anything that goes live.
            </p>
          </div>

          {clientId !== 'claude-desktop' ? (
            <details className="developer-mcp__details developer-mcp__manual">
              <summary>Configure {client.label} manually (skip the install script)</summary>
              <p className="panel-copy">
                Paste this if you prefer not to run the installer. Config file:{' '}
                <code>{client.configPath}</code>
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
                    {copied === 'config' ? 'Copied' : 'Copy Cursor JSON'}
                  </button>
                </>
              )}
            </details>
          ) : null}

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
