import React, { useEffect, useMemo, useState } from 'react'
import { Container } from '@harnessio/uicore'
import css from './McpSetup.module.scss'

type Transport = 'stdio' | 'sse'

type CatalogItem = {
  surface: 'tool' | 'resource' | 'prompt'
  name?: string
  uri?: string
  domain: string
  description: string
  requires?: string[]
  notes?: string
}

type CatalogSection = {
  tools: CatalogItem[]
  resources: CatalogItem[]
  prompts: CatalogItem[]
}

type CatalogCounts = {
  active_tools: number
  active_resources: number
  active_prompts: number
  blocked_tools: number
  blocked_resources: number
  blocked_prompts: number
  coming_soon_tools: number
  coming_soon_resources: number
  coming_soon_prompts: number
}

type CatalogResponse = {
  server_name: string
  protocol_version: string
  counts: CatalogCounts
  active: CatalogSection
  blocked: CatalogSection
  coming_soon: CatalogSection
}

type CatalogStatus = 'available' | 'blocked' | 'coming-soon'

const emptySection: CatalogSection = { tools: [], resources: [], prompts: [] }

const emptyCatalog: CatalogResponse = {
  server_name: 'solodev',
  protocol_version: '2024-11-05',
  counts: {
    active_tools: 0,
    active_resources: 0,
    active_prompts: 0,
    blocked_tools: 0,
    blocked_resources: 0,
    blocked_prompts: 0,
    coming_soon_tools: 0,
    coming_soon_resources: 0,
    coming_soon_prompts: 0
  },
  active: emptySection,
  blocked: emptySection,
  coming_soon: emptySection
}

export default function McpSetup() {
  const [activeTransport, setActiveTransport] = useState<Transport>('stdio')
  const [serverStatus, setServerStatus] = useState<'checking' | 'online' | 'offline'>('checking')
  const [copiedBlock, setCopiedBlock] = useState<string | null>(null)
  const [catalog, setCatalog] = useState<CatalogResponse>(emptyCatalog)
  const origin = typeof window === 'undefined' ? 'http://localhost:3000' : window.location.origin

  useEffect(() => {
    fetch('/api/v1/system/mcp/catalog')
      .then(async res => {
        if (!res.ok) {
          throw new Error('Catalog unavailable')
        }
        const nextCatalog = (await res.json()) as CatalogResponse
        setCatalog(nextCatalog)
        setServerStatus('online')
      })
      .catch(() => {
        setCatalog(emptyCatalog)
        setServerStatus('offline')
      })
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
      "url": "${origin}/mcp",
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
  const availableCount = catalog.counts.active_tools + catalog.counts.active_resources + catalog.counts.active_prompts
  const blockedCount = catalog.counts.blocked_tools + catalog.counts.blocked_resources + catalog.counts.blocked_prompts
  const comingSoonCount =
    catalog.counts.coming_soon_tools + catalog.counts.coming_soon_resources + catalog.counts.coming_soon_prompts

  const copyToClipboard = (text: string, blockId: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedBlock(blockId)
      setTimeout(() => setCopiedBlock(null), 2000)
    })
  }

  const renderItems = (section: CatalogSection, status: CatalogStatus) => {
    const items = [...section.tools, ...section.resources, ...section.prompts]

    if (items.length === 0) {
      return <div className={css.emptyState}>Nothing in this section yet.</div>
    }

    return items.map(item => {
      const identifier = item.name || item.uri || 'unknown'
      return (
        <div key={`${status}-${identifier}`} className={css.toolItem}>
          <div className={css.toolItemHeader}>
            <code className={css.toolName}>{identifier}</code>
            <span className={css.toolBadge} data-status={status}>
              {status === 'available' ? 'Available now' : status === 'blocked' ? 'Blocked by config' : 'Coming soon'}
            </span>
          </div>
          <span className={css.toolMeta}>
            {item.domain}
            {item.surface === 'resource' ? ' resource' : item.surface === 'prompt' ? ' prompt' : ' tool'}
          </span>
          <span className={css.toolDesc}>{item.description}</span>
          {item.requires && item.requires.length > 0 && (
            <span className={css.toolRequires}>Requires: {item.requires.join(', ')}</span>
          )}
          {item.notes && <span className={css.toolNotes}>{item.notes}</span>}
        </div>
      )
    })
  }

  return (
    <Container className={css.main}>
      <section className={css.hero}>
        <div className={css.heroCopy}>
          <span className={css.eyebrow}>AI client setup</span>
          <h1 className={css.title}>Plug SoloDev directly into Claude, Cursor, or any MCP-native tool.</h1>
          <p className={css.subtitle}>
            The live MCP surface stays honest for real clients, while this page keeps a running roadmap of what is
            available now, what is blocked by config, and what is still coming online.
          </p>
          <div className={css.heroChips}>
            <span>{availableCount} available now</span>
            <span>{blockedCount} blocked by config</span>
            <span>{comingSoonCount} coming soon</span>
          </div>
        </div>

        <div className={css.statusCard}>
          <div className={css.statusRow}>
            <span className={css.statusDot} data-status={serverStatus} />
            <div>
              <div className={css.statusText}>
                {serverStatus === 'checking'
                  ? 'Checking SoloDev MCP catalog'
                  : serverStatus === 'online'
                  ? 'SoloDev MCP catalog online'
                  : 'SoloDev MCP catalog offline'}
              </div>
              <div className={css.statusMeta}>{origin}</div>
            </div>
          </div>
          <ol className={css.checklist}>
            <li>Pick the transport that matches where your editor is running.</li>
            <li>Copy the config block into your MCP client.</li>
            <li>Use the catalog below to see which paths are ready versus still in flight.</li>
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
              Streamable HTTP
              <span>Best for shared or remote hosts</span>
            </button>
          </div>
          <p className={css.transportDesc}>
            {activeTransport === 'stdio'
              ? 'The client launches SoloDev as a subprocess and talks over standard input/output. Use this when the binary and your editor live on the same machine.'
              : 'The client connects to the running SoloDev server over HTTP at the mounted MCP endpoint. Use this when SoloDev is running separately from the editor.'}
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
              <span className={css.metricValue}>{availableCount}</span>
              <span className={css.metricLabel}>Live MCP surfaces</span>
            </div>
            <div>
              <span className={css.metricValue}>{serverStatus === 'online' ? 'Ready' : 'Pending'}</span>
              <span className={css.metricLabel}>Catalog handshake</span>
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
              <button
                type="button"
                className={css.copyButton}
                onClick={() => copyToClipboard(activeClaudeConfig, 'claude')}>
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
            <span className={css.catalogEyebrow}>MCP roadmap</span>
            <h2 className={css.sectionTitle}>Track the real SoloDev MCP surface as each path ships end to end.</h2>
          </div>
          <p className={css.catalogSummary}>
            Active items match the actual MCP registration. Blocked items are implemented but currently unavailable in
            this runtime. Coming soon stays visible here so you can track what still needs backend completion.
          </p>
        </div>

        <div className={css.catalogColumns}>
          <div className={css.toolDomain}>
            <h3 className={css.toolDomainTitle}>Available Now</h3>
            {renderItems(catalog.active, 'available')}
          </div>

          <div className={css.toolDomain}>
            <h3 className={css.toolDomainTitle}>Blocked by Config</h3>
            {renderItems(catalog.blocked, 'blocked')}
          </div>

          <div className={css.toolDomain}>
            <h3 className={css.toolDomainTitle}>Coming Soon</h3>
            {renderItems(catalog.coming_soon, 'coming-soon')}
          </div>
        </div>
      </section>
    </Container>
  )
}
