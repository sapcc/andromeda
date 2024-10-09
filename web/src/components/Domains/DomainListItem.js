import React, {useMemo, useState} from "react"

import {
    DataGridCell,
    DataGridRow,
    Icon,
    Stack,
    ContextMenu,
    MenuItem,
    Pill,
} from "@cloudoperators/juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {JsonModal, ListItemSpinner, ListItemStatus} from "../Components";

const DomainListItem = ({domain, setError}) => {
    const setPanel = urlStore((state) => state.openPanel)
    const [setTab, setDomain] = urlStore((state) => [state.setTab, state.setDomain])
    const auth = authStore((state) => state.auth)
    const [showJson, setShowJson] = useState(false)
    const queryClient = useQueryClient()
    const createdAt = useMemo(() => {
        if (domain.created_at) {
            return DateTime.fromISO(domain.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [domain.created_at])
    const updatedAt = useMemo(() => {
        if (domain.updated_at) {
            return DateTime.fromISO(domain.updated_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [domain.updated_at])

    const {mutate} = useMutation({mutationFn: deleteItem})

    const handleEditDomainClick = () => setPanel("Domain", domain.id)
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
                    const queryKey = ["domains"]
                    queryClient
                        .setQueryDefaults(queryKey, {refetchInterval: 5000})
                    queryClient.invalidateQueries({
                        queryKey: queryKey
                    })
                        .then()
                },
                onError: setError
            }
        )
    }

    return (
        <>
            <DataGridRow>
                <DataGridCell>
                    <ListItemSpinner data={domain} />
                </DataGridCell>
                <DataGridCell>
                    <Stack direction="vertical">
                        <p>{domain.fqdn}</p>
                        {domain.cname_target && <small><a target="_blank" href={"https://"+domain.cname_target}>{domain.cname_target}</a></small>}
                    </Stack>
                </DataGridCell>
                <DataGridCell>{domain.record_type}</DataGridCell>
                <DataGridCell>{domain.provider}</DataGridCell>
                <DataGridCell>{createdAt}</DataGridCell>
                <DataGridCell>{updatedAt}</DataGridCell>
                <DataGridCell><Pill
                    onClick={() => {if(domain.pools.length) {setDomain(domain.id);  setTab(1); }}}
                    pillKeyLabel="#"
                    pillValueLabel={domain.pools.length || "0"}
                /></DataGridCell>
                <DataGridCell><ListItemStatus data={domain} /></DataGridCell>
                <DataGridCell className="jn-items-end">
                    {/* Use <Stack> to align and space elements: */}
                    <Stack gap="1.5">
                        <Icon
                            icon="edit"
                            size="18"
                            className="leading-none self-center"
                            onClick={handleEditDomainClick}
                        />
                        <ContextMenu>
                            <MenuItem
                                icon="deleteForever"
                                label="Delete"
                                onClick={handleDeleteDomainClick}
                            />
                            <MenuItem
                                icon="info"
                                label="JSON"
                                onClick={() => setShowJson(!showJson)}
                            />
                        </ContextMenu>
                    </Stack>
                </DataGridCell>
            </DataGridRow>
            {showJson && JsonModal(domain)}
        </>
    )
}

export default DomainListItem
