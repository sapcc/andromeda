import React, {useState} from "react"

import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../actions"
import {Box, Button, Icon, JsonViewer, Label, Modal, Stack, TextInput} from "juno-ui-components"
import {Error} from "../Components";
import {continents, countries} from "countries-list";
import DatacenterSelect from "./DatacenterSelect";

const NewGeographicMapModal = () => {
    const auth = authStore((state) => state.auth)
    const closeModal = urlStore((state) => state.closeModal)
    const queryClient = useQueryClient()
    const [formState, setFormState] = useState({
        name: "",
        default_datacenter: "",
        scope: "private",
        provider: "akamai",
    })
    const [datacenter, setDatacenter] = useState()
    const [error, setError] = useState()
    const [assignments, setAssignments] = useState({})
    const [expandContinent, setExpandContinent] = useState(Object.fromEntries(Object.keys(continents).map(c => [c, false])))
    const mutation = useMutation({mutationFn: createItem})

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "geomaps",
                token: auth?.token,
                endpoint: auth?.endpoint,
                formState: {
                    geomap: {
                        ...formState,
                        assignments: Object.keys(assignments).map(key => ({country: key, datacenter: assignments[key]}))
                    }
                },
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["geomaps", data.geomap.id], data)
                    queryClient
                        .invalidateQueries({queryKey: ["geomaps"]})
                        .then(closeModal)
                },
                onError: setError
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    const addAssignments = (entries) => {
        if (datacenter) {
            setAssignments({
                ...assignments,
                ...Object.fromEntries(entries)
            })
            setError(undefined)
        } else {
            setError({message: "Select 'Assign Datacenter' first"})
        }
    }

    const ContinentList = ({continent}) => {
        var vals = Object.entries(countries)
            .filter(([key, value]) => value.continent === continent)
            .filter(c => !(c[0] in assignments))

        return (
            <Stack className="pl-12 pt-1" direction="vertical" alignment="start" gap="1">
                {vals.map(([key, country]) => {
                    return (
                        <Button key={key} icon="addCircle" size="small" label={`${country.name} ${country.emoji}`}
                                onClick={() => addAssignments([[key, datacenter]])}/>
                    )
                })}
            </Stack>
        )
    }

    return (
        <Modal
            title="Add new Geographical Map"
            size="large"
            open
            onCancel={closeModal}
            confirmButtonLabel="Save new Geographical Map"
            onConfirm={onSubmit}
        >
            {/* Error Bar */}
            <Error error={error} />

            <Stack direction="vertical" gap="1.5">

                <TextInput
                    label="Name"
                    name="name"
                    value={formState.name}
                    onChange={handleChange}
                />

                <DatacenterSelect setError={setError}
                                  setDatacenter={(value) => setFormState({
                                      ...formState, default_datacenter: value})}
                                  label="Default Datacenter"
                />

                <Box>
                    <Label text="Datacenter - Country association"/>
                    <DatacenterSelect className="jn-mt-2" setDatacenter={setDatacenter} setError={setError} label="Assign Datacenter" />
                    <Stack direction="horizontal" gap="2" className="jn-py-2">
                        <Stack direction="vertical" gap="1" className="basis-1/3">
                            {Object.keys(continents).filter((continentKey, _) => Object.entries(countries)
                                .filter(([countryKey, value]) => value.continent === continentKey && !(countryKey in assignments)).length > 0).map((key, index) => (
                                <div key={index} className="jn-p-1">
                                    <Icon icon={expandContinent[key] ? "expandMore" : "chevronRight"} onClick={() => {
                                        setExpandContinent({...expandContinent, [key]: !expandContinent[key]})
                                    }}
                                    />
                                    <Button icon="addCircle" size="small" label={continents[key]} onClick={(e) => {
                                        addAssignments(
                                            Object.entries(countries)
                                                .filter(([_, value]) => value.continent === key)
                                                .map(([countryKey, _]) => [countryKey, datacenter])
                                        )
                                    }} />
                                    {expandContinent[key] && <ContinentList continent={key}/>}
                                </div>
                            ))}
                        </Stack>
                        <Stack direction="vertical" distribution="between" gap="1" className="basis-2/3">
                            <Label text="Assignments" />
                            <JsonViewer data={assignments} />
                            <Button
                                className="mx-0"
                                label="Clear assignments"
                                size="small"
                                onClick={() => setAssignments({})}
                                variant="primary-danger"
                            />
                        </Stack>
                    </Stack>
                </Box>

            </Stack>
        </Modal>
    )
}

export default NewGeographicMapModal
