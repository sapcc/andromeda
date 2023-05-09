import React, {useMemo} from "react"
import {DataGridCell, DataGridRow, Icon, Stack} from "juno-ui-components"
import {authStore, useStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {currentState, push} from "url-state-provider"
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {ListItemSpinner, ListItemStatus} from "../Components";

const GeographicMapListItem = ({geomap, setError}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
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
    const mutation = useMutation(deleteItem)

    const handleEditGeographicMapClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {
            ...urlState,
            currentPanel: "GeographicMap",
            id: geomap.id,
        })
    }

    const handleShowGeographicMapClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {
            ...urlState,
            currentModal: "ShowGeographicMap",
            id: geomap.id,
        })
    }

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
                <ListItemSpinner data={geomap} />
            </DataGridCell>
            <DataGridCell>{geomap.scope}</DataGridCell>
            <DataGridCell>{geomap.provider}</DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={geomap} />
            </DataGridCell>
            <DataGridCell>
                <Stack gap="1.5">
                    {geomap.project_id === auth?.project?.id && (
                        <>
                            <Icon
                                icon="openInNew"
                                size="18"
                                className="leading-none"
                                onClick={handleShowGeographicMapClick}
                            />
                            {/*<Icon
                                icon="edit"
                                size="18"
                                className="leading-none"
                                onClick={handleEditGeographicMapClick}
                            />*/}
                            <Icon
                                icon="deleteForever"
                                size="18"
                                className="leading-none"
                                onClick={handleDeleteGeographicMapClick}
                            />
                        </>
                    )}
                </Stack>
            </DataGridCell>
        </DataGridRow>
    )
}

export default GeographicMapListItem
