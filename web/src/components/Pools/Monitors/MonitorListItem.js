// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useMemo} from "react"

import {DataGridCell, DataGridRow, Icon, Stack} from "@cloudoperators/juno-ui-components"
import {authStore, urlStore} from "../../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {deleteItem} from "../../../actions"
import {DateTime} from "luxon";
import {ListItemSpinner, ListItemStatus} from "../../Components";

const MonitorListItem = ({monitor, setError}) => {
    const auth = authStore((state) => state.auth)
    const openPanel = urlStore((state) => state.openPanel)
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

    const {mutate} = useMutation({mutationFn: deleteItem})

    const handleEditMonitorClick = () => {
        openPanel("Monitor", monitor.id)
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
                    queryClient.invalidateQueries({
                        queryKey: queryKey
                    })
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
                <ListItemSpinner data={monitor} />
            </DataGridCell>
            <DataGridCell>{createdAt}</DataGridCell>
            <DataGridCell>{updatedAt}</DataGridCell>
            <DataGridCell>
                <ListItemStatus data={monitor} />
            </DataGridCell>
            <DataGridCell className="jn-items-end">
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
