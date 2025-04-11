import React from "react";

import {authStore} from "../../store";
import {useInfiniteQuery} from "@tanstack/react-query";
import {fetchAll, nextPageParam} from "../../actions";
import {Button, Icon} from "@cloudoperators/juno-ui-components";

const DomainMenu = ({formState, setFormState, setError}) => {
    const auth = authStore((state) => state.auth)

    const {
        isLoading,
        data,
        hasNextPage,
        fetchNextPage,
        isFetching
    } = useInfiniteQuery({
        queryKey: ["domains"],
        queryFn: fetchAll,
        getNextPageParam: nextPageParam,
        meta: auth,
        onError: setError
    })

    const toggleDomain = (id) => {
        let domains = [...formState.domains]
        if (domains.includes(id)) {
            domains = domains.filter(item => item !== id)
        } else {
            domains.push(id)
        }

        setFormState({
            ...formState, domains: domains,
        })
    }

    return (
        <div>
            <table className="table-auto w-full jn-text-left">
                <thead className="jn-bg-theme-background-lvl-2">
                    <tr>
                        <th>Name/ID</th>
                        <th>FQDN</th>
                        <th>Mode</th>
                        <th>Provider</th>
                    </tr>
                </thead>
                <tbody>
                {data?.pages.map((group, i) => group.domains.map((domain, index) => (
                    <tr onClick={(e) => {
                        e.preventDefault();
                        toggleDomain(domain.id)
                    }}
                        className={`cursor-pointer hover:jn-bg-theme-background-lvl-3 hover:jn-text-theme-accent ${formState.domains.includes(domain.id) && "jn-text-theme-accent"}`}
                    >
                        <td className="jn-inline-flex">
                            <Icon
                                icon={formState.domains.includes(domain.id) ? "checkCircle" : "addCircle"}
                                className="jn-mr-2"
                            />
                            {`${domain.name || domain.id}`}</td>
                        <td>{domain.fqdn}</td>
                        <td>{domain.mode}</td>
                        <td className={"place-self-end"}>{domain.provider}</td>
                    </tr>
                )))
                }
                < /tbody>
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
export default DomainMenu
