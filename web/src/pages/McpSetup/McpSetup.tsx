import React, { useEffect, useMemo, useState } from 'react'
import { Container } from '@harnessio/uicore'
import css from './McpSetup.module.scss'

type Transport = 'stdio' | 'sse'

interface ToolInfo {
  name: string
  description: string
}

const DOMAIN_TOOLS: Record<string, ToolInfo[]> = {
  Pipelines: [
    { name: 'pipeline_list', description: 'List all pipelines in a repository' },
    { name: 'pipeline_create', description: 'Create a new CI/CD pipeline' },
    { name: 'pipeline_trigger', description: 'Trigger a pipeline execution' }
  ],
  Security: [
    { name: 'security_scan', description: 'Run security scan on repository' },
    { name: 'security_findings', description: 'List security findings' },
    { name: 'security_resolve', description: 'Mark a finding as resolved' }
  ],
  'Quality Gates': [
    { name: 'quality_rules_list', description: 'List quality gate rules' },
    { name: 'quality_rules_create', description: 'Create a quality gate rule' },
    { name: 'quality_check', description: 'Run quality gate check' }
  ],
  'Error Tracker': [
    { name: 'errors_list', description: 'List tracked errors' },
    { name: 'errors_resolve', description: 'Resolve an error' },
    { name: 'errors_assign', description: 'Assign an error to a developer' }
  ],
  Remediation: [
    { name: 'remediation_suggest', description: 'Get AI-suggested fixes' },
    { name: 'remediation_apply', description: 'Apply a remediation' },
    { name: 'remediation_status', description: 'Check remediation status' }
  ],
  'Health Monitor': [
    { name: 'health_check', description: 'Run system health checks' },
    { name: 'health_metrics', description: 'Get health metrics' },
    { name: 'health_alerts', description: 'List health alerts' }
  ],
  'Feature Flags': [
    { name: 'flags_list', description: 'List feature flags' },
    { name: 'flags_toggle', description: 'Toggle a feature flag' },
    { name: 'flags_create', description: 'Create a new feature flag' }
  ],
  'Tech Debt': [
    { name: 'debt_list', description: 'List tech debt items' },
    { name: 'debt_prioritize', description: 'Prioritize tech debt items' },
    { name: 'debt_create', description: 'Track a new tech debt item' }
  ]
}

