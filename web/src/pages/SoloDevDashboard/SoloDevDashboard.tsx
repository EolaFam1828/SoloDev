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

import React, { useEffect, useMemo, useState } from 'react'
import { useHistory } from 'react-router-dom'
import { Container } from '@harnessio/uicore'
import { useAppContext } from 'AppContext'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import { DomainCard, type DomainStatus } from './DomainCard'
import css from './SoloDevDashboard.module.scss'

interface DomainData {
  icon: string
  title: string
  summary: string
  metric: string
  status: DomainStatus
  statusLabel: string
  accentColor: string
  capability: string
  actionLabel: string
  key: string
}

const DOMAIN_DATA: DomainData[] = [
  {
    icon: '\u25B6',
    title: 'Pipelines',
    summary: 'Keep delivery visible with pipeline inventory, trigger health, and execution handoffs in one lane.',
    metric: '3 active flows',
    status: 'healthy',
    statusLabel: 'Healthy',
    accentColor: 'var(--solodev-pipeline-blue)',
    capability: 'Ship velocity',
    actionLabel: 'Open pipelines',
    key: 'pipelines'
  },
  {
    icon: '\u26A1',
    title: 'Security',
    summary: 'Surface new findings fast, keep remediation moving, and keep risky repos from drifting quietly.',
    metric: '12 active findings',
    status: 'warning',
    statusLabel: 'Warning',
    accentColor: 'var(--solodev-security-red)',
    capability: 'Risk posture',
    actionLabel: 'Stay in command',
    key: 'security'
  },
  {
    icon: '\u2713',
    title: 'Quality Gates',
    summary: 'Translate team standards into enforceable gates and keep release confidence visible at a glance.',
    metric: '5 guardrails live',
    status: 'passing',
    statusLabel: 'Passing',
    accentColor: 'var(--solodev-quality-green)',
    capability: 'Release discipline',
    actionLabel: 'Inspect rules',
    key: 'quality'
  },
  {
    icon: '\u2715',
    title: 'Error Tracker',
    summary: 'Turn incidents into structured work instead of scattered screenshots, logs, and Slack archaeology.',
    metric: '2 unresolved incidents',
    status: 'warning',
    statusLabel: 'Warning',
    accentColor: 'var(--solodev-error-orange)',
    capability: 'Production signal',
    actionLabel: 'Trace issues',
    key: 'errors'
  },
  {
    icon: '\u2692',
    title: 'Remediation',
    summary: 'Queue fixes, route AI-assisted repairs, and keep operators focused on the highest leverage work.',
    metric: '1 patch waiting',
    status: 'active',
    statusLabel: 'Active',
    accentColor: 'var(--solodev-remediation-cyan)',
    capability: 'Action engine',
    actionLabel: 'Review queue',
    key: 'remediation'
  },
  {
    icon: '\u2665',
    title: 'Health Monitor',
    summary: 'Track platform heartbeat, repo readiness, and route-level health without leaving the control surface.',
    metric: 'All systems steady',
    status: 'healthy',
    statusLabel: 'Healthy',
    accentColor: 'var(--solodev-health-emerald)',
    capability: 'Service heartbeat',
    actionLabel: 'See health',
    key: 'health'
  },
  {
    icon: '\u2691',
    title: 'Feature Flags',
    summary: 'Ship progressively, keep experiments legible, and flip features without turning rollouts into guesswork.',
    metric: '8 live flags',
    status: 'active',
    statusLabel: 'Active',
    accentColor: 'var(--solodev-flags-purple)',
    capability: 'Progressive delivery',
    actionLabel: 'View flags',
    key: 'flags'
  },
  {
    icon: '\u2630',
    title: 'Tech Debt',
    summary: 'Keep the backlog honest with visible debt inventory, prioritization pressure, and sprint-ready context.',
    metric: '15 tracked items',
    status: 'warning',
    statusLabel: 'Medium',
    accentColor: 'var(--solodev-debt-amber)',
    capability: 'Sustainable velocity',
    actionLabel: 'Triage debt',
    key: 'debt'
  }
]

