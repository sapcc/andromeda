// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useState} from "react"

import {Badge, Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Stack,} from "@cloudoperators/juno-ui-components"
import MonitorListItem from "./MonitorListItem"
import {authStore, urlStore} from "../../../store"
import {fetchAll, nextPageParam} from "../../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import {Error, Loading} from "../../Components";

const MonitorList = () => {
    const auth = authStore((state) => state.auth)
    const openModal = urlStore((state) => state.openModal)
    const [poolId, clearSelectedPool] = urlStore((state) => [state.pool, state.clearPool])
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
        queryKey: ["monitors", {pool_id: poolId}],
        queryFn: fetchAll,
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: setError,
        onSuccess: () => setError(undefined),
    })

    return (
        <Stack direction="vertical" className="basis-1 md:basis-1/2 mt-6">
            <Stack gap="2" className="px-2 py-2 bg-theme-background-lvl-1">
                <div className="text-lg text-theme-high">
                    <strong>Associated Monitors<Badge>Pool {poolId}</Badge></strong>
                </div>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={() => openModal("NewMonitorsItem")}
                    label="Add a Monitor"
                />
                <Button
                    icon="close"
                    onClick={clearSelectedPool}
                />
            </Stack>

            {/* Error Bar */}
            <Error error={error} />

            {/* Loading indicator for page content */}
            <Loading isLoading={isLoading} />

            {isSuccess && data.pages[0]?.monitors.length ? (
                <DataGrid columns={5}>
                    <DataGridRow>
                        <DataGridHeadCell>ID/Name</DataGridHeadCell>
                        <DataGridHeadCell>Created</DataGridHeadCell>
                        <DataGridHeadCell>Updated</DataGridHeadCell>
                        <DataGridHeadCell>Status</DataGridHeadCell>
                        <DataGridHeadCell className="items-end">Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.monitors.map((monitor, index) => (
                            <MonitorListItem key={index} monitor={monitor} setError={setError}/>)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="p-4">There are no Monitors to display.</div>
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
        </Stack>
    )
}

export default MonitorList
