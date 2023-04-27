import React, {useMemo} from "react"

import {DataGridCell, DataGridRow, Icon, Spinner, Stack} from "juno-ui-components"
import {authStore, useStore} from "../../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {currentState, push} from "url-state-provider"
import {deleteItem} from "../../../actions"
import {DateTime} from "luxon";
import {ListItemSpinner, ListItemStatus} from "../../Components";

const MonitorListItem = ({monitor, setError}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const queryClient = useQueryClient()
    const createdAt = useMemo(() => {
        if (monitor.created_at) {
            return DateTime.fromISO(monitor.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [monitor.created_at])
    const updatedAt = useMemo(() => {
        if (monitor.updated_at) {
            return DateTime.fromISO(monitor.updated_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [monitor.updated_at])

    const {mutate} = useMutation(deleteItem)

    const handleEditMonitorClick = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {
            ...urlState,
            currentPanel: "Monitor",
            id: monitor.id,
        })
        queryClient
            .setQueryDefaults(["monitors", {pool_id: monitor.pool_id}], {refetchInterval: 5000})
        queryClient
            .setQueryDefaults(["monitors"], {refetchInterval: 5000})
        queryClient
            .setQueryDefaults(["domains"], {refetchInterval: 5000})
    }

    const handleDeleteMonitorClick = () => {
        mutate(
            {
                key: "monitors",
                id: monitor.id,
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: () => {
                    const queryKey = ["monitors", {pool_id: monitor.pool_id}]
                    queryClient
                        .setQueryDefaults(queryKey, {refetchInterval: 5000})
                    queryClient
                        .invalidateQueries(queryKey)
                        .then()
                },
                onError: (err) => {
                    setError(err)
                }
            }
        )
    }

    return (
        <DataGridRow>
            <DataGridCell>
                <ListItemSpinner data={monitor} maxLength="6" />
            </DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={monitor} />
            </DataGridCell>
            <DataGridCell>
                {/* Use <Stack> to align and space elements: */}
                <Stack gap="1.5">
                    <Icon
                        icon="edit"
                        size="18"
                        className="leading-none"
                        onClick={handleEditMonitorClick}
                    />
                    <Icon
                        icon="deleteForever"
                        size="18"
                        className="leading-none"
                        onClick={handleDeleteMonitorClick}
                    />
                </Stack>
            </DataGridCell>
        </DataGridRow>
    )
}

export default MonitorListItem
