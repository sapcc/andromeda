import React from "react"

import {DataGridCell, DataGridRow, Icon, Stack} from "juno-ui-components"
import {authStore, useStore} from "../../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {currentState, push} from "url-state-provider"
import {deleteItem, fetchItem} from "../../../actions"
import {ListItemSpinner, ListItemStatus} from "../../Components";

const MemberListItem = ({member, setError}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const queryClient = useQueryClient()

    const queryDatacenter = useQuery(
        ["datacenters", member.datacenter_id],
        fetchItem,
        {
            enabled: 'datacenter_id' in member,
            meta: auth,
            onError: setError,
        }
    )
    const mutation = useMutation(deleteItem)

    const handleEditMemberClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {
            ...urlState,
            currentPanel: "Member",
            id: member.id,
        })
    }

    const handleDeleteMemberClick = () => {
        mutation.mutate(
            {
                key: "members",
                id: member.id,
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: () => {
                    const queryKey= ["members", {pool_id: member.pool_id}]
                    queryClient
                        .setQueryDefaults(queryKey, {refetchInterval: 5000})
                    queryClient
                        .invalidateQueries(queryKey)
                        .then()
                },
                onError: setError,
            }
        )
    }

    return (
        <DataGridRow className={member.admin_state_up ? "" : "text-theme-background-lvl-4"}>
            <DataGridCell>
                <ListItemSpinner data={member} />
            </DataGridCell>
            <DataGridCell>{member.address}</DataGridCell>
            <DataGridCell>{member.port}</DataGridCell>
            <DataGridCell>{queryDatacenter.data?.datacenter.name}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={member} />
            </DataGridCell>
            <DataGridCell>
                {/* Use <Stack> to align and space elements: */}
                <Stack gap="1.5">
                    <Icon
                        icon="edit"
                        size="18"
                        className="leading-none"
                        onClick={handleEditMemberClick}
                    />
                    <Icon
                        icon="deleteForever"
                        size="18"
                        className="leading-none"
                        onClick={handleDeleteMemberClick}
                    />
                </Stack>
            </DataGridCell>
        </DataGridRow>
    )
}

export default MemberListItem
