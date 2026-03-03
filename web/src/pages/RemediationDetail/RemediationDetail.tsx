/*
 * Copyright 2026 EolaFam1828. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useCallback, useEffect, useState } from 'react'
import { useHistory, useParams } from 'react-router-dom'
import { Container } from '@harnessio/uicore'
import { useAppContext } from 'AppContext'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import css from './RemediationDetail.module.scss'

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
  url?: string
  last_error?: string
  started_at?: number
  completed_at?: number
}

interface RemediationFull {
  id: number
  identifier: string
  repo_id: number
  title: string
  description?: string
  status: string
  trigger_source: string
  trigger_ref?: string
  branch: string
  commit_sha?: string
  error_log: string
  source_code?: string
  file_path?: string
  ai_model?: string
  ai_response?: string
  patch_diff?: string
  fix_branch?: string
  pr_link?: string
  confidence?: number
  tokens_used?: number
  duration_ms?: number
  delivery?: RemediationDelivery
  validation?: RemediationValidation
  created?: number
  updated?: number
}

const TRIGGER_LABELS: Record<string, string> = {
  pipeline: 'Pipeline Failure',
  error_tracker: 'Runtime Error',
  security_scan: 'Security Finding',
  quality_gate: 'Quality Gate',
  manual: 'Manual'
}

function statusClass(status: string): string {
  const map: Record<string, string> = {
    pending: css.badge_pending,
    processing: css.badge_processing,
    completed: css.badge_completed,
    applied: css.badge_applied,
    failed: css.badge_failed,
    dismissed: css.badge_dismissed
  }
  return map[status] || ''
}

function deliveryClass(state: string): string {
  const map: Record<string, string> = {
    not_attempted: css.delivery_not_attempted,
    branch_ready: css.delivery_branch_ready,
    applied: css.delivery_applied,
    failed: css.delivery_failed
  }
  return map[state] || ''
}

function validationClass(state: string): string {
  const map: Record<string, string> = {
    not_attempted: css.validation_not_attempted,
    queued: css.validation_queued,
    running: css.validation_running,
    passed: css.validation_passed,
    failed: css.validation_failed,
    unavailable: css.validation_unavailable
  }
  return map[state] || ''
}

function formatTime(ms?: number): string {
  if (!ms) return 'n/a'
  return new Date(ms).toLocaleString()
}

export default function RemediationDetail() {
  const { routes } = useAppContext()
  const history = useHistory()
  const { space } = useGetRepositoryMetadata()
  const { remediationId } = useParams<{ remediationId: string }>()
  const [rem, setRem] = useState<RemediationFull | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [actionLoading, setActionLoading] = useState(false)

  const fetchDetail = useCallback(() => {
    if (!space || !remediationId) return
    setLoading(true)
    setError(null)
    fetch(`/api/v1/spaces/${space}/remediations/${remediationId}`)
      .then(res => {
        if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
        return res.json()
      })
      .then((data: RemediationFull) => setRem(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [space, remediationId])

  useEffect(() => { fetchDetail() }, [fetchDetail])

  const doApply = useCallback(async () => {
    if (!space || !remediationId) return
    setActionLoading(true)
    try {
      const res = await fetch(`/api/v1/spaces/${space}/remediations/${remediationId}/apply`, { method: 'POST' })
      if (!res.ok) {
        const body = await res.json().catch(() => ({ message: res.statusText }))
        throw new Error(body.message || res.statusText)
      }
      fetchDetail()
    } catch (err: any) {
      setError(`Apply failed: ${err.message}`)
    } finally {
      setActionLoading(false)
    }
  }, [space, remediationId, fetchDetail])

  const doValidate = useCallback(async () => {
    if (!space || !remediationId) return
    setActionLoading(true)
    try {
      const res = await fetch(`/api/v1/spaces/${space}/remediations/${remediationId}/validate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({})
      })
      if (!res.ok) {
        const body = await res.json().catch(() => ({ message: res.statusText }))
        throw new Error(body.message || res.statusText)
      }
      fetchDetail()
    } catch (err: any) {
      setError(`Validate failed: ${err.message}`)
    } finally {
      setActionLoading(false)
    }
  }, [space, remediationId, fetchDetail])

  const doStatusUpdate = useCallback(async (newStatus: string) => {
    if (!space || !remediationId) return
    setActionLoading(true)
    try {
      const res = await fetch(`/api/v1/spaces/${space}/remediations/${remediationId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: newStatus })
      })
      if (!res.ok) {
        const body = await res.json().catch(() => ({ message: res.statusText }))
        throw new Error(body.message || res.statusText)
      }
      fetchDetail()
    } catch (err: any) {
      setError(`Status update failed: ${err.message}`)
    } finally {
      setActionLoading(false)
    }
  }, [space, remediationId, fetchDetail])

  if (loading) return <Container className={css.main}><div className={css.loadingState}>Loading remediation...</div></Container>
  if (error && !rem) return <Container className={css.main}><div className={css.errorState}>{error}</div></Container>
  if (!rem) return <Container className={css.main}><div className={css.errorState}>Remediation not found</div></Container>

  const canApply = rem.status === 'completed' || (rem.status === 'applied' && rem.delivery?.state === 'failed')
  const canValidate = rem.status === 'applied' && rem.fix_branch &&
    (!rem.validation || rem.validation.state === 'not_attempted' || rem.validation.state === 'failed' || rem.validation.state === 'unavailable')
  const canDismiss = rem.status !== 'dismissed'
  const canReopen = rem.status === 'dismissed'
  const deliveryState = rem.delivery?.state || 'not_attempted'
  const validationState = rem.validation?.state || 'not_attempted'

  return (
    <Container className={css.main}>
      <div className={css.header}>
        <span
          className={css.backLink}
          onClick={() => routes.toSOLODEVRemediationQueue && history.push(routes.toSOLODEVRemediationQueue({ space: space || '' }))}>
          &larr; Back to Queue
        </span>

        <div className={css.titleRow}>
          <div>
            <h1 className={css.title}>{rem.title}</h1>
            <span className={css.identifier}>{rem.identifier}</span>
          </div>
          <div className={css.actions}>
            {canApply && (
              <button className={`${css.actionBtn} ${css.actionBtn_primary}`} disabled={actionLoading} onClick={doApply}>
                {rem.delivery?.state === 'failed' ? 'Retry Delivery' : 'Apply'}
              </button>
            )}
            {canValidate && (
              <button className={`${css.actionBtn} ${css.actionBtn_validate}`} disabled={actionLoading} onClick={doValidate}>
                Validate
              </button>
            )}
            {rem.pr_link && (
              <a className={css.actionBtn} href={rem.pr_link} target="_blank" rel="noreferrer">
                Open PR
              </a>
            )}
            {canDismiss && (
              <button className={`${css.actionBtn} ${css.actionBtn_danger}`} disabled={actionLoading} onClick={() => doStatusUpdate('dismissed')}>
                Dismiss
              </button>
            )}
            {canReopen && (
              <button className={css.actionBtn} disabled={actionLoading} onClick={() => doStatusUpdate('completed')}>
                Reopen
              </button>
            )}
          </div>
        </div>
      </div>

      {error && <div className={css.errorText} style={{ marginBottom: 16, position: 'relative', zIndex: 1 }}>{error}</div>}

      <div className={css.grid}>
        <div className={css.panel}>
          <div className={css.panelLabel}>Status</div>
          <span className={`${css.badge} ${statusClass(rem.status)}`}>{rem.status}</span>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Delivery</div>
          <span className={`${css.deliveryState} ${deliveryClass(deliveryState)}`}>
            {deliveryState.replace(/_/g, ' ')}
          </span>
          {rem.delivery?.last_error && (
            <div className={css.errorText} style={{ marginTop: 8, fontSize: 12 }}>
              {rem.delivery.last_error}
            </div>
          )}
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Validation</div>
          <span className={`${css.validationState} ${validationClass(validationState)}`}>
            {validationState.replace(/_/g, ' ')}
          </span>
          {rem.validation?.pipeline_identifier && (
            <div className={css.validationMeta}>Pipeline: {rem.validation.pipeline_identifier}</div>
          )}
          {rem.validation?.execution_status && (
            <div className={css.validationMeta}>CI: {rem.validation.execution_status}</div>
          )}
          {rem.validation?.last_error && (
            <div className={css.errorText} style={{ marginTop: 8, fontSize: 12 }}>
              {rem.validation.last_error}
            </div>
          )}
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Trigger</div>
          <div className={css.panelValue}>
            {TRIGGER_LABELS[rem.trigger_source] || rem.trigger_source}
            {rem.trigger_ref && <span style={{ color: '#6e7681', marginLeft: 8 }}>#{rem.trigger_ref}</span>}
          </div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Confidence</div>
          <div className={css.panelValue}>
            {typeof rem.confidence === 'number' && rem.confidence > 0 ? `${Math.round(rem.confidence * 100)}%` : 'n/a'}
          </div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Branch</div>
          <div className={css.panelValue}>{rem.branch || 'n/a'}</div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Fix Branch</div>
          <div className={css.panelValue}>{rem.fix_branch || 'n/a'}</div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>File Path</div>
          <div className={css.panelValue}>{rem.file_path || 'n/a'}</div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>AI Model</div>
          <div className={css.panelValue}>{rem.ai_model || 'n/a'}</div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Created</div>
          <div className={css.panelValue}>{formatTime(rem.created)}</div>
        </div>

        <div className={css.panel}>
          <div className={css.panelLabel}>Updated</div>
          <div className={css.panelValue}>{formatTime(rem.updated)}</div>
        </div>

        {rem.pr_link && (
          <div className={`${css.panel} ${css.panelFull}`}>
            <div className={css.panelLabel}>Pull Request</div>
            <a className={css.link} href={rem.pr_link} target="_blank" rel="noreferrer">{rem.pr_link}</a>
          </div>
        )}
      </div>

      {rem.patch_diff && (
        <div className={css.panel} style={{ marginBottom: 16, position: 'relative', zIndex: 1 }}>
          <div className={css.panelLabel}>Patch Diff</div>
          <div className={css.diffBlock}>{rem.patch_diff}</div>
        </div>
      )}

      {rem.ai_response && (
        <div className={css.panel} style={{ marginBottom: 16, position: 'relative', zIndex: 1 }}>
          <div className={css.panelLabel}>AI Response</div>
          <div className={css.aiResponse}>{rem.ai_response}</div>
        </div>
      )}
    </Container>
  )
}
