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

import React, { useEffect, useMemo, useState } from 'react'
import { Render } from 'react-jsx-match'
import { useHistory, useRouteMatch } from 'react-router-dom'
import { BookmarkBook, Settings } from 'iconoir-react'

import { Container, Layout } from '@harnessio/uicore'

import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import { useGetSpaceParam } from 'hooks/useGetSpaceParam'
import { useStrings } from 'framework/strings'
import type { SpaceSpaceOutput } from 'services/code'
import { SpaceSelector } from 'components/SpaceSelector/SpaceSelector'
import { useAppContext } from 'AppContext'
import { isGitRev } from 'utils/GitUtils'
import { NavMenuItem } from './NavMenuItem'
import css from './DefaultMenu.module.scss'

interface DomainLink {
  icon: string
  label: string
  key: string
  accentColor: string
}

const DOMAIN_LINKS: DomainLink[] = [
  { icon: '\u25B6', label: 'Pipelines', key: 'pipelines', accentColor: 'var(--solodev-pipeline-blue)' },
  { icon: '\u26A1', label: 'Security', key: 'security', accentColor: 'var(--solodev-security-red)' },
  { icon: '\u2713', label: 'Quality Gates', key: 'quality', accentColor: 'var(--solodev-quality-green)' },
  { icon: '\u2715', label: 'Error Tracker', key: 'errors', accentColor: 'var(--solodev-error-orange)' },
  { icon: '\u2692', label: 'Remediation', key: 'remediation', accentColor: 'var(--solodev-remediation-cyan)' },
  { icon: '\u2665', label: 'Health Monitor', key: 'health', accentColor: 'var(--solodev-health-emerald)' },
  { icon: '\u2691', label: 'Feature Flags', key: 'flags', accentColor: 'var(--solodev-flags-purple)' },
  { icon: '\u2630', label: 'Tech Debt', key: 'debt', accentColor: 'var(--solodev-debt-amber)' }
]