export default function McpSetup() {
  const [activeTransport, setActiveTransport] = useState<Transport>('stdio')
  const [serverStatus, setServerStatus] = useState<'checking' | 'online' | 'offline'>('checking')
  const [copiedBlock, setCopiedBlock] = useState<string | null>(null)
  const origin = typeof window === 'undefined' ? 'http://localhost:3000' : window.location.origin

  useEffect(() => {
    fetch('/api/v1/system/config')
      .then(res => setServerStatus(res.ok ? 'online' : 'offline'))
      .catch(() => setServerStatus('offline'))
  }, [])

  const stdioConfig = useMemo(
    () => `{
  "mcpServers": {
    "solodev": {
      "command": "./solodev",
      "args": ["mcp"],
      "env": {
        "SOLODEV_URL": "${origin}",
        "SOLODEV_TOKEN": "<your-token>"
      }
    }
  }
}`,
    [origin]
  )

  const sseConfig = useMemo(
    () => `{
  "mcpServers": {
    "solodev": {
      "url": "${origin}/api/v1/mcp/sse",
      "headers": {
        "Authorization": "Bearer <your-token>"
      }
    }
  }
}`,
    [origin]
  )

  const cursorConfig = useMemo(
    () => `{
  "mcpServers": {
    "solodev": {
      "command": "./solodev",
      "args": ["mcp"],
      "env": {
        "SOLODEV_URL": "${origin}",
        "SOLODEV_TOKEN": "<your-token>"
      }
    }
  }
}`,
    [origin]
  )

  const activeClaudeConfig = activeTransport === 'stdio' ? stdioConfig : sseConfig

  const copyToClipboard = (text: string, blockId: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedBlock(blockId)
      setTimeout(() => setCopiedBlock(null), 2000)
    })
  }

  return (
    <Container className={css.main}>
      <section className={css.hero}>
        <div className={css.heroCopy}>
          <span className={css.eyebrow}>AI client setup</span>
          <h1 className={css.title}>Plug SoloDev directly into Claude, Cursor, or any MCP-native tool.</h1>
          <p className={css.subtitle}>
            Use stdio for local speed, switch to SSE for shared environments, and keep all 24 SoloDev tools available
            behind one reliable connection profile.
          </p>
          <div className={css.heroChips}>
            <span>24 MCP tools</span>
            <span>8 DevOps domains</span>
            <span>Local + remote transports</span>
          </div>
        </div>

        <div className={css.statusCard}>
          <div className={css.statusRow}>
            <span className={css.statusDot} data-status={serverStatus} />
            <div>
              <div className={css.statusText}>
                {serverStatus === 'checking'
                  ? 'Checking SoloDev transport'
                  : serverStatus === 'online'
                    ? 'SoloDev server online'
                    : 'SoloDev server offline'}
              </div>
              <div className={css.statusMeta}>{origin}</div>
            </div>
          </div>
          <ol className={css.checklist}>
            <li>Pick a connection method that matches your environment.</li>
            <li>Copy the generated config block into your client.</li>
            <li>Paste a token and reconnect to verify the handshake.</li>
          </ol>
          {serverStatus === 'offline' && (
            <p className={css.statusHint}>
              Start the server with <code>./solodev server</code> before testing the client connection.
            </p>
          )}
        </div>
      </section>

      <section className={css.setupGrid}>
        <div className={css.sectionCard}>
          <h2 className={css.sectionTitle}>Connection Method</h2>
          <div className={css.tabs}>
            <button
              type="button"
              className={`${css.tab} ${activeTransport === 'stdio' ? css.activeTab : ''}`}
              onClick={() => setActiveTransport('stdio')}>
              stdio
              <span>Recommended for local work</span>
            </button>
            <button
              type="button"
              className={`${css.tab} ${activeTransport === 'sse' ? css.activeTab : ''}`}
              onClick={() => setActiveTransport('sse')}>
              SSE
              <span>Best for shared or remote hosts</span>
            </button>
          </div>
          <p className={css.transportDesc}>
            {activeTransport === 'stdio'
              ? 'The client launches SoloDev as a subprocess and talks over standard input/output. Use this when the binary and your editor live on the same machine.'
              : 'The client connects to the running SoloDev server over HTTP. Use this when the server is shared, remote, or managed separately from the editor.'}
          </p>
        </div>

        <div className={css.sectionCard}>
          <h2 className={css.sectionTitle}>Verification Targets</h2>
          <div className={css.metricList}>
            <div>
              <span className={css.metricValue}>{activeTransport === 'stdio' ? 'Local' : 'Remote'}</span>
              <span className={css.metricLabel}>Deployment mode</span>
            </div>
            <div>
              <span className={css.metricValue}>{serverStatus === 'online' ? 'Ready' : 'Pending'}</span>
              <span className={css.metricLabel}>Handshake state</span>
            </div>
            <div>
              <span className={css.metricValue}>Bearer</span>
              <span className={css.metricLabel}>Auth pattern</span>
            </div>
          </div>
        </div>
      </section>

      <section className={css.configGrid}>
        <div className={css.section}>
          <h2 className={css.sectionTitle}>Claude Desktop</h2>
          <p className={css.configHint}>
            Add to <code>~/Library/Application Support/Claude/claude_desktop_config.json</code>
          </p>
          <div className={css.codeBlock}>
            <div className={css.codeHeader}>
              <span>claude_desktop_config.json</span>
              <button type="button" className={css.copyButton} onClick={() => copyToClipboard(activeClaudeConfig, 'claude')}>
                {copiedBlock === 'claude' ? 'Copied' : 'Copy'}
              </button>
            </div>
            <pre className={css.code}>{activeClaudeConfig}</pre>
          </div>
        </div>

        <div className={css.section}>
          <h2 className={css.sectionTitle}>Cursor</h2>
          <p className={css.configHint}>
            Add to <code>.cursor/mcp.json</code> in your project root
          </p>
          <div className={css.codeBlock}>
            <div className={css.codeHeader}>
              <span>mcp.json</span>
              <button type="button" className={css.copyButton} onClick={() => copyToClipboard(cursorConfig, 'cursor')}>
                {copiedBlock === 'cursor' ? 'Copied' : 'Copy'}
              </button>
            </div>
            <pre className={css.code}>{cursorConfig}</pre>
          </div>
        </div>
      </section>

      <section className={css.section}>
        <div className={css.catalogHeader}>
          <div>
            <span className={css.catalogEyebrow}>Tool catalog</span>
            <h2 className={css.sectionTitle}>24 tools organized by the work your team actually does.</h2>
          </div>
          <p className={css.catalogSummary}>
            The catalog is already split by operating concern, which makes the MCP surface easier to discover and easier to
            automate against.
          </p>
        </div>

        <div className={css.toolGrid}>
          {Object.entries(DOMAIN_TOOLS).map(([domain, tools]) => (
            <div key={domain} className={css.toolDomain}>
              <h3 className={css.toolDomainTitle}>{domain}</h3>
              {tools.map(tool => (
                <div key={tool.name} className={css.toolItem}>
                  <code className={css.toolName}>{tool.name}</code>
                  <span className={css.toolDesc}>{tool.description}</span>
                </div>
              ))}
            </div>
          ))}
        </div>
      </section>
    </Container>
  )
}
