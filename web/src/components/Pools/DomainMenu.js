import React from "react";

import {authStore, useStore} from "../../store";
import {useInfiniteQuery} from "@tanstack/react-query";
import {fetchAll, nextPageParam} from "../../actions";
import {Menu} from "juno-ui-components/build/Menu";
import {MenuItem} from "juno-ui-components/build/MenuItem";


const DomainMenu = ({formState, setFormState, setError}) => {
    const auth = authStore((state) => state.auth)

    const {
        isLoading,
        data,
        hasNextPage,
        fetchNextPage,
        isFetching
    } = useInfiniteQuery(["domains"], fetchAll, {
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
        <Menu variant="small">
            {/* Render items: */}
            {data?.pages.map((group, i) => group.domains.map((domain, index) => (
                <MenuItem
                    key={domain.id}
                    icon={formState.domains.includes(domain.id) ? "checkCircle" : "addCircle"}
                    label={`${domain.name || domain.id} (${domain.fqdn})`}
                    onClick={() => toggleDomain(domain.id)}
                    className={formState.domains.includes(domain.id) ? "jn-text-theme-info" : ""}
                />
            )))}
            <MenuItem
                label={isLoading ? "Loading..." :
                    isFetching ? 'Loading more...'
                        : hasNextPage
                            ? 'Load More'
                            : 'Nothing more to load'}
                onClick={hasNextPage ? () => fetchNextPage() : null}
                icon={hasNextPage ? "expandMore" : "info"}
            />
        </Menu>
    )
}
export default DomainMenu
