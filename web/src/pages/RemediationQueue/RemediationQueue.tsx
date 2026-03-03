/*
 * Copyright 2026 EolaFam1828. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useCallback, useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'
import { Container } from '@harnessio/uicore'
import { useAppContext } from 'AppContext'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import css from './RemediationQueue.module.scss'

interface RemediationDelivery {
  mode: string
  state: string
  pr_number?: number
  last_error?: string
  attempted_at?: number
}

interface RemediationValidation {
  state: string
  pipeline_identifier?: string
  execution_number?: number
  execution_status?: string
  last_error?: string
  started_at?: number
  completed_at?: number
}

interface RemediationItem {
  identifier: string
  title: string
  status: string
  trigger_source: string
  trigger_ref?: string
  confidence?: number
  pr_link?: string
  fix_branch?: string
  delivery?: RemediationDelivery
  validation?: RemediationValidation
  created?: number
  updated?: number
}

const STATUSES = ['all', 'pending', 'processing', 'completed', 'applied', 'failed', 'dismissed'] as const
const TRIGGERS = ['all', 'pipeline', 'error_tracker', 'security_scan', 'manual'] as const

const TRIGGER_LABELS: Record<string, string> = {
  pipeline: 'Pipeline',
  error_tracker: 'Error',
  security_scan: 'Security',
  quality_gate: 'Quality',
  manual: 'Manual'
}

function statusClass(status: string): string {
  switch (status) {
    case 'pending': return css.badge_pending
    case 'processing': return css.badge_processing
    case 'completed': return css.badge_completed
    case 'applied': return css.badge_applied
    case 'failed': return css.badge_failed
    case 'dismissed': return css.badge_dismissed
    default: return ''
  }
}

function deliveryClass(state: string): string {
  switch (state) {
    case 'applied': return css.deliveryChip_applied
    case 'branch_ready': return css.deliveryChip_branch_ready
    case 'failed': return css.deliveryChip_failed
    default: return ''
  }
}

function validationClass(state: string): string {
  switch (state) {
    case 'passed': return css.validationChip_passed
    case 'failed': return css.validationChip_failed
    case 'running': return css.validationChip_running
    case 'queued': return css.validationChip_queued
    case 'unavailable': return css.validationChip_unavailable
    default: return ''
  }
}

function formatTime(ms?: number): string {
  if (!ms) return ''
  const d = new Date(ms)
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

export default function RemediationQueue() {
  const { routes } = useAppContext()
  const history = useHistory()
  const { space } = useGetRepositoryMetadata()
  const [items, setItems] = useState<RemediationItem[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [triggerFilter, setTriggerFilter] = useState<string>('all')

  const fetchItems = useCallback(() => {
    if (!space) return
    setLoading(true)
    const params = new URLSearchParams()
    if (statusFilter !== 'all') params.set('status', statusFilter)
    if (triggerFilter !== 'all') params.set('trigger_source', triggerFilter)
    params.set('limit', '50')

    fetch(`/api/v1/spaces/${space}/remediations?${params}`)
      .then(res => (res.ok ? res.json() : []))
      .then((data: RemediationItem[]) => setItems(Array.isArray(data) ? data : []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [space, statusFilter, triggerFilter])

  useEffect(() => { fetchItems() }, [fetchItems])

  return (
    <Container className={css.main}>
      <div className={css.header}>
        <div className={css.titleRow}>
          <h1 className={css.title}>Remediation Queue</h1>
          <span
            className={css.backLink}
            onClick={() => routes.toSOLODEVDashboard && history.push(routes.toSOLODEVDashboard({ space: space || '' }))}>
            Back to Dashboard
          </span>
        </div>
        <p className={css.subtitle}>
          {items.length} remediation{items.length !== 1 ? 's' : ''} matching filters
        </p>
      </div>

      <div className={css.filters}>
        {STATUSES.map(s => (
          <button
            key={s}
            className={css.filterBtn}
            data-active={statusFilter === s}
            onClick={() => setStatusFilter(s)}>
            {s === 'all' ? 'All Status' : s}
          </button>
        ))}
      </div>

      <div className={css.filters}>
        {TRIGGERS.map(t => (
          <button
            key={t}
            className={css.filterBtn}
            data-active={triggerFilter === t}
            onClick={() => setTriggerFilter(t)}>
            {t === 'all' ? 'All Sources' : TRIGGER_LABELS[t] || t}
          </button>
        ))}
      </div>

      {loading ? (
        <div className={css.loadingState}>Loading remediations...</div>
      ) : items.length === 0 ? (
        <div className={css.emptyState}>No remediations match the current filters.</div>
      ) : (
        <div className={css.list}>
          {items.map(item => (
            <div
              key={item.identifier}
              className={css.row}
              onClick={() => {
                if (routes.toSOLODEVRemediationDetail) {
                  history.push(routes.toSOLODEVRemediationDetail({ space: space || '', remediationId: item.identifier }))
                }
              }}>
              <div className={css.rowPrimary}>
                <div className={css.rowTitle}>{item.title}</div>
                <div className={css.rowMeta}>
                  <span className={`${css.badge} ${statusClass(item.status)}`}>{item.status}</span>
                  {item.delivery && item.delivery.state !== 'not_attempted' && (
                    <span className={`${css.deliveryChip} ${deliveryClass(item.delivery.state)}`}>
                      {item.delivery.state.replace('_', ' ')}
                    </span>
                  )}
                  {item.validation && item.validation.state !== 'not_attempted' && (
                    <span className={`${css.validationChip} ${validationClass(item.validation.state)}`}>
                      {item.validation.state}
                    </span>
                  )}
                  <span className={css.triggerBadge}>{TRIGGER_LABELS[item.trigger_source] || item.trigger_source}</span>
                  {item.created ? <span>{formatTime(item.created)}</span> : null}
                </div>
              </div>
              <div className={css.rowActions}>
                {typeof item.confidence === 'number' && item.confidence > 0 && (
                  <span className={css.confidence}>{Math.round(item.confidence * 100)}%</span>
                )}
                {item.pr_link ? (
                  <a
                    className={css.prLink}
                    href={item.pr_link}
                    rel="noreferrer"
                    target="_blank"
                    onClick={e => e.stopPropagation()}>
                    PR
                  </a>
                ) : null}
              </div>
            </div>
          ))}
        </div>
      )}
    </Container>
  )
}
