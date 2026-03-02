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
import {
  Container,
  Layout,
  Text,
  Button,
  ButtonVariation,
  Formik,
  useToaster,
  FormInput,
  FormikForm
} from '@harnessio/uicore'
import { useGet, useMutate } from 'restful-react'
import type { FormikState } from 'formik'
import type { RepoRepositoryOutput } from 'services/code'
import { useStrings } from 'framework/strings'
import { useAppContext } from 'AppContext'
import { useGetSpaceParam } from 'hooks/useGetSpaceParam'
import { getErrorMessage, permissionProps } from 'utils/Utils'
import { NavigationCheck } from 'components/NavigationCheck/NavigationCheck'
import { LoadingSpinner } from 'components/LoadingSpinner/LoadingSpinner'
import css from './SecurityScanSettings.module.scss'

interface SecurityScanProps {
  repoMetadata: RepoRepositoryOutput | undefined
  activeTab: string
}

interface FormData {
  secretScanEnable: boolean
  verifyCommitterIdentity: boolean
  findingRemediationMode: SecurityFindingRemediationMode
}

type SecurityFindingRemediationMode = 'manual' | 'critical_high_auto' | 'all_auto'

const remediationModeOptions: Array<{
  value: SecurityFindingRemediationMode
  label: string
  description: string
}> = [
  {
    value: 'manual',
    label: 'Manual only',
    description: 'Security findings are recorded, but fix tasks are only created when a user triggers them.'
  },
  {
    value: 'critical_high_auto',
    label: 'Auto-fix critical/high',
    description: 'SoloDev creates fix tasks automatically for critical and high severity findings only.'
  },
  {
    value: 'all_auto',
    label: 'Auto-fix all findings',
    description: 'SoloDev creates fix tasks automatically for every open finding in this repository.'
  }
]

const SecurityScanSettings = (props: SecurityScanProps) => {
  const { repoMetadata, activeTab } = props
  const { hooks, standalone, routingId } = useAppContext()
  const { getString } = useStrings()
  const { showError, showSuccess } = useToaster()
  const space = useGetSpaceParam()
  const permPushResult = hooks?.usePermissionTranslate?.(
    {
      resource: {
        resourceType: 'CODE_REPOSITORY',
        resourceIdentifier: repoMetadata?.identifier as string
      },
      permissions: ['code_repo_edit']
    },
    [space]
  )
  const { data: securitySettings, loading: securitySettingsLoading } = useGet({
    path: `/api/v1/repos/${repoMetadata?.path}/+/settings/security`,
    queryParams: { routingId: routingId },
    lazy: !activeTab
  })
  const { mutate: updateSecuritySettings, loading: isUpdating } = useMutate({
    verb: 'PATCH',
    path: `/api/v1/repos/${repoMetadata?.path}/+/settings/security`,
    queryParams: { routingId: routingId }
  })

  const handleSubmit = async (
    formData: FormData,
    resetForm: (nextState?: Partial<FormikState<FormData>> | undefined) => void
  ) => {
    try {
      const payload = {
        secret_scanning_enabled: !!formData?.secretScanEnable,
        principal_committer_match: !!formData?.verifyCommitterIdentity,
        finding_remediation_mode: formData?.findingRemediationMode
      }
      const response = await updateSecuritySettings(payload)
      showSuccess(getString('securitySettings.updateSuccess'), 1500)
      resetForm({
        values: {
          secretScanEnable: !!response?.secret_scanning_enabled,
          verifyCommitterIdentity: !!response?.principal_committer_match,
          findingRemediationMode: response?.finding_remediation_mode || 'manual'
        }
      })
    } catch (exception) {
      showError(getErrorMessage(exception), 1500, getString('securitySettings.failedToUpdate'))
    }
  }
  return (
    <Container className={css.main}>
      <LoadingSpinner visible={securitySettingsLoading || isUpdating} />
      {securitySettings && (
        <Formik<FormData>
          formName="securityScanSettings"
          initialValues={{
            secretScanEnable: !!securitySettings?.secret_scanning_enabled,
            verifyCommitterIdentity: !!securitySettings?.principal_committer_match,
            findingRemediationMode: securitySettings?.finding_remediation_mode || 'manual'
          }}
          onSubmit={(formData, { resetForm }) => {
            handleSubmit(formData, resetForm)
          }}>
          {formik => {
            return (
              <FormikForm>
                <Layout.Vertical padding={{ top: 'medium' }}>
                  <Container padding="medium" margin="medium" className={css.generalContainer}>
                    <Layout.Horizontal
                      spacing={'medium'}
                      padding={{ left: 'medium' }}
                      flex={{ alignItems: 'center', justifyContent: 'flex-start' }}>
                      <FormInput.Toggle
                        {...permissionProps(permPushResult, standalone)}
                        key={'secretScanEnable'}
                        style={{ margin: '0px' }}
                        label=""
                        name="secretScanEnable"></FormInput.Toggle>
                      <Text className={css.title}>{getString('securitySettings.scanningSecret')}</Text>
                      <Text className={css.text}>{getString('securitySettings.scanningSecretDesc')}</Text>
                    </Layout.Horizontal>
                  </Container>
                  <Container padding="medium" margin="medium" className={css.generalContainer}>
                    <Layout.Horizontal
                      spacing={'medium'}
                      padding={{ left: 'medium' }}
                      flex={{ alignItems: 'center', justifyContent: 'flex-start' }}>
                      <FormInput.Toggle
                        {...permissionProps(permPushResult, standalone)}
                        key={'verifyCommitterIdentity'}
                        style={{ margin: '0px' }}
                        label=""
                        name="verifyCommitterIdentity"></FormInput.Toggle>
                      <Text className={css.title}>{getString('securitySettings.verifyCommitterIdentity')}</Text>
                      <Text className={css.text}>{getString('securitySettings.verifyCommitterIdentityDesc')}</Text>
                    </Layout.Horizontal>
                  </Container>
                  <Container padding="medium" margin="medium" className={css.generalContainer}>
                    <Layout.Vertical spacing="small" padding={{ left: 'medium' }}>
                      <Text className={css.title}>Security remediation mode</Text>
                      <Text className={css.text}>
                        Choose when SoloDev should create AI fix tasks from new security findings for this repository.
                      </Text>
                      <Container margin={{ top: 'medium' }}>
                        <FormInput.RadioGroup
                          {...permissionProps(permPushResult, standalone)}
                          name="findingRemediationMode"
                          key={formik.values.findingRemediationMode}
                          label=""
                          className={css.radioContainer}
                          items={remediationModeOptions.map(option => ({
                            label: (
                              <Container>
                                <Layout.Vertical spacing="xsmall">
                                  <Text className={css.optionTitle}>{option.label}</Text>
                                  <Text className={css.optionText}>{option.description}</Text>
                                </Layout.Vertical>
                              </Container>
                            ),
                            value: option.value
                          }))}
                        />
                      </Container>
                    </Layout.Vertical>
                  </Container>
                </Layout.Vertical>
                <Layout.Horizontal margin={'medium'} spacing={'medium'}>
                  <Button
                    variation={ButtonVariation.PRIMARY}
                    text={getString('save')}
                    onClick={() => formik.submitForm()}
                    disabled={formik.isSubmitting}
                    {...permissionProps(permPushResult, standalone)}
                  />
                </Layout.Horizontal>
                <NavigationCheck when={formik.dirty} />
              </FormikForm>
            )
          }}
        </Formik>
      )}
    </Container>
  )
}

export default SecurityScanSettings
