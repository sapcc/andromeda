import React, {useMemo, useState} from "react"

import {DataGridCell, DataGridRow, Icon, Stack, Tooltip, TooltipContent, TooltipTrigger} from "juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {JsonModal, ListItemSpinner, ListItemStatus} from "../Components";
import {ContextMenu} from "juno-ui-components/build/ContextMenu";
import {MenuItem} from "juno-ui-components/build/MenuItem";

const DomainListItem = ({domain, setError}) => {
    const setPanel = urlStore((state) => state.openPanel)
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
                    <Stack gap="1.5" alignment="center">
                        {domain.fqdn}
                        {domain.cname_target && (
                            <Tooltip triggerEvent="hover" variant="info"
                            >
                                <TooltipTrigger asChild>
                                    <Icon icon="help" size="18" />
                                </TooltipTrigger>
                                <TooltipContent>Use CNAME record: {domain.cname_target}</TooltipContent>
                            </Tooltip>
                        )}
                    </Stack>
                </DataGridCell>
                <DataGridCell>{domain.record_type}</DataGridCell>
                <DataGridCell>{domain.provider}</DataGridCell>
                <DataGridCell>{createdAt}</DataGridCell>
                <DataGridCell>{updatedAt}</DataGridCell>
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
                        <Icon
                            icon="openInNew"
                            size="18"
                            href={domain.fqdn}
                            target="_blank"
                            className="leading-none self-center"
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
