import React, {useState} from "react"

import {Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Message, Spinner, Stack,} from "juno-ui-components"
import GeographicMapListItem from "./GeographicMapListItem"
import {authStore, useStore} from "../../store"
import {currentState, push} from "url-state-provider"
import {fetchAll, nextPageParam} from "../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import {Error, Loading} from "../Components";

const GeographicMapList = () => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const [error, setError] = useState()
    const {
        data,
        isSuccess,
        fetchNextPage,
        hasNextPage,
        isLoading,
        isFetching,
        isFetchingNextPage,
    } = useInfiniteQuery(
        ["geomaps"],
        fetchAll,
        {
            getNextPageParam: nextPageParam,
            meta: auth,
            onError: setError,
            onSuccess: () => setError(undefined),
        }
    )

    const handleNewGeographicMapClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {...urlState, currentModal: "NewGeographicMapsItem"})
    }

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
                    <strong>Geographic Maps</strong>
                </div>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={handleNewGeographicMapClick}
                    label="Add a Geographic Map"
                />
            </Stack>
            {isSuccess ? (
                <DataGrid columns={7}>
                    <DataGridRow>
                        <DataGridHeadCell>ID/Name</DataGridHeadCell>
                        <DataGridHeadCell>Scope</DataGridHeadCell>
                        <DataGridHeadCell>Provider</DataGridHeadCell>
                        <DataGridHeadCell>Created</DataGridHeadCell>
                        <DataGridHeadCell>Updated</DataGridHeadCell>
                        <DataGridHeadCell>Status</DataGridHeadCell>
                        <DataGridHeadCell className="jn-items-end">Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.geomaps.map((geomap, index) => (
                            <GeographicMapListItem key={index} geomap={geomap} setError={setError}/>)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="jn-p-4">There are no geographical maps to display.</div>
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
                >
                    {isFetching ? <Spinner variant="primary"/> : null}
                    {isFetchingNextPage
                        ? 'Loading more...'
                        : hasNextPage
                            ? 'Load More'
                            : 'Nothing more to load'}
                </Button>
            </Box>
        </>
    )
}

export default GeographicMapList
