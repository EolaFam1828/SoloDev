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

import React from 'react'
import css from './DomainCard.module.scss'

export type DomainStatus = 'healthy' | 'warning' | 'critical' | 'active' | 'passing'

interface DomainCardProps {
  icon: string
  title: string
  summary: string
  metric: string
  status: DomainStatus
  statusLabel: string
  accentColor: string
  capability: string
  actionLabel: string
  onClick?: () => void
}

const STATUS_INDICATORS: Record<DomainStatus, { symbol: string; className: keyof typeof css }> = {
  healthy: { symbol: '\u25CF', className: 'healthy' },
  passing: { symbol: '\u25CF', className: 'passing' },
  active: { symbol: '\u25CF', className: 'active' },
  warning: { symbol: '\u25B2', className: 'warning' },
  critical: { symbol: '\u25CF', className: 'critical' }
}

export const DomainCard: React.FC<DomainCardProps> = ({
  icon,
  title,
  summary,
  metric,
  status,
  statusLabel,
  accentColor,
  capability,
  actionLabel,
  onClick
}) => {
  const indicator = STATUS_INDICATORS[status]

  return (
    <div
      className={css.card}
      style={{ '--card-accent': accentColor } as React.CSSProperties}
      onClick={onClick}
      role={onClick ? 'button' : 'article'}
      tabIndex={onClick ? 0 : -1}
      aria-disabled={!onClick}
      onKeyDown={e => e.key === 'Enter' && onClick?.()}>
      <div className={css.cardTopline}>
        <span className={css.cardIcon}>{icon}</span>
        <span className={css.cardEyebrow}>{capability}</span>
        <span className={css.cardMetric}>{metric}</span>
      </div>
      <div className={css.cardBody}>
        <div className={css.cardTitle}>{title}</div>
        <p className={css.cardSummary}>{summary}</p>
      </div>
      <div className={css.cardFooter}>
        <div className={`${css.cardStatus} ${css[indicator.className]}`}>
          <span className={css.statusDot}>{indicator.symbol}</span>
          <span>{statusLabel}</span>
        </div>
        <span className={css.cardAction}>{actionLabel}</span>
      </div>
      <div className={css.cardAccent} />
      <div className={css.cardGlow} />
    </div>
  )
}
