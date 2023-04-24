import React, {useState} from "react"

import {
    Box,
    Button,
    ContentAreaToolbar,
    DataGrid,
    DataGridHeadCell,
    DataGridRow,
    Message,
    Spinner,
    Stack,
} from "juno-ui-components"
import DomainListItem from "./DomainListItem"
import {authStore, useStore} from "../../store"
import {currentState, push} from "url-state-provider"
import {fetchAll, nextPageParam} from "../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import {Error, Loading} from "../Components";

const DomainList = ({domains}) => {
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
        ["domains"],
        fetchAll,
        {
            getNextPageParam: nextPageParam,
            meta: auth,
            onError: setError,
            onSuccess: () => setError(undefined),
        },
    )

    const handleNewDomainClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {...urlState, currentModal: "NewDomainsItem"})
    }

    return (
        <>
            {/* Error Bar */}
            <Error error={error} />

            {/* Loading indicator for page content */}
            <Loading isLoading={isLoading} />

            <ContentAreaToolbar>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={handleNewDomainClick}
                    label="Add a Domain"
                />
            </ContentAreaToolbar>
            {isSuccess ? (
                <DataGrid columns={8}>
                    <DataGridRow>
                        <DataGridHeadCell>ID/Name</DataGridHeadCell>
                        <DataGridHeadCell>FQDN</DataGridHeadCell>
                        <DataGridHeadCell>Record Type</DataGridHeadCell>
                        <DataGridHeadCell>Provider</DataGridHeadCell>
                        <DataGridHeadCell>Created</DataGridHeadCell>
                        <DataGridHeadCell>Updated</DataGridHeadCell>
                        <DataGridHeadCell>Status</DataGridHeadCell>
                        <DataGridHeadCell>Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.domains.map((domain, index) => (
                            <DomainListItem key={index} domain={domain} setError={setError}/>)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="jn-p-4">There are no domains to display.</div>
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

export default DomainList
