import React, {useState} from "react"

import {Badge, Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Stack,} from "juno-ui-components"
import MemberListItem from "./MemberListItem"
import {authStore, urlStore} from "../../../store"
import {fetchAll, nextPageParam} from "../../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import {Error, Loading} from "../../Components";

const MemberList = () => {
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
        queryKey: ["members", {pool_id: poolId}],
        queryFn: fetchAll,
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: setError,
        onSuccess: () => setError(undefined),
    })

    return (
        <Stack direction="vertical" className="basis-1 md:basis-1/2 mt-6">
            <Stack gap="2" className="jn-px-2 jn-py-2 jn-bg-theme-background-lvl-1">
                <div className="jn-text-lg jn-text-theme-high">
                    <strong>Associated Members<Badge>Pool {poolId}</Badge></strong>
                </div>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={() => openModal("NewMembersItem")}
                    label="Add a Member"
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


            {isSuccess && data.pages[0]?.members.length ? (
                <DataGrid columns={6}>
                    <DataGridRow>
                        <DataGridHeadCell>ID/Name</DataGridHeadCell>
                        <DataGridHeadCell>Address</DataGridHeadCell>
                        <DataGridHeadCell>Port</DataGridHeadCell>
                        <DataGridHeadCell>Datacenter</DataGridHeadCell>
                        <DataGridHeadCell>Status</DataGridHeadCell>
                        <DataGridHeadCell className="jn-items-end">Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.members.map((member, index) => (
                            <MemberListItem key={index} member={member} setError={setError}/>)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="jn-p-4">There are no members to display.</div>
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

export default MemberList
