import React, {useMemo} from "react"

import {Badge, DataGridCell, DataGridRow, Icon, Spinner, Stack} from "juno-ui-components"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {currentState, push} from "url-state-provider"
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {authStore, useStore} from "../../store"
import {ListItemSpinner, ListItemStatus} from "../Components";

const PoolListItem = ({pool, setSelectedPool, isActive, setError}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const urlState = currentState(urlStateKey)
    const queryClient = useQueryClient()
    const createdAt = useMemo(() => {
        if (pool.created_at) {
            return DateTime.fromISO(pool.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [pool.created_at])
    const updatedAt = useMemo(() => {
        if (pool.updated_at) {
            return DateTime.fromISO(pool.updated_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [pool.updated_at])

    const {mutate} = useMutation(deleteItem)

    const handleEditPoolClick = () => {
        push(urlStateKey, {
            ...urlState,
            currentPanel: "Pool",
            id: pool.id,
        })
    }

    const handleDeletePoolClick = () => {
        mutate(
            {
                key: "pools",
                id: pool.id,
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: () => {
                    const queryKey = ["pools"]
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

    const handlePoolClick = () => {
        push(urlStateKey, {
            ...urlState,
            pool: pool.id,
        })
        setSelectedPool(pool.id)
    }

    return (
        <DataGridRow>
            <DataGridCell>
                <ListItemSpinner data={pool} onClick={handlePoolClick} className={`cursor-pointer ${isActive ? "jn-text-theme-accent" : "hover:text-theme-accent"}`}/>
            </DataGridCell>
            <DataGridCell>{pool.domains?.length || 0}/{pool.members?.length || 0}/{pool.monitors?.length || 0}</DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={pool} />
            </DataGridCell>
            <DataGridCell>
                {/* Use <Stack> to align and space elements: */}
                <Stack gap="1.5">
                    <Icon
                        icon="edit"
                        size="18"
                        className="leading-none"
                        onClick={handleEditPoolClick}
                    />
                    <Icon
                        icon="deleteForever"
                        size="18"
                        className="leading-none"
                        onClick={handleDeletePoolClick}
                    />
                </Stack>
            </DataGridCell>
        </DataGridRow>
    )
}

export default PoolListItem