export default function SoloDevDashboard() {
  const { routes } = useAppContext()
  const history = useHistory()
  const { space, repoMetadata } = useGetRepositoryMetadata()
  const [mcpConnected, setMcpConnected] = useState<boolean | null>(null)
  const repoPath = repoMetadata?.path || ''

  useEffect(() => {
    fetch('/api/v1/system/config')
      .then(res => {
        setMcpConnected(res.ok)
      })
      .catch(() => {
        setMcpConnected(false)
      })
  }, [])

  const headlineStats = useMemo(
    () => [
      { label: 'Domains', value: '8' },
      { label: 'AI Tools', value: '24' },
      { label: 'Resources', value: '8' },
      { label: 'Project', value: space || 'SoloDev' }
    ],
    [space]
  )

  return (
    <Container className={css.main}>
      <section className={css.hero}>
        <div className={css.heroCopy}>
          <span className={css.eyebrow}>SoloDev mission control</span>
          <h1 className={css.title}>A sharper command surface for your entire delivery system.</h1>
          <p className={css.subtitle}>
            Pipelines, quality, security, remediation, and MCP tooling should feel like one instrument panel, not a pile
            of tabs. SoloDev is your high-signal control plane.
          </p>
          <div className={css.heroActions}>
            <button
              className={css.primaryAction}
              onClick={() => {
                if (space && routes.toSOLODEVMcpSetup) {
                  history.push(routes.toSOLODEVMcpSetup({ space }))
                }
              }}>
              Connect an AI client
            </button>
            <button
              className={css.secondaryAction}
              onClick={() => {
                if (space) {
                  history.push(routes.toCODESpaceSettings({ space }))
                }
              }}>
              Project settings
            </button>
            {repoPath && (
              <button className={css.secondaryAction} onClick={() => history.push(routes.toCODEPipelines({ repoPath }))}>
                Pipeline inventory
              </button>
            )}
          </div>
        </div>

        <div className={css.heroPanel}>
          <div className={css.statusBadge}>
            <span className={css.statusDot} data-connected={mcpConnected === true} />
            <span>
              {mcpConnected === true ? 'MCP online and ready' : mcpConnected === false ? 'MCP unreachable' : 'Checking MCP reachability'}
            </span>
          </div>
          <div className={css.statGrid}>
            {headlineStats.map(stat => (
              <div key={stat.label} className={css.statCard}>
                <span className={css.statValue}>{stat.value}</span>
                <span className={css.statLabel}>{stat.label}</span>
              </div>
            ))}
          </div>
          <div className={css.heroNote}>
            SoloDev should feel like a flagship product: warm light mode when you want clarity, a deep command-room dark
            theme when you want focus, and zero leftover platform branding.
          </div>
        </div>
      </section>

      <section className={css.sectionHeader}>
        <div>
          <span className={css.sectionEyebrow}>Domains</span>
          <h2 className={css.sectionTitle}>Everything that drives shipping, quality, and stability.</h2>
        </div>
        <p className={css.sectionSummary}>
          Each lane stays visible, status-first, and ready for the next expansion without feeling like placeholder filler.
        </p>
      </section>

      <div className={css.grid}>
        {DOMAIN_DATA.map(domain => (
          <DomainCard
            key={domain.key}
            icon={domain.icon}
            title={domain.title}
            summary={domain.summary}
            metric={domain.metric}
            status={domain.status}
            statusLabel={domain.statusLabel}
            accentColor={domain.accentColor}
            capability={domain.capability}
            actionLabel={domain.actionLabel}
            onClick={() => {
              if (domain.key === 'pipelines' && repoPath) {
                history.push(routes.toCODEPipelines({ repoPath }))
                return
              }

              if (space) {
                history.push(routes.toCODESpaceSettings({ space }))
              }
            }}
          />
        ))}
      </div>
    </Container>
  )
}
