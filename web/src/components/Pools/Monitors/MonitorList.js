import React, {useState} from "react"

import {Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Stack,} from "juno-ui-components"
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
    } = useInfiniteQuery(
        ["monitors", {pool_id: poolId}],
        fetchAll,
        {
            getNextPageParam: nextPageParam,
            meta: auth,
            onError: setError,
            onSuccess: () => setError(undefined),
        }
    )

    return (
        <Stack direction="vertical" className="basis-1 md:basis-1/2">
            <Stack
                className="jn-px-6 jn-py-3 mt-6 jn-bg-theme-background-lvl-1">
                <div className="jn-text-lg jn-text-theme-high">
                    <strong>Associated Monitors</strong>
                </div>
                <Stack direction="horizontal" alignment="end" gap="2">
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
                        <DataGridHeadCell>Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.monitors.map((monitor, index) => (
                            <MonitorListItem key={index} monitor={monitor} setError={setError}/>)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="jn-p-4">There are no Monitors to display.</div>
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
