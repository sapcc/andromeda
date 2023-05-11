import React, {useMemo, useState} from "react"

import {authStore, useStore} from "../../store"
import {useQuery} from '@tanstack/react-query'
import {fetchItem} from "../../actions"
import {Badge, DataGrid, DataGridCell, DataGridHeadCell, DataGridRow, Modal, Spinner, Stack} from "juno-ui-components"
import {currentState, push} from "url-state-provider"
import {Error} from "../Components";
import {continents, countries} from "countries-list";
import {DateTime} from "luxon";

const ShowGeographicMapModal = ({id}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const [error, setError] = useState()
    const [geomap, setGeomap] = useState({})
    const {isSuccess, isLoading} = useQuery(
        ["geomaps", id],
        fetchItem,
        {
            meta: auth,
            onSuccess: (data) => setGeomap(data.geomap),
            onError: setError
        })

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


    const closeShowGeographicMapModal = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {...urlState, currentModal: ""})
    }

    return (
        <Modal
            title={`Geographical map`}
            size="large"
            open
            onCancel={closeShowGeographicMapModal}
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
                            <DataGridRow>
                                <DataGridCell>
                                    {o.datacenter}
                                </DataGridCell>
                                <DataGridCell>
                                    {countries[o.country].name} {countries[o.country].emoji}
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