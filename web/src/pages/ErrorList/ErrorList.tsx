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
import { Container, FlexExpander, Layout, PageBody, PageHeader, TableV2 as Table, Text, Tag } from '@harnessio/uicore'
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
import css from './ErrorList.module.scss'

interface ErrorGroup {
  id: string
  title: string
  message?: string
  frequency: number
  lastOccurrence: number
  firstOccurrence: number
  status: 'open' | 'resolved' | 'ignored'
}

const getStatusIntent = (status: string): Intent => {
  switch (status) {
    case 'resolved':
      return Intent.SUCCESS
    case 'ignored':
      return Intent.NONE
    case 'open':
      return Intent.DANGER
    default:
      return Intent.NONE
  }
}

const ErrorList = () => {
  const space = useGetSpaceParam()
  const { getString } = useStrings()
  const [searchTerm, setSearchTerm] = useState<string | undefined>()
  const pageBrowser = useQueryParams<PageBrowserProps>()
  const pageInit = pageBrowser.page ? parseInt(pageBrowser.page) : 1
  const [page, setPage] = usePageIndex(pageInit)

  const {
    data: errors,
    error,
    loading,
    refetch,
    response
  } = useGet<ErrorGroup[]>({
    path: `/api/v1/spaces/${space}/+/errors`,
    queryParams: { page, limit: LIST_FETCHING_LIMIT, query: searchTerm }
  })

  const columns: Column<ErrorGroup>[] = useMemo(
    () => [
      {
        Header: getString('name') || 'Error Title',
        width: 'calc(100% - 360px)',
        Cell: ({ row }: CellProps<ErrorGroup>) => {
          const record = row.original
          return (
            <Container className={css.nameContainer}>
              <Layout.Horizontal spacing="small" style={{ flexGrow: 1 }}>
                <Layout.Vertical flex className={css.name}>
                  <Text className={css.errorTitle} lineClamp={1}>
                    <Keywords value={searchTerm}>{record.title}</Keywords>
                  </Text>
                  {record.message && (
                    <Text className={css.desc} lineClamp={1}>
                      {record.message}
                    </Text>
                  )}
                </Layout.Vertical>
              </Layout.Horizontal>
            </Container>
          )
        }
      },
      {
        Header: 'Occurrences',
        width: '110px',
        Cell: ({ row }: CellProps<ErrorGroup>) => {
          const count = row.original.frequency
          return (
            <Text
              color={count > 10 ? Color.RED_600 : count > 5 ? Color.ORANGE_600 : Color.BLACK}
              lineClamp={1}
              font={{ weight: 'bold' }}>
              {count}
            </Text>
          )
        },
        disableSortBy: true
      },
      {
        Header: 'Status',
        width: '100px',
        Cell: ({ row }: CellProps<ErrorGroup>) => {
          return (
            <Tag intent={getStatusIntent(row.original.status)} minimal round>
              {row.original.status.toUpperCase()}
            </Tag>
          )
        },
        disableSortBy: true
      },
      {
        Header: getString('updatedDate') || 'Last Seen',
        width: '140px',
        Cell: ({ row }: CellProps<ErrorGroup>) => {
          return (
            <Layout.Horizontal style={{ alignItems: 'center' }}>
              <Text color={Color.BLACK} lineClamp={1} width={120}>
                {formatDate(row.original.lastOccurrence)}
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
      <PageHeader title="Error Tracker" />
      <PageBody
        className={cx({ [css.withError]: !!error })}
        error={error ? getErrorMessage(error) : null}
        retryOnError={voidFn(refetch)}
        noData={{
          when: () => errors?.length === 0 && searchTerm === undefined,
          image: noDataImage,
          message: 'No errors tracked'
        }}>
        <LoadingSpinner visible={loading && !searchTerm} />

        <Container padding="xlarge">
          <Layout.Horizontal spacing="large" className={css.layout}>
            <FlexExpander />
            <SearchInputWithSpinner loading={loading} query={searchTerm} setQuery={setSearchTerm} />
          </Layout.Horizontal>

          <Container margin={{ top: 'medium' }}>
            {!!errors?.length && (
              <Table<ErrorGroup>
                className={css.table}
                columns={columns}
                data={errors || []}
                getRowClassName={row => cx(css.row, !row.original.message && css.noDesc)}
              />
            )}
            <NoResultCard showWhen={() => !!errors && errors?.length === 0 && !!searchTerm?.length} forSearch={true} />
          </Container>
          <ResourceListingPagination response={response} page={page} setPage={setPage} />
        </Container>
      </PageBody>
    </Container>
  )
}

export default ErrorList
