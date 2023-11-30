import React, {useMemo, useState} from "react"
import {DataGridCell, DataGridRow, Icon, Stack} from "juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {JsonModal, ListItemSpinner, ListItemStatus} from "../Components";
import {ContextMenu} from "juno-ui-components/build/ContextMenu";
import {MenuItem} from "juno-ui-components/build/MenuItem";

const GeographicMapListItem = ({geomap, setError}) => {
    const auth = authStore((state) => state.auth)
    const [setPanel, setModalId] = urlStore((state) => [state.openPanel, state.openModalWithId])
    const [showJson, setShowJson] = useState(false)
    const queryClient = useQueryClient()
    const createdAt = useMemo(() => {
        if (geomap.created_at) {
            return DateTime.fromISO(geomap.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [geomap.created_at])
    const updatedAt = useMemo(() => {
        if (geomap.updated_at) {
            return DateTime.fromISO(geomap.updated_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [geomap.updated_at])
    const mutation = useMutation({mutationFn: deleteItem})

    const handleEditGeographicMapClick = () => setPanel("GeographicMap", geomap.id)
    const handleShowGeographicMapClick = () => setModalId("ShowGeographicMap", geomap.id)
    const handleDeleteGeographicMapClick = () => {
        mutation.mutate(
            {
                key: "geomaps",
                endpoint: auth?.endpoint,
                id: geomap.id,
                token: auth?.token,
            },
            {
                onSuccess: () => {
                    const queryKey = ["geomaps"]
                    queryClient
                        .setQueryDefaults(queryKey, {refetchInterval: 5000})
                    // refetch geomaps
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
                <ListItemSpinner data={geomap} />
            </DataGridCell>
            <DataGridCell>{geomap.assignments.length}</DataGridCell>
            <DataGridCell>{geomap.scope}</DataGridCell>
            <DataGridCell>{geomap.provider}</DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={geomap} />
            </DataGridCell>
            <DataGridCell className="jn-items-end">
                <Stack gap="1.5">
                    {geomap.project_id === auth?.project?.id && (
                        <>
                            <Icon
                                icon="openInNew"
                                size="18"
                                className="leading-none self-center"
                                onClick={handleShowGeographicMapClick}
                            />
                            {/*<Icon
                                icon="edit"
                                size="18"
                                className="leading-none self-center"
                                onClick={handleEditGeographicMapClick}
                            />*/}
                            <ContextMenu>
                                <MenuItem
                                    icon="deleteForever"
                                    label="Delete"
                                    onClick={handleDeleteGeographicMapClick}
                                />
                                <MenuItem
                                    icon="info"
                                    label="JSON"
                                    onClick={() => setShowJson(!showJson)}
                                />
                            </ContextMenu>
                        </>
                    )}
                </Stack>
            </DataGridCell>
        </DataGridRow>
        {showJson && JsonModal(geomap)}
        </>
    )
}

export default GeographicMapListItem
