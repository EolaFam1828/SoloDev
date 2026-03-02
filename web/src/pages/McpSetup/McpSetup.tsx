/*
 * Copyright 2024 Harness, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useEffect, useState } from 'react'
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

const STDIO_CONFIG = `{
  "mcpServers": {
    "solodev": {
      "command": "./solodev",
      "args": ["mcp"],
      "env": {
        "SOLODEV_URL": "http://localhost:3000",
        "SOLODEV_TOKEN": "<your-token>"
      }
    }
  }
}`

const SSE_CONFIG = `{
  "mcpServers": {
    "solodev": {
      "url": "http://localhost:3000/api/v1/mcp/sse",
      "headers": {
        "Authorization": "Bearer <your-token>"
      }
    }
  }
}`

const CURSOR_CONFIG = `{
  "mcpServers": {
    "solodev": {
      "command": "./solodev",
      "args": ["mcp"],
      "env": {
        "SOLODEV_URL": "http://localhost:3000",
        "SOLODEV_TOKEN": "<your-token>"
      }
    }
  }
}`

export default function McpSetup() {
  const [activeTransport, setActiveTransport] = useState<Transport>('stdio')
  const [serverStatus, setServerStatus] = useState<'checking' | 'online' | 'offline'>('checking')
  const [copiedBlock, setCopiedBlock] = useState<string | null>(null)

  useEffect(() => {
    fetch('/api/v1/system/config')
      .then(res => setServerStatus(res.ok ? 'online' : 'offline'))
      .catch(() => setServerStatus('offline'))
  }, [])

  const copyToClipboard = (text: string, blockId: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedBlock(blockId)
      setTimeout(() => setCopiedBlock(null), 2000)
    })
  }

  return (
    <Container className={css.main}>
      <div className={css.header}>
        <h1 className={css.title}>MCP Setup</h1>
        <p className={css.subtitle}>Connect your AI tools to SoloDev's 24 DevOps tools across 8 domains</p>
      </div>

      {/* Connection Status */}
      <div className={css.statusCard}>
        <div className={css.statusRow}>
          <span className={css.statusDot} data-status={serverStatus} />
          <span className={css.statusText}>
            SoloDev Server:{' '}
            {serverStatus === 'checking' ? 'Checking...' : serverStatus === 'online' ? 'Online' : 'Offline'}
          </span>
        </div>
        {serverStatus === 'offline' && (
          <p className={css.statusHint}>Start the server with: <code>./solodev</code></p>
        )}
      </div>

      {/* Transport Tabs */}
      <div className={css.section}>
        <h2 className={css.sectionTitle}>Connection Method</h2>
        <div className={css.tabs}>
          <button
            className={`${css.tab} ${activeTransport === 'stdio' ? css.activeTab : ''}`}
            onClick={() => setActiveTransport('stdio')}>
            stdio (Recommended)
          </button>
          <button
            className={`${css.tab} ${activeTransport === 'sse' ? css.activeTab : ''}`}
            onClick={() => setActiveTransport('sse')}>
            SSE (Server-Sent Events)
          </button>
        </div>

        {activeTransport === 'stdio' ? (
          <div className={css.transportInfo}>
            <p className={css.transportDesc}>
              The stdio transport runs SoloDev as a subprocess. Best for local development with Claude Desktop and Cursor.
            </p>
          </div>
        ) : (
          <div className={css.transportInfo}>
            <p className={css.transportDesc}>
              The SSE transport connects to a running SoloDev server over HTTP. Best for remote servers and shared environments.
            </p>
          </div>
        )}
      </div>

      {/* Config Blocks */}
      <div className={css.section}>
        <h2 className={css.sectionTitle}>Claude Desktop</h2>
        <p className={css.configHint}>
          Add to <code>~/Library/Application Support/Claude/claude_desktop_config.json</code>
        </p>
        <div className={css.codeBlock}>
          <div className={css.codeHeader}>
            <span>claude_desktop_config.json</span>
            <button
              className={css.copyButton}
              onClick={() =>
                copyToClipboard(activeTransport === 'stdio' ? STDIO_CONFIG : SSE_CONFIG, 'claude')
              }>
              {copiedBlock === 'claude' ? 'Copied!' : 'Copy'}
            </button>
          </div>
          <pre className={css.code}>{activeTransport === 'stdio' ? STDIO_CONFIG : SSE_CONFIG}</pre>
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
            <button className={css.copyButton} onClick={() => copyToClipboard(CURSOR_CONFIG, 'cursor')}>
              {copiedBlock === 'cursor' ? 'Copied!' : 'Copy'}
            </button>
          </div>
          <pre className={css.code}>{CURSOR_CONFIG}</pre>
        </div>
      </div>

      {/* Tool Catalog */}
      <div className={css.section}>
        <h2 className={css.sectionTitle}>Tool Catalog</h2>
        <p className={css.subtitle}>24 tools across 8 DevOps domains</p>

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
      </div>
    </Container>
  )
}
