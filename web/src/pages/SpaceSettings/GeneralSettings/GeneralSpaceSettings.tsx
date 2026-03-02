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
import cx from 'classnames'
import { Button, ButtonSize, ButtonVariation, Container, Text, TextInput, useToaster } from '@harnessio/uicore'
import { Color, Intent } from '@harnessio/design-system'
import { useMutate } from 'restful-react'
import { useHistory } from 'react-router-dom'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import { useGetSpace } from 'services/code'
import { useAppContext } from 'AppContext'
import { useStrings } from 'framework/strings'
import { ACCESS_MODES, getErrorMessage, permissionProps } from 'utils/Utils'
import useDeleteSpaceModal from '../DeleteSpaceModal/DeleteSpaceModal'
import css from './GeneralSpaceSettings.module.scss'

export default function GeneralSpaceSettings() {
  const { space } = useGetRepositoryMetadata()
  const { routes, standalone, hooks } = useAppContext()
  const history = useHistory()
  const { showError, showSuccess } = useToaster()
  const { getString } = useStrings()
  const { openModal: openDeleteSpaceModal } = useDeleteSpaceModal()
  const { data, refetch } = useGetSpace({ space_ref: encodeURIComponent(space), lazy: !space })
  const [editName, setEditName] = useState(ACCESS_MODES.VIEW)
  const [editDesc, setEditDesc] = useState(ACCESS_MODES.VIEW)
  const [draftName, setDraftName] = useState('')
  const [draftDesc, setDraftDesc] = useState('')

  const currentName = data?.identifier || space
  const currentDescription = data?.description || ''

  useEffect(() => {
    setDraftName(currentName)
    setDraftDesc(currentDescription)
  }, [currentDescription, currentName])

  const { mutate: patchSpace } = useMutate({
    verb: 'PATCH',
    path: `/api/v1/spaces/${space}`
  })
  const { mutate: updateName } = useMutate({
    verb: 'POST',
    path: `/api/v1/spaces/${space}/move`
  })

  const permEditResult = hooks?.usePermissionTranslate?.(
    {
      resource: {
        resourceType: 'CODE_REPOSITORY'
      },
      permissions: ['code_repo_edit']
    },
    [space]
  )
  const permDeleteResult = hooks?.usePermissionTranslate?.(
    {
      resource: {
        resourceType: 'CODE_REPOSITORY'
      },
      permissions: ['code_repo_delete']
    },
    [space]
  )

  const dashboardHref = useMemo(() => {
    return currentName && routes.toSOLODEVDashboard ? routes.toSOLODEVDashboard({ space: currentName }) : ''
  }, [currentName, routes])

  const settingsHref = useMemo(() => {
    return currentName ? routes.toCODESpaceSettings({ space: currentName }) : ''
  }, [currentName, routes])

  const saveName = () => {
    const nextName = draftName.trim()

    if (!nextName) {
      showError('Enter a project name.')
      return
    }

    updateName({ uid: nextName })
      .then(() => {
        showSuccess(getString('spaceUpdate'))
        setEditName(ACCESS_MODES.VIEW)
        history.push(routes.toCODESpaceSettings({ space: nextName }))
      })
      .catch(err => {
        showError(getErrorMessage(err))
      })
  }

  const saveDescription = () => {
    patchSpace({ description: draftDesc.trim() })
      .then(() => {
        showSuccess(getString('spaceUpdate'))
        setEditDesc(ACCESS_MODES.VIEW)
        refetch()
      })
      .catch(err => {
        showError(getErrorMessage(err))
      })
  }

  return (
    <Container className={css.mainCtn}>
      <section className={css.hero}>
        <div className={css.heroCopy}>
          <span className={css.eyebrow}>Project identity</span>
          <h1 className={css.title}>Shape the SoloDev control surface around this project.</h1>
          <p className={css.subtitle}>
            Clean up naming, keep the route structure legible, and isolate destructive actions so the settings page
            feels deliberate instead of inherited.
          </p>
        </div>
        <div className={css.heroMeta}>
          <div className={css.metaCard}>
            <span className={css.metaLabel}>Project slug</span>
            <span className={css.metaValue}>{currentName}</span>
          </div>
          <div className={css.metaCard}>
            <span className={css.metaLabel}>Dashboard route</span>
            <code className={css.metaCode}>{dashboardHref || '-'}</code>
          </div>
          <div className={css.metaCard}>
            <span className={css.metaLabel}>Settings route</span>
            <code className={css.metaCode}>{settingsHref || '-'}</code>
          </div>
        </div>
      </section>

      <section className={css.contentGrid}>
        <div className={css.card}>
          <div className={css.cardHeader}>
            <div>
              <h2 className={css.cardTitle}>Project profile</h2>
              <p className={css.cardText}>Rename the project, sharpen the description, and keep the navigation crisp.</p>
            </div>
          </div>

          <div className={css.settingRow}>
            <div className={css.settingIntro}>
              <h3 className={css.settingTitle}>{getString('name')}</h3>
              <p className={css.settingHelp}>Used in the dashboard header, route structure, and navigation rail.</p>
            </div>
            <div className={css.settingBody}>
              {editName === ACCESS_MODES.EDIT ? (
                <div className={css.editBlock}>
                  <TextInput
                    name="name"
                    value={draftName}
                    className={css.input}
                    onChange={evt => {
                      setDraftName((evt.currentTarget as HTMLInputElement).value)
                    }}
                  />
                  <div className={css.buttonRow}>
                    <Button
                      type="button"
                      text={getString('save')}
                      variation={ButtonVariation.PRIMARY}
                      size={ButtonSize.SMALL}
                      disabled={!draftName.trim() || draftName.trim() === currentName}
                      onClick={saveName}
                      {...permissionProps(permEditResult, standalone)}
                    />
                    <Button
                      type="button"
                      text={getString('cancel')}
                      variation={ButtonVariation.TERTIARY}
                      size={ButtonSize.SMALL}
                      onClick={() => {
                        setDraftName(currentName)
                        setEditName(ACCESS_MODES.VIEW)
                      }}
                    />
                  </div>
                </div>
              ) : (
                <div className={css.readOnlyBlock}>
                  <div>
                    <div className={css.readOnlyValue}>{currentName}</div>
                    <div className={css.readOnlyHint}>Project slug stays readable across SoloDev URLs.</div>
                  </div>
                  <Button
                    type="button"
                    text={getString('edit')}
                    icon="Edit"
                    variation={ButtonVariation.SECONDARY}
                    size={ButtonSize.SMALL}
                    onClick={() => setEditName(ACCESS_MODES.EDIT)}
                    {...permissionProps(permEditResult, standalone)}
                  />
                </div>
              )}
            </div>
          </div>

          <div className={css.divider} />

          <div className={css.settingRow}>
            <div className={css.settingIntro}>
              <h3 className={css.settingTitle}>{getString('description')}</h3>
              <p className={css.settingHelp}>Give operators a quick sentence about the mission, scope, or ownership.</p>
            </div>
            <div className={css.settingBody}>
              {editDesc === ACCESS_MODES.EDIT ? (
                <div className={css.editBlock}>
                  <textarea
                    value={draftDesc}
                    className={css.textArea}
                    placeholder="Describe what this project owns and why it matters."
                    onChange={evt => {
                      setDraftDesc(evt.currentTarget.value)
                    }}
                  />
                  <div className={css.buttonRow}>
                    <Button
                      type="button"
                      text={getString('save')}
                      variation={ButtonVariation.PRIMARY}
                      size={ButtonSize.SMALL}
                      disabled={draftDesc.trim() === currentDescription.trim()}
                      onClick={saveDescription}
                      {...permissionProps(permEditResult, standalone)}
                    />
                    <Button
                      type="button"
                      text={getString('cancel')}
                      variation={ButtonVariation.TERTIARY}
                      size={ButtonSize.SMALL}
                      onClick={() => {
                        setDraftDesc(currentDescription)
                        setEditDesc(ACCESS_MODES.VIEW)
                      }}
                    />
                  </div>
                </div>
              ) : (
                <div className={css.readOnlyBlock}>
                  <div>
                    <div className={cx(css.readOnlyValue, css.descriptionValue)}>
                      {currentDescription || 'No project description yet.'}
                    </div>
                    <div className={css.readOnlyHint}>A strong description turns the settings view into an actual operating brief.</div>
                  </div>
                  <Button
                    type="button"
                    text={getString('edit')}
                    icon="Edit"
                    variation={ButtonVariation.SECONDARY}
                    size={ButtonSize.SMALL}
                    onClick={() => setEditDesc(ACCESS_MODES.EDIT)}
                    {...permissionProps(permEditResult, standalone)}
                  />
                </div>
              )}
            </div>
          </div>
        </div>

        <div className={css.card}>
          <div className={css.cardHeader}>
            <div>
              <h2 className={css.cardTitle}>Quick routes</h2>
              <p className={css.cardText}>Jump straight into the places you are most likely to tune after renaming the project.</p>
            </div>
          </div>

          <button
            type="button"
            className={css.quickLink}
            onClick={() => {
              if (currentName && routes.toSOLODEVDashboard) {
                history.push(routes.toSOLODEVDashboard({ space: currentName }))
              }
            }}>
            <span className={css.quickLinkTitle}>Return to dashboard</span>
            <span className={css.quickLinkText}>Back to the SoloDev mission control surface for this project.</span>
          </button>

          <button
            type="button"
            className={css.quickLink}
            onClick={() => {
              if (currentName && routes.toSOLODEVMcpSetup) {
                history.push(routes.toSOLODEVMcpSetup({ space: currentName }))
              }
            }}>
            <span className={css.quickLinkTitle}>Open MCP setup</span>
            <span className={css.quickLinkText}>Reconfigure Claude, Cursor, or another MCP-native client.</span>
          </button>

          <div className={css.noteCard}>
            <span className={css.noteLabel}>Theme control</span>
            <p className={css.noteText}>
              The new light and dark toggle lives in the left rail so this project can move from studio-bright to
              command-room dark without leaving the page.
            </p>
          </div>
        </div>
      </section>

      <section className={cx(css.card, css.dangerCard)}>
        <div className={css.dangerHeader}>
          <div>
            <h2 className={css.cardTitle}>Danger zone</h2>
            <p className={css.cardText}>Delete the project only when you are certain nothing downstream still depends on it.</p>
          </div>
          <Button
            type="button"
            intent={Intent.DANGER}
            variation={ButtonVariation.SECONDARY}
            text={getString('deleteSpace')}
            onClick={() => openDeleteSpaceModal()}
            {...permissionProps(permDeleteResult, standalone)}
          />
        </div>

        <div className={css.warningBanner}>
          <Text
            icon="main-issue"
            iconProps={{ size: 16, color: Color.ORANGE_500, margin: { right: 'small' } }}
            color={Color.ORANGE_900}>
            {getString('spaceSetting.intentText', {
              space: currentName
            })}
          </Text>
        </div>
      </section>
    </Container>
  )
}
