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
  metric: string
  status: DomainStatus
  statusLabel: string
  accentColor: string
  onClick?: () => void
}

const STATUS_INDICATORS: Record<DomainStatus, { symbol: string; className: string }> = {
  healthy: { symbol: '\u25CF', className: 'healthy' },
  passing: { symbol: '\u25CF', className: 'passing' },
  active: { symbol: '\u25CF', className: 'active' },
  warning: { symbol: '\u25B2', className: 'warning' },
  critical: { symbol: '\u25CF', className: 'critical' }
}

export const DomainCard: React.FC<DomainCardProps> = ({
  icon,
  title,
  metric,
  status,
  statusLabel,
  accentColor,
  onClick
}) => {
  const indicator = STATUS_INDICATORS[status]

  return (
    <div
      className={css.card}
      style={{ '--card-accent': accentColor } as React.CSSProperties}
      onClick={onClick}
      role="button"
      tabIndex={0}
      onKeyDown={e => e.key === 'Enter' && onClick?.()}>
      <div className={css.cardHeader}>
        <span className={css.cardIcon}>{icon}</span>
        <span className={css.cardTitle}>{title}</span>
      </div>
      <div className={css.cardMetric}>{metric}</div>
      <div className={`${css.cardStatus} ${css[indicator.className]}`}>
        <span className={css.statusDot}>{indicator.symbol}</span>
        <span>{statusLabel}</span>
      </div>
    </div>
  )
}
