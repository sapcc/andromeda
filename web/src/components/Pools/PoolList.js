// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useState} from "react"

import {Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Stack, Select, SelectOption, SearchInput} from "@cloudoperators/juno-ui-components"
import PoolListItem from "./PoolListItem"
import {authStore, urlStore} from "../../store"
import {fetchAll, nextPageParam} from "../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import MemberList from "./Members/MemberList";
import MonitorList from "./Monitors/MonitorList";
import {Error, Loading} from "../Components";

const PoolList = () => {
    const auth = authStore((state) => state.auth)
    const setModal = urlStore((state) => state.openModal)
    const selectedPool = urlStore((state) => state.pool)
    const [domain, setDomain] = urlStore((state) => [state.domain, state.setDomain])
    const [error, setError] = useState()

    const {
        data,
        isLoading,
        isSuccess,
        fetchNextPage,
        hasNextPage,
        isFetching,
        isFetchingNextPage,
    } = useInfiniteQuery({
        queryKey: ["pools"],
        queryFn: fetchAll,
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: (err) => setError(err),
        onSuccess: () => setError(undefined),
    })

    return (
        <>
            {/* Error Bar */}
            <Error error={error} />

            {/* Loading indicator for page content */}
            <Loading isLoading={isLoading} />
            <Stack gap="2" className="jn-px-2 jn-py-2 jn-bg-theme-background-lvl-1">
                <div className="jn-text-lg jn-text-theme-high">
                    <strong>Pools</strong>
                </div>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={() => setModal("NewPoolsItem")}
                    label="Add a Pool"
                />
                <span>
                    <Select
                        placeholder="Filter for Domain..."
                        onChange={(e) => setDomain(e === "All" ? null : e)}
                        value={domain}
                    >
                        <SelectOption value="All" label="All" />
                        {/* Filter for domains */}
                        {isSuccess && [...new Set(data.pages.map((group, _) =>
                            group.pools.map((pool, _) =>
                                pool.domains
                            ).flat()).flat())].map((domain, i) => (
                                <SelectOption key={i} value={domain} />
                            )
                        )}
                    </Select>
                </span>
            </Stack>
            {isSuccess ? (
                <DataGrid columns={6}>
                <DataGridRow>
                        <DataGridHeadCell>ID/Name</DataGridHeadCell>
                        <DataGridHeadCell>#Domains/#Members/#Monitors</DataGridHeadCell>
                        <DataGridHeadCell>Created</DataGridHeadCell>
                        <DataGridHeadCell>Updated</DataGridHeadCell>
                        <DataGridHeadCell>Status</DataGridHeadCell>
                        <DataGridHeadCell className="jn-items-end">Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.pools.map((pool, index) => (domain == null || pool.domains.includes(domain)) && (
                            <PoolListItem
                                key={index}
                                pool={pool}
                                isActive={selectedPool === pool.id}
                                setError={setError}
                            />)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="jn-p-4">There are no pools to display.</div>
            )
            }
            <Box>
                <Button
                    variant="subdued"
                    size="small"
                    icon="expandMore"
                    onClick={() => fetchNextPage()}
                    disabled={!hasNextPage || isFetchingNextPage}
                    className="whitespace-nowrap"
                    progress={isFetching}
                    label={isFetchingNextPage
                        ? 'Loading more...'
                        : hasNextPage
                            ? 'Load More'
                            : 'Nothing more to load'}
                />
            </Box>

            {selectedPool && (
                <Stack direction="vertical" gap="2">
                    <MemberList />
                    <MonitorList />
                </Stack>
            )}
        </>
    )
}

export default PoolList
