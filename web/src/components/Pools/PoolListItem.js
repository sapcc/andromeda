import React, {useMemo} from "react"

import {DataGridCell, DataGridRow, Icon, Stack} from "@cloudoperators/juno-ui-components"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {deleteItem} from "../../actions"
import {DateTime} from "luxon";
import {authStore, urlStore} from "../../store"
import {ListItemSpinner, ListItemStatus} from "../Components";

const PoolListItem = ({pool, isActive, setError}) => {
    const auth = authStore((state) => state.auth)
    const [openPanel, setSelectedPool] = urlStore((state) => [state.openPanel, state.setPool])
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

    const {mutate} = useMutation({mutationFn: deleteItem})

    const handlePoolClick = () => setSelectedPool(pool.id)
    const handleEditPoolClick = () => openPanel("Pool", pool.id)
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
                <ListItemSpinner
                    data={pool} onClick={handlePoolClick}
                    className={`cursor-pointer ${isActive ? "jn-text-theme-accent" : "hover:text-theme-accent"}`}
                />
            </DataGridCell>
            <DataGridCell>{pool.domains?.length || 0}/{pool.members?.length || 0}/{pool.monitors?.length || 0}</DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={pool} />
            </DataGridCell>
            <DataGridCell className="jn-items-end">
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
