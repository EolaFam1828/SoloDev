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
  useToaster
} from '@harnessio/uicore'
import { Color } from '@harnessio/design-system'
import cx from 'classnames'
import type { CellProps, Column } from 'react-table'
import Keywords from 'react-keywords'
import { useGet } from 'restful-react'
import { CheckCircle, DeleteCircle } from 'iconoir-react'
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
import css from './MonitorList.module.scss'

interface Monitor {
  id: string
  name: string
  description?: string
  url: string
  status: 'up' | 'down' | 'degraded'
  uptime: number
  lastChecked: number
  responseTime: number
}

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'up':
      return <CheckCircle width={20} height={20} color="#16B671" />
    case 'down':
      return <DeleteCircle width={20} height={20} color="#E74C3C" />
    case 'degraded':
      return <DeleteCircle width={20} height={20} color="#F39C12" />
    default:
      return null
  }
}

const getStatusColor = (status: string): string => {
  switch (status) {
    case 'up':
      return Color.GREEN_600
    case 'down':
      return Color.RED_600
    case 'degraded':
      return Color.YELLOW_600
    default:
      return Color.GREY_600
  }
}

const MonitorList = () => {
  const space = useGetSpaceParam()
  const { getString } = useStrings()
  const { showSuccess } = useToaster()
  const [searchTerm, setSearchTerm] = useState<string | undefined>()
  const pageBrowser = useQueryParams<PageBrowserProps>()
  const pageInit = pageBrowser.page ? parseInt(pageBrowser.page) : 1
  const [page, setPage] = usePageIndex(pageInit)

  const {
    data: monitors,
    error,
    loading,
    refetch,
    response
  } = useGet<Monitor[]>({
    path: `/api/v1/spaces/${space}/+/monitors`,
    queryParams: { page, limit: LIST_FETCHING_LIMIT, query: searchTerm }
  })

  const NewMonitorButton = (
    <Button
      text="New Monitor"
      variation={ButtonVariation.PRIMARY}
      icon="plus"
      onClick={() => {
        showSuccess('Create monitor functionality to be implemented')
      }}
    />
  )

  const columns: Column<Monitor>[] = useMemo(
    () => [
      {
        Header: getString('name') || 'Monitor Name',
        width: 'calc(100% - 380px)',
        Cell: ({ row }: CellProps<Monitor>) => {
          const record = row.original
          return (
            <Container className={css.nameContainer}>
              <Layout.Horizontal spacing="small" style={{ flexGrow: 1 }}>
                <Layout.Vertical flex className={css.name}>
                  <Text className={css.monitorName} lineClamp={1}>
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
        width: '100px',
        Cell: ({ row }: CellProps<Monitor>) => {
          return (
            <Layout.Horizontal spacing="small" style={{ alignItems: 'center' }}>
              {getStatusIcon(row.original.status)}
              <Text color={getStatusColor(row.original.status)} lineClamp={1}>
                {row.original.status.toUpperCase()}
              </Text>
            </Layout.Horizontal>
          )
        },
        disableSortBy: true
      },
      {
        Header: 'Uptime',
        width: '100px',
        Cell: ({ row }: CellProps<Monitor>) => {
          return (
            <Text color={Color.BLACK} lineClamp={1}>
              {row.original.uptime.toFixed(2)}%
            </Text>
          )
        },
        disableSortBy: true
      },
      {
        Header: 'Response Time',
        width: '130px',
        Cell: ({ row }: CellProps<Monitor>) => {
          return (
            <Text color={Color.BLACK} lineClamp={1}>
              {row.original.responseTime}ms
            </Text>
          )
        },
        disableSortBy: true
      },
      {
        Header: getString('updatedDate') || 'Last Checked',
        width: '140px',
        Cell: ({ row }: CellProps<Monitor>) => {
          return (
            <Layout.Horizontal style={{ alignItems: 'center' }}>
              <Text color={Color.BLACK} lineClamp={1} width={120}>
                {formatDate(row.original.lastChecked)}
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
      <PageHeader title="Uptime Monitor" />
      <PageBody
        className={cx({ [css.withError]: !!error })}
        error={error ? getErrorMessage(error) : null}
        retryOnError={voidFn(refetch)}
        noData={{
          when: () => monitors?.length === 0 && searchTerm === undefined,
          image: noDataImage,
          message: 'No monitors configured',
          button: NewMonitorButton
        }}>
        <LoadingSpinner visible={loading && !searchTerm} />

        <Container padding="xlarge">
          <Layout.Horizontal spacing="large" className={css.layout}>
            {NewMonitorButton}
            <FlexExpander />
            <SearchInputWithSpinner loading={loading} query={searchTerm} setQuery={setSearchTerm} />
          </Layout.Horizontal>

          <Container margin={{ top: 'medium' }}>
            {!!monitors?.length && (
              <Table<Monitor>
                className={css.table}
                columns={columns}
                data={monitors || []}
                getRowClassName={row => cx(css.row, !row.original.description && css.noDesc)}
              />
            )}
            <NoResultCard
              showWhen={() => !!monitors && monitors?.length === 0 && !!searchTerm?.length}
              forSearch={true}
            />
          </Container>
          <ResourceListingPagination response={response} page={page} setPage={setPage} />
        </Container>
      </PageBody>
    </Container>
  )
}

export default MonitorList
