import React, {useEffect, useState} from "react"

import {Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Message, Spinner, Stack,} from "juno-ui-components"
import PoolListItem from "./PoolListItem"
import {authStore, useStore} from "../../store"
import {currentState, push} from "url-state-provider"
import {fetchAll, nextPageParam} from "../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import MemberList from "./Members/MemberList";
import MonitorList from "./Monitors/MonitorList";
import {Error, Loading} from "../Components";

const PoolList = () => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
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
        ["pools"],
        fetchAll,
        {
            getNextPageParam: nextPageParam,
            meta: auth,
            onError: setError,
            onSuccess: () => setError(undefined),
        }
    )
    const handleNewPoolClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {...urlState, currentModal: "NewPoolsItem"})
    }

    const [selectedPool, setSelectedPool] = useState("");
    useEffect(() => {
        const urlState = currentState(urlStateKey)
        if (urlState?.pool) setSelectedPool(urlState.pool)
    }, [urlStateKey])

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
                    onClick={handleNewPoolClick}
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
                        <DataGridHeadCell>Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.pools.map((pool, index) => (
                            <PoolListItem
                                key={index}
                                pool={pool}
                                setSelectedPool={setSelectedPool}
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
                <Stack direction="horizontal" gap="2">
                    <MemberList poolID={selectedPool} setSelectedPool={setSelectedPool} />
                    <MonitorList poolID={selectedPool} setSelectedPool={setSelectedPool} />
                </Stack>
            )}
        </>
    )
}

export default PoolList
