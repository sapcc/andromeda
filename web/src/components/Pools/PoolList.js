import React, {useState} from "react"

import {Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Stack,} from "@cloudoperators/juno-ui-components"
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

            <Stack
                distribution="between"
                direction="horizontal"
                alignment="center"
                className="jn-px-6 jn-py-3 mt-6 jn-bg-theme-background-lvl-1">
                <div className="jn-text-lg jn-text-theme-high">
                    <strong>Pools</strong>
                </div>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={() => setModal("NewPoolsItem")}
                    label="Add a Pool"
                />
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
                        group.pools.map((pool, index) => (
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
