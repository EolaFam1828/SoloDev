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
import { Render } from 'react-jsx-match'
import { Container, Layout } from '@harnessio/uicore'
import { Flag, Bug, Shield, Activity, AlertTriangle, CheckCircle } from 'iconoir-react'
import { moduleRoutes } from 'ModuleRoutes'
import { NavMenuItem } from './NavMenuItem'
import css from './ModuleMenu.module.scss'

interface ModuleMenuProps {
  space?: string
}

export const ModuleMenu: React.FC<ModuleMenuProps> = ({ space }) => {
  return (
    <Render when={!!space}>
      <Container className={css.moduleSection}>
        <Layout.Vertical spacing="small">
          <NavMenuItem
            label="Feature Flags"
            to={moduleRoutes.toFeatureFlags({ space: space as string })}
            customIcon={<Flag />}
          />

          <NavMenuItem
            label="Technical Debt"
            to={moduleRoutes.toTechDebt({ space: space as string })}
            customIcon={<Bug />}
          />

          <NavMenuItem
            label="Security Scanner"
            to={moduleRoutes.toSecurityScans({ space: space as string })}
            customIcon={<Shield />}
          />

          <NavMenuItem
            label="Uptime Monitor"
            to={moduleRoutes.toMonitors({ space: space as string })}
            customIcon={<Activity />}
          />

          <NavMenuItem
            label="Error Tracker"
            to={moduleRoutes.toErrors({ space: space as string })}
            customIcon={<AlertTriangle />}
          />

          <NavMenuItem
            label="Quality Gates"
            to={moduleRoutes.toQualityGates({ space: space as string })}
            customIcon={<CheckCircle />}
          />
        </Layout.Vertical>
      </Container>
    </Render>
  )
}
