import React, {useState} from "react"

import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../actions"
import {Checkbox, Modal, Select, SelectOption, Stack, TextInput} from "juno-ui-components"
import {Error} from "../Components";

const NewDatacenterModal = () => {
    const auth = authStore((state) => state.auth)
    const closeModal = urlStore((state) => state.closeModal)
    const queryClient = useQueryClient()
    const [formState, setFormState] = useState({
        name: undefined,
        admin_state_up: true,
        continent: undefined,
        country: undefined,
        state_or_province: undefined,
        city: undefined,
        longitude: 13.4,
        latitude: 52.52,
        provider: "akamai",
    })

    const {error, mutate} = useMutation(createItem)
    const onSubmit = () => {
        mutate(
            {
                key: "datacenters",
                token: auth?.token,
                endpoint: auth?.endpoint,
                formState: {"datacenter": formState},
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["datacenters", data.datacenter.id], data)
                    queryClient
                        .invalidateQueries( {queryKey: ["datacenters"]})
                        .then(closeModal)
                }
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Modal
            title="Add new Datacenter"
            open
            onCancel={closeModal}
            confirmButtonLabel="Save new Datacenter"
            onConfirm={onSubmit}
        >
            {/* Error Bar */}
            <Error error={error} />

            <Checkbox
                id="selectable"
                label="Enabled"
                checked={formState.admin_state_up}
                onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
            />
            <TextInput
                label="Name"
                name="name"
                value={formState.name}
                onChange={handleChange}
            />
            <Select
                label="Continent"
                helptext="A two-letter code that specifies the continent where the data center maps to."
                value={formState.continent}
                onChange={(target) => setFormState({...formState, continent: target})}
            >
                <SelectOption label="Unknown"/>
                <SelectOption key="AF" label="AF - Africa" value="AF"/>
                <SelectOption key="AN" label="AN - Antarctica" value="AN"/>
                <SelectOption key="AS" label="AS - Asia" value="AS"/>
                <SelectOption key="EU" label="EU - Europe" value="EU"/>
                <SelectOption key="NA" label="NA - North america" value="NA"/>
                <SelectOption key="OC" label="OC - Oceania" value="OC"/>
                <SelectOption key="SA" label="SA - South america" value="SA"/>
            </Select>
            <TextInput
                label="Country"
                name="country"
                helptext="A two-letter ISO 3166 country code that specifies the country where the data center maps to."
                value={formState.country}
                onChange={handleChange}
            />
            <TextInput
                label="State or Province"
                name="state_or_province"
                helptext="The name of the state or province where the data center is located."
                value={formState.state_or_province}
                onChange={handleChange}
            />
            <TextInput
                label="City"
                name="city"
                helptext="The name of the city where the data center is located."
                value={formState.city}
                onChange={handleChange}
            />
            <Stack gap="2">
                <TextInput
                    label="Longitude"
                    type="number"
                    helptext="Specifies the geographical longitude of the data center's position."
                    value={formState.longitude?.toString()}
                    onChange={(event) => setFormState({...formState, longitude: parseFloat(event.target.value)})}
                />
                <TextInput
                    label="Latitude"
                    type="number"
                    helptext="Specifies the geographic latitude of the data center's position."
                    value={formState.latitude?.toString()}
                    onChange={(event) => setFormState({...formState, latitude: parseFloat(event.target.value)})}
                />
            </Stack>
        </Modal>
    )
}

export default NewDatacenterModal
