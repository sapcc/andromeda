// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useState} from "react"

import {Box, Button, DataGrid, DataGridHeadCell, DataGridRow, Spinner, Stack,} from "@cloudoperators/juno-ui-components"
import DatacenterListItem from "./DatacenterListItem"
import {authStore, urlStore} from "../../store"
import {fetchAll, nextPageParam} from "../../actions";
import {useInfiniteQuery} from '@tanstack/react-query';
import {Error, Loading} from "../Components";

const DatacenterList = () => {
    const setModal = urlStore((state) => state.openModal)
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
    } = useInfiniteQuery({
        queryKey: ["datacenters"],
        queryFn: fetchAll,
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: setError,
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
                    <strong>Datacenters</strong>
                </div>
                <Button
                    variant="primary"
                    icon="addCircle"
                    onClick={() => setModal("NewDatacentersItem")}
                    label="Add a Datacenter"
                />
            </Stack>
            {isSuccess ? (
                <DataGrid columns={10}>
                    <DataGridRow>
                        <DataGridHeadCell>ID/Name</DataGridHeadCell>
                        <DataGridHeadCell>Continent</DataGridHeadCell>
                        <DataGridHeadCell>Country</DataGridHeadCell>
                        <DataGridHeadCell>State/Province</DataGridHeadCell>
                        <DataGridHeadCell>City</DataGridHeadCell>
                        <DataGridHeadCell>Location</DataGridHeadCell>
                        <DataGridHeadCell>Created</DataGridHeadCell>
                        <DataGridHeadCell>Updated</DataGridHeadCell>
                        <DataGridHeadCell>Status</DataGridHeadCell>
                        <DataGridHeadCell>Options</DataGridHeadCell>
                    </DataGridRow>

                    {/* Render items: */}
                    {data.pages.map((group, i) =>
                        group.datacenters.map((datacenter, index) => (
                            <DatacenterListItem key={index} datacenter={datacenter} setError={setError}/>)
                        )
                    )}
                </DataGrid>
            ) : (
                <div className="jn-p-4">There are no datacenters to display.</div>
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

export default DatacenterList
