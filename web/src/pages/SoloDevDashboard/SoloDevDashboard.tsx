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
import { useHistory } from 'react-router-dom'
import { Container } from '@harnessio/uicore'
import { useAppContext } from 'AppContext'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
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

const DOMAIN_DATA: DomainData[] = [
  {
    icon: '\u25B6',
    title: 'Pipelines',
    metric: '3 active',
    status: 'healthy',
    statusLabel: 'Healthy',
    accentColor: '#58a6ff',
    key: 'pipelines'
  },
  {
    icon: '\u26A1',
    title: 'Security',
    metric: '12 findings',
    status: 'warning',
    statusLabel: 'Warning',
    accentColor: '#f85149',
    key: 'security'
  },
  {
    icon: '\u2713',
    title: 'Quality Gates',
    metric: '5 rules',
    status: 'passing',
    statusLabel: 'Passing',
    accentColor: '#3fb950',
    key: 'quality'
  },
  {
    icon: '\u2715',
    title: 'Error Tracker',
    metric: '2 unresolved',
    status: 'warning',
    statusLabel: 'Warning',
    accentColor: '#d29922',
    key: 'errors'
  },
  {
    icon: '\u2692',
    title: 'Remediation',
    metric: '1 pending',
    status: 'active',
    statusLabel: 'Active',
    accentColor: '#39d2c0',
    key: 'remediation'
  },
  {
    icon: '\u2665',
    title: 'Health Monitor',
    metric: 'All passing',
    status: 'healthy',
    statusLabel: 'Healthy',
    accentColor: '#2ea043',
    key: 'health'
  },
  {
    icon: '\u2691',
    title: 'Feature Flags',
    metric: '8 flags',
    status: 'active',
    statusLabel: 'Active',
    accentColor: '#bc8cff',
    key: 'flags'
  },
  {
    icon: '\u2630',
    title: 'Tech Debt',
    metric: '15 items',
    status: 'warning',
    statusLabel: 'Medium',
    accentColor: '#e3b341',
    key: 'debt'
  }
]

export default function SoloDevDashboard() {
  const { routes } = useAppContext()
  const history = useHistory()
  const { space } = useGetRepositoryMetadata()
  const [mcpConnected, setMcpConnected] = useState<boolean | null>(null)

  useEffect(() => {
    fetch('/api/v1/system/config')
      .then(res => {
        setMcpConnected(res.ok)
      })
      .catch(() => {
        setMcpConnected(false)
      })
  }, [])

  return (
    <Container className={css.main}>
      {/* MCP Banner */}
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

      {/* Header */}
      <div className={css.header}>
        <h1 className={css.title}>DevOps Dashboard</h1>
        <p className={css.subtitle}>8 domains &middot; 24 MCP tools &middot; 8 resources</p>
      </div>

      {/* Domain Grid */}
      <div className={css.grid}>
        {DOMAIN_DATA.map(domain => (
          <DomainCard
            key={domain.key}
            icon={domain.icon}
            title={domain.title}
            metric={domain.metric}
            status={domain.status}
            statusLabel={domain.statusLabel}
            accentColor={domain.accentColor}
            onClick={() => {
              // Pipelines links to existing pipelines page if we have a repo
              // Others are stubs that stay on dashboard for now
            }}
          />
        ))}
      </div>
    </Container>
  )
}
