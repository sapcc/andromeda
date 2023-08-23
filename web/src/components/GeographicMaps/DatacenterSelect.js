import React from "react";
import {useInfiniteQuery} from "@tanstack/react-query";
import {fetchAll, nextPageParam} from "../../actions";
import {authStore} from "../../store";
import {Button, SelectOption, Select} from "juno-ui-components";


const DatacenterSelect = ({setDatacenter, setError, label = "Select Datacenter", className}) => {
    const auth = authStore((state) => state.auth)
    const {
        isLoading,
        data,
        hasNextPage,
        fetchNextPage,
        isFetching
    } = useInfiniteQuery(["datacenters"], fetchAll, {
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: setError
    })

    return (
        <Select
            onValueChange={(value) => setDatacenter(value)}
            label={label}
            required="true"
            position="align-items"
            className={className}
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