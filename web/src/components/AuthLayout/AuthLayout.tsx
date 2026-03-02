/*
 * Copyright 2023 Harness, Inc.
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
import { Container } from '@harnessio/uicore'
import css from './AuthLayout.module.scss'

const AuthLayout: React.FC<React.PropsWithChildren<unknown>> = props => {
  return (
    <div className={css.layout}>
      <div className={css.brandColumn}>
        <div className={css.brandContent}>
          <div className={css.logoRow}>
            <svg width="32" height="32" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="28" height="28" rx="6" fill="#7C3AED"/>
              <path d="M8 8l6 6-6 6" stroke="#fff" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"/>
              <path d="M16 20h4" stroke="#fff" strokeWidth="2.5" strokeLinecap="round"/>
            </svg>
            <span className={css.logoText}>SoloDev</span>
          </div>
          <h1 className={css.heroTitle}>Your AI-Native<br/>DevOps Platform</h1>
          <p className={css.heroSub}>8 domains &middot; 24 MCP tools &middot; One interface</p>
          <div className={css.domainGrid}>
            {['Pipelines', 'Security', 'Quality Gates', 'Error Tracker', 'Remediation', 'Health', 'Feature Flags', 'Tech Debt'].map(d => (
              <span key={d} className={css.domainTag}>{d}</span>
            ))}
          </div>
        </div>
      </div>
      <div className={css.cardColumn}>
        <div className={css.card}>
          <Container className={css.cardChildren}>{props.children}</Container>
        </div>
      </div>
    </div>
  )
}

export default AuthLayout
