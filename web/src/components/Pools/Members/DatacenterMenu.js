import React from "react";

import {useInfiniteQuery} from "@tanstack/react-query";
import {fetchAll, nextPageParam} from "../../../actions";
import {authStore} from "../../../store";
import {Menu, MenuItem} from "@cloudoperators/juno-ui-components";

const DatacenterMenu = ({formState, setFormState, setError}) => {
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

    const onDatacenterClick = (id) => {

        setFormState({
            ...formState,
            datacenter_id: id === formState.datacenter_id ? undefined : id,
        })
    }

    return (
        <Menu variant="small">
            {/* Render items: */}
            {data?.pages.map((group, i) => group.datacenters.map((datacenter, index) => (
                <MenuItem
                    key={datacenter.id}
                    icon={formState.datacenter_id === datacenter.id ? "checkCircle" : "addCircle"}
                    label={`${datacenter.name || datacenter.id}`}
                    onClick={(e) => {e.preventDefault(); onDatacenterClick(datacenter.id)}}
                    className={formState.datacenter_id === datacenter.id ? "jn-text-theme-info" : ""}
                />
            )))}
            {hasNextPage && (
                <MenuItem
                    label={isLoading ? "Loading..." :
                        isFetching ? 'Loading more...'
                            : hasNextPage
                                ? 'Load More'
                                : 'Nothing more to load'}
                    onClick={hasNextPage ? () => fetchNextPage() : undefined}
                    icon={hasNextPage ? "expandMore" : "info"}
                />
            )}
        </Menu>
    )
}

export default DatacenterMenu