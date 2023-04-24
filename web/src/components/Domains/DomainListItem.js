import React, {useMemo} from "react"

import {DataGridCell, DataGridRow, Icon, Spinner, Stack} from "juno-ui-components"
import {authStore, useStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {currentState, push} from "url-state-provider"
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {ListItemSpinner, ListItemStatus} from "../Components";

const DomainListItem = ({domain, setError}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const queryClient = useQueryClient()
    const createdAt = useMemo(() => {
        if (domain.created_at) {
            return DateTime.fromSQL(domain.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [domain.created_at])
    const updatedAt = useMemo(() => {
        if (domain.updated_at) {
            return DateTime.fromSQL(domain.updated_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [domain.updated_at])

    const {mutate} = useMutation(deleteItem)

    const handleEditDomainClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {
            ...urlState,
            currentPanel: "Domain",
            id: domain.id,
        })
    }

    const handleDeleteDomainClick = () => {
        mutate(
            {
                key: "domains",
                endpoint: auth?.endpoint,
                id: domain.id,
                token: auth?.token
            },
            {
                onSuccess: () => {
                    const queryKey = ["domains", endpoint]
                    queryClient
                        .setQueryDefaults(queryKey, {refetchInterval: 5000})
                    queryClient
                        .invalidateQueries(queryKey)
                        .then()
                },
                onError: setError
            }
        )
    }

    return (
        <DataGridRow>
            <DataGridCell>
                <ListItemSpinner data={domain} />
            </DataGridCell>
            <DataGridCell>{domain.fqdn}</DataGridCell>
            <DataGridCell>{domain.record_type}</DataGridCell>
            <DataGridCell>{domain.provider}</DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell><ListItemStatus data={domain} /></DataGridCell>
            <DataGridCell>
                {/* Use <Stack> to align and space elements: */}
                <Stack gap="1.5">
                    <Icon
                        icon="edit"
                        size="18"
                        className="leading-none"
                        onClick={handleEditDomainClick}
                    />
                    <Icon
                        icon="deleteForever"
                        size="18"
                        className="leading-none"
                        onClick={handleDeleteDomainClick}
                    />
                    <Icon
                        icon="openInNew"
                        size="18"
                        href={domain.fqdn}
                        target="_blank"
                        className="leading-none"
                    />
                </Stack>
            </DataGridCell>
        </DataGridRow>
    )
}

export default DomainListItem