export const DefaultMenu: React.FC = () => {
  const history = useHistory()
  const { routes, standalone, isCurrentSessionPublic } = useAppContext()
  const [selectedSpace, setSelectedSpace] = useState<SpaceSpaceOutput | undefined>()
  const spaceFromRoute = useGetSpaceParam()
  const { repoMetadata, gitRef, commitRef } = useGetRepositoryMetadata()
  const { getString } = useStrings()
  const repoPath = useMemo(() => repoMetadata?.path || '', [repoMetadata])
  const routeMatch = useRouteMatch()
  const isCommitSelected = useMemo(() => routeMatch.path === '/:space*/:repoName/commit/:commitRef*', [routeMatch])
  const activeSpacePath = selectedSpace?.path || spaceFromRoute || repoMetadata?.path?.split('/')[0] || ''
  const activeSpaceLabel = selectedSpace?.identifier || activeSpacePath

  const isFilesSelected = useMemo(
    () =>
      !isCommitSelected &&
      (routeMatch.path === '/:space*/:repoName' || routeMatch.path.startsWith('/:space*/:repoName/edit')),
    [routeMatch, isCommitSelected]
  )
  const isWebhookSelected = useMemo(() => routeMatch.path.startsWith('/:space*/:repoName/webhook'), [routeMatch])
  const _gitRef = useMemo(() => {
    const ref = commitRef || gitRef
    return !isGitRev(ref) ? ref : ''
  }, [commitRef, gitRef])

  const isDashboardSelected = useMemo(() => routeMatch.path.includes('/dashboard'), [routeMatch])
  const isMcpSetupSelected = useMemo(() => routeMatch.path.includes('/mcp-setup'), [routeMatch])

  useEffect(() => {
    if (!selectedSpace && activeSpacePath) {
      setSelectedSpace({
        id: -1,
        identifier: activeSpacePath,
        path: activeSpacePath
      })
    }
  }, [activeSpacePath, selectedSpace])

  return (
    <Container className={css.main}>
      <Layout.Vertical spacing="small">
        {/* Space Selector */}
        <Render when={!isCurrentSessionPublic}>
          <SpaceSelector
            onSelect={(_selectedSpace, isUserAction) => {
              setSelectedSpace(_selectedSpace)
              if (_selectedSpace.path === '' && _selectedSpace.id === -1) {
                setSelectedSpace(undefined)
              }
              if (isUserAction) {
                history.push(
                  routes.toSOLODEVDashboard
                    ? routes.toSOLODEVDashboard({ space: _selectedSpace.path as string })
                    : routes.toCODERepositories({ space: _selectedSpace.path as string })
                )
              }
            }}
          />
        </Render>

        {/* Dashboard Link */}
        <Render when={!!activeSpacePath}>
          <div className={css.projectPill}>
            <span className={css.projectLabel}>Project</span>
            <strong>{activeSpaceLabel}</strong>
          </div>
          <div className={css.sectionLabel}>OVERVIEW</div>
          <NavMenuItem
            label="Dashboard"
            to={
              routes.toSOLODEVDashboard
                ? routes.toSOLODEVDashboard({ space: activeSpacePath })
                : routes.toCODERepositories({ space: activeSpacePath })
            }
            isSelected={isDashboardSelected}
          />
        </Render>

        {/* Domain Navigation */}
        <Render when={!!activeSpacePath}>
          <div className={css.sectionDivider} />
          <div className={css.sectionLabel}>DEVOPS DOMAINS</div>
          {DOMAIN_LINKS.map(domain => (
            <div key={domain.key} className={css.domainItem}>
              <span className={css.domainIcon} style={{ color: domain.accentColor }}>
                {domain.icon}
              </span>
              <NavMenuItem
                label={domain.label}
                to={
                  domain.key === 'pipelines' && repoMetadata
                    ? routes.toCODEPipelines({ repoPath })
                    : routes.toSOLODEVDashboard
                      ? routes.toSOLODEVDashboard({ space: activeSpacePath })
                      : routes.toCODERepositories({ space: activeSpacePath })
                }
              />
            </div>
          ))}
        </Render>

        {/* Repositories */}
        <Render when={!!activeSpacePath}>
          <div className={css.sectionDivider} />
          <div className={css.sectionLabel}>CODE</div>
          <NavMenuItem
            label={getString('repositories')}
            to={routes.toCODERepositories({ space: activeSpacePath })}
            isDeselected={!!repoMetadata}
            isHighlighted={!!repoMetadata}
            customIcon={<BookmarkBook />}
          />
        </Render>

        {/* Repo sub-menu */}
        <Render when={repoMetadata}>
          <Container className={css.repoLinks}>
            <Layout.Vertical spacing="small">
              <NavMenuItem
                data-code-repo-section="files"
                isSubLink
                isSelected={isFilesSelected}
                label={getString('files')}
                to={routes.toCODERepository({ repoPath, gitRef: _gitRef || repoMetadata?.default_branch })}
              />
              <NavMenuItem
                data-code-repo-section="commits"
                isSelected={isCommitSelected}
                isSubLink
                label={getString('commits')}
                to={routes.toCODECommits({ repoPath, commitRef: _gitRef })}
              />
              <NavMenuItem
                data-code-repo-section="branches"
                isSubLink
                label={getString('branches')}
                to={routes.toCODEBranches({ repoPath })}
              />
              <NavMenuItem
                data-code-repo-section="tags"
                isSubLink
                label={getString('tags')}
                to={routes.toCODETags({ repoPath })}
              />
              <NavMenuItem
                data-code-repo-section="pull-requests"
                isSubLink
                label={getString('pullRequests')}
                to={routes.toCODEPullRequests({ repoPath })}
              />
              <NavMenuItem
                data-code-repo-section="branches"
                isSubLink
                isSelected={isWebhookSelected}
                label={getString('webhooks')}
                to={routes.toCODEWebhooks({ repoPath })}
              />
              {standalone && (
                <NavMenuItem
                  data-code-repo-section="pipelines"
                  isSubLink
                  label={getString('pageTitle.pipelines')}
                  to={routes.toCODEPipelines({ repoPath })}
                />
              )}
              <NavMenuItem
                data-code-repo-section="settings"
                isSubLink
                label={getString('manageRepository')}
                to={routes.toCODESettings({ repoPath })}
              />
            </Layout.Vertical>
          </Container>
        </Render>

        {/* MCP Setup */}
        <Render when={!!activeSpacePath}>
          <div className={css.sectionDivider} />
          <div className={css.sectionLabel}>TOOLS</div>
          <NavMenuItem
            label="MCP Setup"
            isSelected={isMcpSetupSelected}
            to={
              routes.toSOLODEVMcpSetup
                ? routes.toSOLODEVMcpSetup({ space: activeSpacePath })
                : routes.toCODERepositories({ space: activeSpacePath })
            }
          />
        </Render>

        {/* Settings & Secrets */}
        {standalone && (
          <Render when={!!activeSpacePath}>
            <NavMenuItem label={getString('pageTitle.secrets')} to={routes.toCODESecrets({ space: activeSpacePath })} />
          </Render>
        )}

        <Render when={!!activeSpacePath}>
          <NavMenuItem
            customIcon={<Settings />}
            label={getString('settings')}
            to={routes.toCODESpaceSettings({ space: activeSpacePath })}
          />
        </Render>
      </Layout.Vertical>
    </Container>
  )
}
