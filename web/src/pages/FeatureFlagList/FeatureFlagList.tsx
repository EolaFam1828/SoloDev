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

import React, { useMemo, useState } from 'react'
import {
  ButtonVariation,
  Container,
  FlexExpander,
  Layout,
  PageBody,
  PageHeader,
  Switch,
  TableV2 as Table,
  Text,
  useToaster
} from '@harnessio/uicore'
import { Color } from '@harnessio/design-system'
import cx from 'classnames'
import type { CellProps, Column } from 'react-table'
import Keywords from 'react-keywords'
import { useGet, useMutate } from 'restful-react'
import { useStrings } from 'framework/strings'
import { LoadingSpinner } from 'components/LoadingSpinner/LoadingSpinner'
import { SearchInputWithSpinner } from 'components/SearchInputWithSpinner/SearchInputWithSpinner'
import { NoResultCard } from 'components/NoResultCard/NoResultCard'
import { LIST_FETCHING_LIMIT, PageBrowserProps, formatDate, getErrorMessage, voidFn } from 'utils/Utils'
import { usePageIndex } from 'hooks/usePageIndex'
import { useQueryParams } from 'hooks/useQueryParams'
import { useGetSpaceParam } from 'hooks/useGetSpaceParam'
import { ResourceListingPagination } from 'components/ResourceListingPagination/ResourceListingPagination'
import { Button } from '@harnessio/uicore'
import noDataImage from '../RepositoriesListing/no-repo.svg?url'
import css from './FeatureFlagList.module.scss'

interface FeatureFlag {
  id: string
  name: string
  description?: string
  enabled: boolean
  createdAt: number
  updatedAt: number
}

const FeatureFlagList = () => {
  const space = useGetSpaceParam()
  const { getString } = useStrings()
  const { showSuccess, showError } = useToaster()
  const [searchTerm, setSearchTerm] = useState<string | undefined>()
  const pageBrowser = useQueryParams<PageBrowserProps>()
  const pageInit = pageBrowser.page ? parseInt(pageBrowser.page) : 1
  const [page, setPage] = usePageIndex(pageInit)

  const {
    data: flags,
    error,
    loading,
    refetch,
    response
  } = useGet<FeatureFlag[]>({
    path: `/api/v1/spaces/${space}/+/featureflags`,
    queryParams: { page, limit: LIST_FETCHING_LIMIT, query: searchTerm }
  })

  const NewFlagButton = (
    <Button
      text="New Feature Flag"
      variation={ButtonVariation.PRIMARY}
      icon="plus"
      onClick={() => {
        // Navigate to new flag creation page or open modal
        showSuccess('Create feature flag functionality to be implemented')
      }}
    />
  )

  const { mutate: toggleFlag } = useMutate({
    verb: 'PUT',
    path: `/api/v1/spaces/${space}/+/featureflags/:id`
  })

  const columns: Column<FeatureFlag>[] = useMemo(
    () => [
      {
        Header: getString('name') || 'Name',
        width: 'calc(100% - 260px)',
        Cell: ({ row }: CellProps<FeatureFlag>) => {
          const record = row.original
          return (
            <Container className={css.nameContainer}>
              <Layout.Horizontal spacing="small" style={{ flexGrow: 1 }}>
                <Layout.Vertical flex className={css.name}>
                  <Text className={css.flagName} lineClamp={1}>
                    <Keywords value={searchTerm}>{record.name}</Keywords>
                  </Text>
                  {record.description && (
                    <Text className={css.desc} lineClamp={1}>
                      {record.description}
                    </Text>
                  )}
                </Layout.Vertical>
              </Layout.Horizontal>
            </Container>
          )
        }
      },
      {
        Header: 'Status',
        width: '80px',
        Cell: ({ row }: CellProps<FeatureFlag>) => {
          const [isEnabled, setIsEnabled] = useState(row.original.enabled)
          return (
            <Switch
              checked={isEnabled}
              onChange={checked => {
                setIsEnabled(checked)
                toggleFlag({ id: row.original.id }, { enabled: checked })
                  .then(() => {
                    showSuccess('Feature flag updated')
                    refetch()
                  })
                  .catch(err => {
                    setIsEnabled(!checked)
                    showError(getErrorMessage(err))
                  })
              }}
            />
          )
        },
        disableSortBy: true
      },
      {
        Header: getString('updatedDate') || 'Updated',
        width: '180px',
        Cell: ({ row }: CellProps<FeatureFlag>) => {
          return (
            <Layout.Horizontal style={{ alignItems: 'center' }}>
              <Text color={Color.BLACK} lineClamp={1} width={120}>
                {formatDate(row.original.updatedAt)}
              </Text>
            </Layout.Horizontal>
          )
        },
        disableSortBy: true
      }
    ],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [getString, refetch, searchTerm, space]
  )

  return (
    <Container className={css.main}>
      <PageHeader title="Feature Flags" />
      <PageBody
        className={cx({ [css.withError]: !!error })}
        error={error ? getErrorMessage(error) : null}
        retryOnError={voidFn(refetch)}
        noData={{
          when: () => flags?.length === 0 && searchTerm === undefined,
          image: noDataImage,
          message: 'No feature flags found',
          button: NewFlagButton
        }}>
        <LoadingSpinner visible={loading && !searchTerm} />

        <Container padding="xlarge">
          <Layout.Horizontal spacing="large" className={css.layout}>
            {NewFlagButton}
            <FlexExpander />
            <SearchInputWithSpinner loading={loading} query={searchTerm} setQuery={setSearchTerm} />
          </Layout.Horizontal>

          <Container margin={{ top: 'medium' }}>
            {!!flags?.length && (
              <Table<FeatureFlag>
                className={css.table}
                columns={columns}
                data={flags || []}
                getRowClassName={row => cx(css.row, !row.original.description && css.noDesc)}
              />
            )}
            <NoResultCard
              showWhen={() => !!flags && flags?.length === 0 && !!searchTerm?.length}
              forSearch={true}
            />
          </Container>
          <ResourceListingPagination response={response} page={page} setPage={setPage} />
        </Container>
      </PageBody>
    </Container>
  )
}

export default FeatureFlagList
