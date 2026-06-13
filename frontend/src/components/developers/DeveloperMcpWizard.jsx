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
    configPath: '~/.cursor/mcp.json or .cursor/mcp.json',
  },
  {
    id: 'claude-desktop',
    label: 'Claude Desktop',
    configPath: '~/Library/Application Support/Claude/claude_desktop_config.json',
  },
  {
    id: 'claude-code',
    label: 'Claude Code',
    configPath: '.mcp.json (created automatically by claude mcp add)',
  },
]

const AGENT_VERIFY_PROMPT = 'List my JoinQuest games using the MCP tools.'
const AGENT_START_PROMPT = 'I want to create a game on JoinQuest'

const RESTART_STEPS = [
  {
    id: 'cursor',
    steps: [
      'Fully quit Cursor (Cmd+Q on Mac — don’t just close the window).',
      'Reopen your game project folder.',
      'Start a new Agent chat (not an old thread).',
      `Verify MCP: “${AGENT_VERIFY_PROMPT}”`,
      `Start integrating: “${AGENT_START_PROMPT}”`,
    ],
  },
  {
    id: 'claude-desktop',
    steps: [
      'Fully quit Claude Desktop and reopen it.',
      'Start a new conversation.',
      'Confirm joinquest-integration appears under MCP servers.',
      `Verify MCP: “${AGENT_VERIFY_PROMPT}”`,
      `Start integrating: “${AGENT_START_PROMPT}”`,
    ],
  },
  {
    id: 'claude-code',
    steps: [
      'Start a new claude session in your game repo.',
      'Approve joinquest-integration when Claude prompts you.',
      `Verify MCP: “${AGENT_VERIFY_PROMPT}”`,
      `Start integrating: “${AGENT_START_PROMPT}”`,
    ],
  },
]

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
  const restartSteps =
    RESTART_STEPS.find((item) => item.id === clientId)?.steps ?? RESTART_STEPS[0].steps

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

  const installCursorCommand = useMemo(
    () => buildInstallDevCommand({ apiKey: activeKey ?? undefined, client: 'cursor' }),
    [activeKey],
  )

  const installClaudeCommand = useMemo(
    () => buildInstallDevCommand({ apiKey: activeKey ?? undefined, client: 'claude' }),
    [activeKey],
  )

  const installDevCommand = useMemo(
    () =>
      activeKey
        ? installCursorCommand
        : buildInstallDevCommand({ apiKey: undefined, client: 'cursor' }),
    [activeKey, installCursorCommand],
  )

  const installInspectCommand = useMemo(
    () => buildInstallDevInspectCommand({ apiKey: activeKey ?? undefined, client: 'cursor' }),
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
            One script installs the <strong>agent skill</strong> and configures <strong>MCP</strong>{' '}
            (Node.js 20+, uses <code>npx</code>). Generate an API key below, then run the install
            command from your <strong>game repo root</strong>.
          </p>

          <div className="developer-mcp__keys">
            <h3>1. Create an API key</h3>
            <p className="panel-copy">
              Your editor starts a small JoinQuest MCP program and passes this key to it. The MCP
              server uses the key to call joinquest.cc on your behalf (list games, run checks, etc.).
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
                    <strong>One-click setup</strong> — run in your game repo (skill + MCP):
                  </p>
                  <div className="developer-mcp__quick-copy-actions">
                    <button
                      type="button"
                      className="button-primary"
                      onClick={() => void handleCopy(installCursorCommand, 'install-cursor')}
                    >
                      {copied === 'install-cursor' ? 'Copied' : 'Copy Cursor install'}
                    </button>
                    <button
                      type="button"
                      className="button-primary"
                      onClick={() => void handleCopy(installClaudeCommand, 'install-claude')}
                    >
                      {copied === 'install-claude' ? 'Copied' : 'Copy Claude install'}
                    </button>
                    <button
                      type="button"
                      className="button-secondary"
                      onClick={() => void handleCopy(cursorConfigText, 'cursor-json')}
                    >
                      {copied === 'cursor-json' ? 'Copied' : 'Copy Cursor JSON only'}
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
            <h3>2. Install skill + connect MCP</h3>
            <p className="panel-copy">
              Run from your <strong>game repo root</strong>. The script is open source —{' '}
              <a className="auth-link" href={INSTALL_DEV_SCRIPT_GITHUB} target="_blank" rel="noreferrer">
                read it on GitHub
              </a>{' '}
              before you run anything.
            </p>
            <ul className="developer-mcp__explainer">
              <li>
                Downloads <code>.agents/skills/joinquest-integration/</code> from the JoinQuest repo
                (markdown instructions for your agent).
              </li>
              <li>
                With <code>--cursor</code>: merges a <code>joinquest-integration</code> entry into{' '}
                <code>.cursor/mcp.json</code> using your API key (local file only).
              </li>
              <li>
                With <code>--claude</code>: runs <code>claude mcp add</code> for this project (needs
                Claude Code installed).
              </li>
              <li>
                MCP runs via <code>npx @joinquest/mcp-integration</code> — talks to joinquest.cc only
                when your agent uses MCP tools. No other network calls.
              </li>
            </ul>
            <p className="panel-copy">
              <strong>Quick install</strong> (pipes the script into your shell — review first if you
              prefer):
            </p>
            <pre className="developer-mcp__config">{installDevCommand}</pre>
            <button
              type="button"
              className="button-secondary"
              onClick={() => void handleCopy(installDevCommand, 'install')}
            >
              {copied === 'install' ? 'Copied' : 'Copy install command'}
            </button>
            <details className="developer-mcp__details">
              <summary>Prefer to download and review first?</summary>
              <p className="panel-copy">
                Saves the script as <code>install-joinquest-dev.sh</code> so you can read it before
                running:
              </p>
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
              <summary>Optional: global MCP install (offline / pinned version)</summary>
              <pre className="developer-mcp__config">{INSTALL_GLOBAL_MCP}</pre>
            </details>
          </div>

          <div className="developer-mcp__client">
            <h3>3. Restart your editor and test</h3>
            <p className="panel-copy">
              After the install script (or manual config below), your editor must <strong>restart</strong>{' '}
              so it picks up the MCP server. Then confirm the agent can reach JoinQuest before you
              start integrating.
            </p>
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
              implementation, integration checks, and release. It uses MCP to list your games, run
              checks, and save catalog metadata — you approve anything that goes live. Already
              registered? The same prompt continues integration for your existing game.
            </p>
          </div>

          <details className="developer-mcp__details developer-mcp__manual">
            <summary>Configure your editor manually (skip the install script)</summary>
            <p className="panel-copy">
              Paste this if you prefer not to run the installer. Same result: your editor launches{' '}
              <code>npx @joinquest/mcp-integration</code> with your API key in{' '}
              <code>JOINQUEST_API_KEY</code>.
            </p>
            <p className="panel-copy">
              Config file: <code>{client.configPath}</code>
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
                    : clientId === 'cursor'
                      ? 'Copy Cursor JSON'
                      : 'Copy MCP config'}
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
