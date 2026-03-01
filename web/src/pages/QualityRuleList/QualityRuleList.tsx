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
  Button,
  ButtonVariation,
  Container,
  FlexExpander,
  Layout,
  PageBody,
  PageHeader,
  TableV2 as Table,
  Text,
  Tag,
  useToaster
} from '@harnessio/uicore'
import { Color, Intent } from '@harnessio/design-system'
import cx from 'classnames'
import type { CellProps, Column } from 'react-table'
import Keywords from 'react-keywords'
import { useGet } from 'restful-react'
import { useStrings } from 'framework/strings'
import { LoadingSpinner } from 'components/LoadingSpinner/LoadingSpinner'
import { SearchInputWithSpinner } from 'components/SearchInputWithSpinner/SearchInputWithSpinner'
import { NoResultCard } from 'components/NoResultCard/NoResultCard'
import { LIST_FETCHING_LIMIT, PageBrowserProps, formatDate, getErrorMessage, voidFn } from 'utils/Utils'
import { usePageIndex } from 'hooks/usePageIndex'
import { useQueryParams } from 'hooks/useQueryParams'
import { useGetSpaceParam } from 'hooks/useGetSpaceParam'
import { ResourceListingPagination } from 'components/ResourceListingPagination/ResourceListingPagination'
import noDataImage from '../RepositoriesListing/no-repo.svg?url'
import css from './QualityRuleList.module.scss'

interface QualityRule {
  id: string
  name: string
  description?: string
  enforcement: 'required' | 'optional' | 'warning'
  category: string
  enabled: boolean
  createdAt: number
  updatedAt: number
}

const getEnforcementIntent = (enforcement: string): Intent => {
  switch (enforcement) {
    case 'required':
      return Intent.DANGER
    case 'warning':
      return Intent.WARNING
    case 'optional':
      return Intent.PRIMARY
    default:
      return Intent.PRIMARY
  }
}

const QualityRuleList = () => {
  const space = useGetSpaceParam()
  const { getString } = useStrings()
  const { showSuccess } = useToaster()
  const [searchTerm, setSearchTerm] = useState<string | undefined>()
  const pageBrowser = useQueryParams<PageBrowserProps>()
  const pageInit = pageBrowser.page ? parseInt(pageBrowser.page) : 1
  const [page, setPage] = usePageIndex(pageInit)

  const {
    data: rules,
    error,
    loading,
    refetch,
    response
  } = useGet<QualityRule[]>({
    path: `/api/v1/spaces/${space}/+/quality/rules`,
    queryParams: { page, limit: LIST_FETCHING_LIMIT, query: searchTerm }
  })

  const NewRuleButton = (
    <Button
      text="New Quality Rule"
      variation={ButtonVariation.PRIMARY}
      icon="plus"
      onClick={() => {
        showSuccess('Create quality rule functionality to be implemented')
      }}
    />
  )

  const columns: Column<QualityRule>[] = useMemo(
    () => [
      {
        Header: getString('name') || 'Rule Name',
        width: 'calc(100% - 340px)',
        Cell: ({ row }: CellProps<QualityRule>) => {
          const record = row.original
          return (
            <Container className={css.nameContainer}>
              <Layout.Horizontal spacing="small" style={{ flexGrow: 1 }}>
                <Layout.Vertical flex className={css.name}>
                  <Text className={css.ruleName} lineClamp={1}>
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
        Header: 'Category',
        width: '100px',
        Cell: ({ row }: CellProps<QualityRule>) => {
          return (
            <Text color={Color.BLACK} lineClamp={1}>
              {row.original.category}
            </Text>
          )
        },
        disableSortBy: true
      },
      {
        Header: 'Enforcement',
        width: '110px',
        Cell: ({ row }: CellProps<QualityRule>) => {
          return (
            <Tag intent={getEnforcementIntent(row.original.enforcement)} minimal round>
              {row.original.enforcement.toUpperCase()}
            </Tag>
          )
        },
        disableSortBy: true
      },
      {
        Header: 'Status',
        width: '80px',
        Cell: ({ row }: CellProps<QualityRule>) => {
          return (
            <Text
              color={row.original.enabled ? Color.GREEN_600 : Color.GREY_600}
              lineClamp={1}
              font={{ weight: 'bold' }}>
              {row.original.enabled ? 'Enabled' : 'Disabled'}
            </Text>
          )
        },
        disableSortBy: true
      },
      {
        Header: getString('updatedDate') || 'Updated',
        width: '140px',
        Cell: ({ row }: CellProps<QualityRule>) => {
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
    [getString, searchTerm]
  )

  return (
    <Container className={css.main}>
      <PageHeader title="Quality Gates" />
      <PageBody
        className={cx({ [css.withError]: !!error })}
        error={error ? getErrorMessage(error) : null}
        retryOnError={voidFn(refetch)}
        noData={{
          when: () => rules?.length === 0 && searchTerm === undefined,
          image: noDataImage,
          message: 'No quality rules configured',
          button: NewRuleButton
        }}>
        <LoadingSpinner visible={loading && !searchTerm} />

        <Container padding="xlarge">
          <Layout.Horizontal spacing="large" className={css.layout}>
            {NewRuleButton}
            <FlexExpander />
            <SearchInputWithSpinner loading={loading} query={searchTerm} setQuery={setSearchTerm} />
          </Layout.Horizontal>

          <Container margin={{ top: 'medium' }}>
            {!!rules?.length && (
              <Table<QualityRule>
                className={css.table}
                columns={columns}
                data={rules || []}
                getRowClassName={row => cx(css.row, !row.original.description && css.noDesc)}
              />
            )}
            <NoResultCard showWhen={() => !!rules && rules?.length === 0 && !!searchTerm?.length} forSearch={true} />
          </Container>
          <ResourceListingPagination response={response} page={page} setPage={setPage} />
        </Container>
      </PageBody>
    </Container>
  )
}

export default QualityRuleList
