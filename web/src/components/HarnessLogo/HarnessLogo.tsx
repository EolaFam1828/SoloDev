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
import { Container, Layout } from '@harnessio/uicore'
import { Link } from 'react-router-dom'
import { useAppContext } from 'AppContext'
import css from './HarnessLogo.module.scss'

export const HarnessLogo: React.FC = () => {
  const { routes } = useAppContext()

  return (
    <Container className={css.main}>
      <Link to={routes.toCODEHome()}>
        <Layout.Horizontal spacing="small" className={css.layout}>
          <div className={css.logoIcon} aria-hidden="true">
            <div className={css.logoBadge}>
              <svg width="34" height="34" viewBox="0 0 34 34" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="1" y="1" width="32" height="32" rx="10" className={css.badgeFrame} />
                <path d="M10.5 11.5L16.5 17L10.5 22.5" className={css.badgeGlyph} />
                <path d="M18.5 22.5H24" className={css.badgeGlyph} />
                <circle cx="25.5" cy="8.5" r="2.5" className={css.badgeSpark} />
              </svg>
            </div>
          </div>
          <span className={css.copy}>
            <span className={css.text}>SoloDev</span>
            <span className={css.meta}>AI DevOps OS</span>
          </span>
        </Layout.Horizontal>
      </Link>
    </Container>
  )
}
