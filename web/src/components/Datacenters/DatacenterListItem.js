import React, {useMemo} from "react"
import {DataGridCell, DataGridRow, Icon, Stack} from "@cloudoperators/juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {ListItemSpinner, ListItemStatus} from "../Components";

const DatacenterListItem = ({datacenter, setError}) => {
    const auth = authStore((state) => state.auth)
    const setPanel = urlStore((state) => state.openPanel)
    const queryClient = useQueryClient()
    const createdAt = useMemo(() => {
        if (datacenter.created_at) {
            return DateTime.fromISO(datacenter.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [datacenter.created_at])
    const updatedAt = useMemo(() => {
        if (datacenter.updated_at) {
            return DateTime.fromISO(datacenter.updated_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [datacenter.updated_at])
    const mutation = useMutation({mutationFn: deleteItem})

    const handleEditDatacenterClick = () => setPanel("Datacenter", datacenter.id)
    const handleDeleteDatacenterClick = () => {
        mutation.mutate(
            {
                key: "datacenters",
                endpoint: auth?.endpoint,
                id: datacenter.id,
                token: auth?.token,
            },
            {
                onSuccess: () => {
                    const queryKey = ["datacenters"]
                    queryClient
                        .setQueryDefaults(queryKey, {refetchInterval: 5000})
                    // refetch datacenters
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
        <DataGridRow>
            <DataGridCell>
                <ListItemSpinner data={datacenter} />
            </DataGridCell>
            <DataGridCell>{datacenter.continent}</DataGridCell>
            <DataGridCell>{datacenter.country}</DataGridCell>
            <DataGridCell>{datacenter.state_or_province}</DataGridCell>
            <DataGridCell>{datacenter.city}</DataGridCell>
            <DataGridCell>{datacenter.latitude}, {datacenter.longitude}</DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={datacenter} />
            </DataGridCell>
            <DataGridCell>
                {/* Use <Stack> to align and space elements: */}
                <Stack gap="1.5">
                    {datacenter.project_id === auth?.project?.id && (
                        <>
                            <Icon
                                icon="edit"
                                size="18"
                                className="leading-none"
                                onClick={handleEditDatacenterClick}
                            />
                            <Icon
                                icon="deleteForever"
                                size="18"
                                className="leading-none"
                                onClick={handleDeleteDatacenterClick}
                            />
                        </>
                    )}
                    <Icon
                        icon="openInNew"
                        size="18"
                        href={`http://www.google.com/maps/place/${datacenter.latitude},${datacenter.longitude}`}
                        target="_blank"
                        className="leading-none"
                    />
                </Stack>
            </DataGridCell>
        </DataGridRow>
    )
}

export default DatacenterListItem
