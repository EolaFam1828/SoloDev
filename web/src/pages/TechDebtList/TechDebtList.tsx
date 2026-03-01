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
import { Button } from '@harnessio/uicore'
import noDataImage from '../RepositoriesListing/no-repo.svg?url'
import css from './TechDebtList.module.scss'

interface TechDebtItem {
  id: string
  title: string
  description?: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  status: 'open' | 'in_progress' | 'resolved'
  createdAt: number
  updatedAt: number
}

const getSeverityColor = (severity: string): Intent => {
  switch (severity) {
    case 'critical':
      return Intent.DANGER
    case 'high':
      return Intent.WARNING
    case 'medium':
      return Intent.PRIMARY
    case 'low':
      return Intent.SUCCESS
    default:
      return Intent.PRIMARY
  }
}

const getTechDebtStatusColor = (status: string): string => {
  switch (status) {
    case 'resolved':
      return Color.GREEN_600
    case 'in_progress':
      return Color.YELLOW_600
    case 'open':
      return Color.RED_600
    default:
      return Color.GREY_600
  }
}

const TechDebtList = () => {
  const space = useGetSpaceParam()
  const { getString } = useStrings()
  const { showSuccess } = useToaster()
  const [searchTerm, setSearchTerm] = useState<string | undefined>()
  const pageBrowser = useQueryParams<PageBrowserProps>()
  const pageInit = pageBrowser.page ? parseInt(pageBrowser.page) : 1
  const [page, setPage] = usePageIndex(pageInit)

  const {
    data: items,
    error,
    loading,
    refetch,
    response
  } = useGet<TechDebtItem[]>({
    path: `/api/v1/spaces/${space}/+/techdebt`,
    queryParams: { page, limit: LIST_FETCHING_LIMIT, query: searchTerm }
  })

  const NewItemButton = (
    <Button
      text="New Tech Debt Item"
      variation={ButtonVariation.PRIMARY}
      icon="plus"
      onClick={() => {
        showSuccess('Create tech debt item functionality to be implemented')
      }}
    />
  )

  const columns: Column<TechDebtItem>[] = useMemo(
    () => [
      {
        Header: getString('name') || 'Title',
        width: 'calc(100% - 340px)',
        Cell: ({ row }: CellProps<TechDebtItem>) => {
          const record = row.original
          return (
            <Container className={css.nameContainer}>
              <Layout.Horizontal spacing="small" style={{ flexGrow: 1 }}>
                <Layout.Vertical flex className={css.name}>
                  <Text className={css.itemTitle} lineClamp={1}>
                    <Keywords value={searchTerm}>{record.title}</Keywords>
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
        Header: 'Severity',
        width: '100px',
        Cell: ({ row }: CellProps<TechDebtItem>) => {
          return (
            <Tag
              intent={getSeverityColor(row.original.severity)}
              minimal
              round>
              {row.original.severity.toUpperCase()}
            </Tag>
          )
        },
        disableSortBy: true
      },
      {
        Header: 'Status',
        width: '100px',
        Cell: ({ row }: CellProps<TechDebtItem>) => {
          return (
            <Text color={getTechDebtStatusColor(row.original.status)} lineClamp={1}>
              {row.original.status.replace('_', ' ').toUpperCase()}
            </Text>
          )
        },
        disableSortBy: true
      },
      {
        Header: getString('updatedDate') || 'Updated',
        width: '140px',
        Cell: ({ row }: CellProps<TechDebtItem>) => {
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
      <PageHeader title="Technical Debt" />
      <PageBody
        className={cx({ [css.withError]: !!error })}
        error={error ? getErrorMessage(error) : null}
        retryOnError={voidFn(refetch)}
        noData={{
          when: () => items?.length === 0 && searchTerm === undefined,
          image: noDataImage,
          message: 'No technical debt items found',
          button: NewItemButton
        }}>
        <LoadingSpinner visible={loading && !searchTerm} />

        <Container padding="xlarge">
          <Layout.Horizontal spacing="large" className={css.layout}>
            {NewItemButton}
            <FlexExpander />
            <SearchInputWithSpinner loading={loading} query={searchTerm} setQuery={setSearchTerm} />
          </Layout.Horizontal>

          <Container margin={{ top: 'medium' }}>
            {!!items?.length && (
              <Table<TechDebtItem>
                className={css.table}
                columns={columns}
                data={items || []}
                getRowClassName={row => cx(css.row, !row.original.description && css.noDesc)}
              />
            )}
            <NoResultCard
              showWhen={() => !!items && items?.length === 0 && !!searchTerm?.length}
              forSearch={true}
            />
          </Container>
          <ResourceListingPagination response={response} page={page} setPage={setPage} />
        </Container>
      </PageBody>
    </Container>
  )
}

export default TechDebtList
