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
import { type TypesSoloDevOverview, useGetSoloDevOverview } from 'services/code'
import { DomainCard, type DomainStatus } from './DomainCard'
import css from './SoloDevDashboard.module.scss'

interface DomainData {
  icon: string
  title: string
  metric: string
  status: DomainStatus
  statusLabel: string
  accentColor: string
  key: string
}

interface RecentRemediation {
  identifier: string
  title: string
  status: string
  confidence?: number
  pr_link?: string
}

export default function SoloDevDashboard() {
  const { routes } = useAppContext()
  const history = useHistory()
  const { space } = useGetRepositoryMetadata()
  const [mcpConnected, setMcpConnected] = useState<boolean | null>(null)
  const [recentRemediations, setRecentRemediations] = useState<RecentRemediation[]>([])
  const [recentLoading, setRecentLoading] = useState(false)
  const { data: overview } = useGetSoloDevOverview({
    space_ref: space || '',
    lazy: !space
  })

  useEffect(() => {
    fetch('/api/v1/system/config')
      .then(res => {
        setMcpConnected(res.ok)
      })
      .catch(() => {
        setMcpConnected(false)
      })
  }, [])

  useEffect(() => {
    if (!space) {
      setRecentRemediations([])
      return
    }

    const controller = new AbortController()
    setRecentLoading(true)

    fetch(`/api/v1/spaces/${space}/remediations?limit=3`, { signal: controller.signal })
      .then(res => (res.ok ? res.json() : []))
      .then((data: RecentRemediation[]) => {
        setRecentRemediations(Array.isArray(data) ? data.slice(0, 3) : [])
      })
      .catch(() => {
        if (!controller.signal.aborted) {
          setRecentRemediations([])
        }
      })
      .finally(() => {
        if (!controller.signal.aborted) {
          setRecentLoading(false)
        }
      })

    return () => controller.abort()
  }, [space])

  const domainData = useMemo(() => buildDomainData(overview, mcpConnected), [mcpConnected, overview])
  const subtitle = useMemo(() => buildSubtitle(overview), [overview])

  return (
    <Container className={css.main}>
      <div className={css.mcpBanner}>
        <div className={css.mcpStatus}>
          <span className={css.mcpDot} data-connected={mcpConnected === true} />
          <span className={css.mcpLabel}>
            MCP Server {mcpConnected === true ? 'Available' : mcpConnected === false ? 'Unreachable' : 'Checking...'}
          </span>
        </div>
        <button
          className={css.mcpButton}
          onClick={() => {
            if (space && routes.toSOLODEVMcpSetup) {
              history.push(routes.toSOLODEVMcpSetup({ space }))
            }
          }}>
          Connect Client
        </button>
      </div>

      <div className={css.header}>
        <h1 className={css.title}>SoloDev Overview</h1>
        <p className={css.subtitle}>{subtitle}</p>
      </div>

      <div className={css.grid}>
        {domainData.map(domain => (
          <DomainCard
            key={domain.key}
            icon={domain.icon}
            title={domain.title}
            metric={domain.metric}
            status={domain.status}
            statusLabel={domain.statusLabel}
            accentColor={domain.accentColor}
            onClick={
              domain.key === 'remediation' && space && routes.toSOLODEVRemediationQueue
                ? () => history.push(routes.toSOLODEVRemediationQueue({ space }))
                : undefined
            }
          />
        ))}
      </div>

      {overview?.loop_health && (
        <div className={css.loopHealth}>
          <h2 className={css.panelTitle} style={{ marginBottom: 12 }}>Loop Health</h2>
          <div className={css.loopHealthGrid}>
            <div className={css.loopHealthItem}>
              <span className={css.loopHealthValue}>{overview.loop_health.awaiting_apply || 0}</span>
              <span className={css.loopHealthLabel}>Awaiting Apply</span>
            </div>
            <div className={css.loopHealthItem}>
              <span className={css.loopHealthValue}>{overview.loop_health.awaiting_validation || 0}</span>
              <span className={css.loopHealthLabel}>Awaiting Validation</span>
            </div>
            <div className={css.loopHealthItem}>
              <span className={css.loopHealthValue} data-alert={(overview.loop_health.validation_failed || 0) > 0}>
                {overview.loop_health.validation_failed || 0}
              </span>
              <span className={css.loopHealthLabel}>Validation Failed</span>
            </div>
          </div>
        </div>
      )}

      <div className={css.panel}>
        <div className={css.panelHeader}>
          <div>
            <h2 className={css.panelTitle}>Recent Remediations</h2>
            <p className={css.panelSubtitle}>Latest remediation records with live delivery state.</p>
          </div>
          <span
            className={css.prLink}
            style={{ cursor: 'pointer' }}
            onClick={() => {
              if (space && routes.toSOLODEVRemediationQueue) {
                history.push(routes.toSOLODEVRemediationQueue({ space }))
              }
            }}>
            View All
          </span>
        </div>

        {recentLoading ? (
          <div className={css.emptyState}>Loading remediations...</div>
        ) : recentRemediations.length === 0 ? (
          <div className={css.emptyState}>No remediation activity yet.</div>
        ) : (
          <div className={css.remediationList}>
            {recentRemediations.map(remediation => (
              <div
                key={remediation.identifier}
                className={css.remediationItem}
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  if (space && routes.toSOLODEVRemediationDetail) {
                    history.push(routes.toSOLODEVRemediationDetail({ space, remediationId: remediation.identifier }))
                  }
                }}>
                <div className={css.remediationPrimary}>
                  <div className={css.remediationTitle}>{remediation.title}</div>
                  <div className={css.remediationMeta}>
                    <span className={css.remediationBadge}>{remediation.status}</span>
                    <span>
                      Confidence{' '}
                      {typeof remediation.confidence === 'number'
                        ? `${Math.round(remediation.confidence * 100)}%`
                        : 'n/a'}
                    </span>
                  </div>
                </div>
                {remediation.pr_link ? (
                  <a
                    className={css.prLink}
                    href={remediation.pr_link}
                    rel="noreferrer"
                    target="_blank"
                    onClick={e => e.stopPropagation()}>
                    Open PR
                  </a>
                ) : (
                  <span className={css.pendingLabel}>No PR yet</span>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </Container>
  )
}

function buildSubtitle(overview?: TypesSoloDevOverview | null): string {
  const deferred = overview?.deferred_domains?.length || 0
  const tools = overview?.mcp?.tools || 0
  const resources = overview?.mcp?.resources || 0
  const prompts = overview?.mcp?.prompts || 0

  return `${deferred} deferred domains · ${tools} MCP tools · ${resources} resources · ${prompts} prompts`
}

function buildDomainData(
  overview: TypesSoloDevOverview | null | undefined,
  mcpConnected: boolean | null
): DomainData[] {
  const security = overview?.security
  const errors = overview?.errors
  const remediation = overview?.remediation
  const mcp = overview?.mcp
  const deferred = new Set(overview?.deferred_domains || [])

  return [
    plannedCard(
      'pipelines',
      '\u25B6',
      'Pipelines',
      deferred.has('pipelines') ? 'Deferred this sprint' : 'Partially wired',
      '#58a6ff'
    ),
    {
      icon: '\u26A1',
      title: 'Security',
      metric: `${security?.open_findings || 0} open findings`,
      status: securityStatus(security?.open_findings, security?.critical, security?.availability),
      statusLabel: securityLabel(security?.open_findings, security?.critical, security?.availability),
      accentColor: '#f85149',
      key: 'security'
    },
    plannedCard(
      'quality',
      '\u2713',
      'Quality Gates',
      deferred.has('quality') ? 'Catalog only' : 'Partially wired',
      '#3fb950'
    ),
    {
      icon: '\u2715',
      title: 'Error Tracker',
      metric: `${errors?.open || 0} open groups`,
      status: errorStatus(errors?.open, errors?.fatal),
      statusLabel: errorLabel(errors?.open, errors?.fatal, errors?.availability),
      accentColor: '#d29922',
      key: 'errors'
    },
    {
      icon: '\u2692',
      title: 'Remediation',
      metric: `${remediation?.pending || 0} pending · ${remediation?.applied || 0} applied`,
      status: remediationStatus(remediation),
      statusLabel: remediationLabel(remediation),
      accentColor: '#39d2c0',
      key: 'remediation'
    },
    {
      icon: '\u2699',
      title: 'MCP Server',
      metric: `${mcp?.tools || 0} tools · ${mcp?.resources || 0} resources`,
      status: mcpConnected === false ? 'warning' : 'healthy',
      statusLabel: mcpConnected === false ? 'Unavailable' : 'Live',
      accentColor: '#ff7b72',
      key: 'mcp'
    },
    plannedCard(
      'health',
      '\u2665',
      'Health Monitor',
      deferred.has('health') ? 'Deferred this sprint' : 'Partially wired',
      '#2ea043'
    ),
    plannedCard(
      'tech_debt',
      '\u2630',
      'Tech Debt',
      deferred.has('tech_debt') ? 'Deferred this sprint' : 'Partially wired',
      '#e3b341'
    )
  ]
}

function plannedCard(key: string, icon: string, title: string, metric: string, accentColor: string): DomainData {
  return {
    icon,
    title,
    metric,
    status: 'warning',
    statusLabel: 'Planned',
    accentColor,
    key
  }
}

function securityStatus(openFindings?: number, critical?: number, availability?: string): DomainStatus {
  if ((critical || 0) > 0) return 'critical'
  if ((openFindings || 0) > 0) return 'warning'
  if (availability === 'ready') return 'healthy'
  return 'warning'
}

function securityLabel(openFindings?: number, critical?: number, availability?: string): string {
  if ((critical || 0) > 0) return 'Critical'
  if ((openFindings || 0) > 0) return 'Warning'
  if (availability === 'ready') return 'Ready'
  return 'Blocked'
}

function errorStatus(open?: number, fatal?: number): DomainStatus {
  if ((fatal || 0) > 0) return 'critical'
  if ((open || 0) > 0) return 'warning'
  return 'healthy'
}

function errorLabel(open?: number, fatal?: number, availability?: string): string {
  if ((fatal || 0) > 0) return 'Critical'
  if ((open || 0) > 0) return 'Warning'
  if (availability === 'ready') return 'Ready'
  return 'Blocked'
}

function remediationStatus(remediation?: TypesSoloDevOverview['remediation']): DomainStatus {
  if (!remediation) return 'warning'
  if ((remediation.processing || 0) > 0 || (remediation.pending || 0) > 0) return 'active'
  if ((remediation.failed || 0) > 0) return 'warning'
  if ((remediation.applied || 0) > 0 || (remediation.completed || 0) > 0) return 'healthy'
  return remediation.availability === 'ready' ? 'healthy' : 'warning'
}

function remediationLabel(remediation?: TypesSoloDevOverview['remediation']): string {
  if (!remediation) return 'Blocked'
  if ((remediation.processing || 0) > 0 || (remediation.pending || 0) > 0) return 'Active'
  if ((remediation.failed || 0) > 0) return 'Warning'
  if ((remediation.applied || 0) > 0 || (remediation.completed || 0) > 0) return 'Ready'
  return remediation.availability === 'ready' ? 'Idle' : 'Blocked'
}
