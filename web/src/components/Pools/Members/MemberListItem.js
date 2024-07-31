import React, {useEffect} from "react"

import {DataGridCell, DataGridRow, Icon, Stack} from "@cloudoperators/juno-ui-components"
import {authStore, urlStore} from "../../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {deleteItem, fetchItem, updateAttributes} from "../../../actions"
import {ListItemSpinner, ListItemStatus} from "../../Components";

const MemberListItem = ({member, setError}) => {
    const auth = authStore((state) => state.auth)
    const openPanel = urlStore((state) => state.openPanel)
    const queryClient = useQueryClient()

    const queryDatacenter = useQuery({
        queryKey: ["datacenters", member.datacenter_id],
        queryFn: fetchItem,
        enabled: 'datacenter_id' in member,
        meta: auth,
    })
    const mutation = useMutation({mutationFn: deleteItem})

    // update formState when data is fetched
    useEffect(() => {
        if (queryDatacenter.error) {
            setError(queryDatacenter.error)
        }
    }, [queryDatacenter.error]);

    const handleEditMemberClick = () => openPanel("Member", member.id)
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
                    queryClient.invalidateQueries({
                        queryKey: queryKey
                    })
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
            <DataGridCell className="jn-items-end">
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
