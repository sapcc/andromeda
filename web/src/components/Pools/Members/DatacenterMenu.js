import React from "react";

import {useInfiniteQuery} from "@tanstack/react-query";
import {fetchAll, nextPageParam} from "../../../actions";
import {authStore} from "../../../store";
import {Button, Icon} from "@cloudoperators/juno-ui-components";

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
        <div>
            <table className="table-auto w-full jn-text-left">
                <thead className="jn-bg-theme-background-lvl-2">
                    <tr className="">
                        <th>Name/ID</th>
                        <th>Country</th>
                        <th>City</th>
                        <th>Provider</th>
                    </tr>
                </thead>
                <tbody>
                {data?.pages.map((group, i) => group.datacenters.map((datacenter, index) => (
                    <tr
                        onClick={(e) => {
                            e.preventDefault();
                            onDatacenterClick(datacenter.id)
                        }}
                        className={`cursor-pointer hover:jn-bg-theme-background-lvl-3 hover:jn-text-theme-accent ${formState.datacenter_id === datacenter.id && "jn-text-theme-accent"}`}>
                        <td className={"jn-inline-flex"}>
                            <Icon
                                icon={formState.datacenter_id === datacenter.id ? "checkCircle" : "addCircle"}
                                className={"jn-mr-2"}
                            />
                            {`${datacenter.name || datacenter.id}`}</td>
                        <td>{datacenter.country}</td>
                        <td>{datacenter.city}</td>
                        <td className={"place-self-end"}>{datacenter.provider}</td>
                    </tr>
                )))}
                </tbody>
            </table>
            {hasNextPage && (
                <Button
                    className={"w-full"}
                    label={isLoading ? "Loading..." :
                        isFetching ? 'Loading more...'
                            : hasNextPage
                                ? 'Load More'
                                : 'Nothing more to load'}
                    onClick={hasNextPage ? () => fetchNextPage() : undefined}
                    icon={hasNextPage ? "expandMore" : "info"}
                />
            )}
        </div>
    )
}

export default DatacenterMenu