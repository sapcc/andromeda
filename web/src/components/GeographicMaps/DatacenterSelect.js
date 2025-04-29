// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React from "react";
import {useInfiniteQuery} from "@tanstack/react-query";
import {fetchAll, nextPageParam} from "../../actions";
import {authStore} from "../../store";
import {Button, SelectOption, Select} from "@cloudoperators/juno-ui-components";


const DatacenterSelect = ({setDatacenter, setError, label = "Select Datacenter", className}) => {
    const auth = authStore((state) => state.auth)
    const {
        isLoading,
        data,
        hasNextPage,
        fetchNextPage,
        isFetching
    } = useInfiniteQuery({
        queryKey: ["datacenters"],
        queryFn: fetchAll,
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: setError
    })

    return (
        <Select
            onValueChange={(value) => setDatacenter(value)}
            label={label}
            position="align-items"
            className={className}
            required
        >
            {/* Render items: */}
            {data?.pages.map((group, i) => group.datacenters.map((datacenter, index) => (
                <SelectOption
                    value={datacenter.id}
                    label={`${datacenter.name || datacenter.id}`}
                />
            )))}
            {hasNextPage && (
                <Button
                    label={isLoading ? "Loading..." :
                        isFetching ? 'Loading more...'
                            : hasNextPage
                                ? 'Load More'
                                : 'Nothing more to load'}
                    onClick={hasNextPage ? () => fetchNextPage() : undefined}
                    icon={hasNextPage ? "expandMore" : "info"}
                />
            )}
        </Select>
    )
}

export default DatacenterSelect