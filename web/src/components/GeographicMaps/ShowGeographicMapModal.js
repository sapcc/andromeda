// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useEffect, useMemo, useState} from "react"

import {authStore, urlStore} from "../../store"
import {useQuery} from '@tanstack/react-query'
import {fetchItem} from "../../actions"
import {Badge, DataGrid, DataGridCell, DataGridHeadCell, DataGridRow, Modal, Spinner, Stack} from "@cloudoperators/juno-ui-components"
import {Error} from "../Components";
import {continents, countries, getEmojiFlag} from "countries-list";
import {DateTime} from "luxon";

const ShowGeographicMapModal = () => {
    const auth = authStore((state) => state.auth)
    const [id, closeModal] = urlStore((state) => [state.id, state.closeModal])
    const [geomap, setGeomap] = useState({})
    const {data, error, isSuccess, isLoading} = useQuery({
        queryKey: ["geomaps", id],
        queryFn: fetchItem,
        meta: auth,
    })

    // update formState when data is fetched
    useEffect(() => {
        if (data) {
            setGeomap(data.geomap)
        }
    }, [data]);

    const createdAt = useMemo(() => {
        if (geomap.created_at) {
            return DateTime.fromISO(geomap.created_at).toLocaleString(
                DateTime.DATETIME_SHORT
            )
        }
    }, [geomap.created_at])

    return (
        <Modal
            title={`Geographical map`}
            size="large"
            open
            onCancel={closeModal}
        >
            {/* Error Bar */}
            <Error error={error} />

            {isLoading && <Spinner size="large" />}
            {isSuccess && (
                <Stack direction="vertical" gap="1.5">
                    <p>ID <Badge>{geomap.id}</Badge></p>
                    {geomap.name && <p>Name <Badge>{geomap.name}</Badge></p>}
                    <p>Default Datacenter: <Badge>{geomap.default_datacenter}</Badge></p>
                    <p>Created: <Badge>{createdAt}</Badge></p>
                    <p>Updated: <Badge>{createdAt}</Badge></p>
                    <p>Scope: <Badge>{geomap.scope}</Badge></p>
                    <DataGrid columns={3}>
                        <DataGridRow>
                            <DataGridHeadCell>
                                Datacenter
                            </DataGridHeadCell>
                            <DataGridHeadCell>
                                Country
                            </DataGridHeadCell>
                            <DataGridHeadCell>
                                Continent
                            </DataGridHeadCell>
                        </DataGridRow>
                        {geomap.assignments?.map(o => (
                            <DataGridRow key={o.country}>
                                <DataGridCell>
                                    {o.datacenter.substring(0, 10)}...
                                </DataGridCell>
                                <DataGridCell>
                                    {getEmojiFlag(o.country)} {countries[o.country].name}
                                </DataGridCell>
                                <DataGridCell>
                                    {continents[countries[o.country].continent]}
                                </DataGridCell>
                            </DataGridRow>
                        ))}
                    </DataGrid>
                </Stack>
            )}
        </Modal>
    )
}

export default ShowGeographicMapModal
