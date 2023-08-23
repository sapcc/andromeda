import React, {useState} from "react"
import {
    Button,
    CheckboxRow,
    Form,
    PanelBody,
    PanelFooter,
    SelectOption,
    SelectRow,
    Stack,
    TextInputRow,
} from "juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../actions"
import {Error, Loading} from "../Components";
import {produce} from "immer";

const EditDatacenterPanel = ({closeCallback}) => {
    const auth = authStore((state) => state.auth)
    const id = urlStore((state) => state.id)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        admin_state_up: true,
        continent: undefined,
        country: undefined,
        state_or_province: undefined,
        city: undefined,
        longitude: undefined,
        latitude: undefined,
        provider: undefined,
    })

    const {isLoading} = useQuery(
        ["datacenters", id],
        fetchItem,
        {
            meta: auth,
            onError: setError,
            onSuccess: (data) => setFormState(updateAttributes(formState, data.datacenter)),
            refetchOnWindowFocus: false,
        })
    const mutation = useMutation(updateItem)

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "datacenters",
                id: id,
                formState: {"datacenter": formState},
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["datacenters", data.datacenter.id], data)
                    queryClient
                        .invalidateQueries("datacenters")
                        .then(closeCallback)
                },
                onError: setError,
            }
        )
    }

    const handleChange = (event) => {
        setFormState(
            produce((draft) => {
                draft[event.target.name] = event.target.value
            })
        );
        //setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Form onSubmit={(e) => e.preventDefault()}>
            <PanelBody
                footer={
                    <PanelFooter>
                        <Button label="Cancel" variant="subdued" onClick={closeCallback}/>
                        <Button label="Save" type="submit" variant="primary" onClick={onSubmit}/>
                    </PanelFooter>
                }
            >
                {/* Error Bar */}
                <Error error={error} />

                {/* Loading indicator for page content */}
                <Loading isLoading={isLoading || mutation.isLoading} />

                <CheckboxRow
                    id="selectable"
                    label="Enabled"
                    disabled={isLoading}
                    checked={formState?.admin_state_up}
                    onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                />
                <TextInputRow
                    label="Name"
                    disabled={isLoading}
                    name="name"
                    value={formState?.name}
                    onChange={handleChange}
                />
                <SelectRow
                    label="Continent"
                    disabled={isLoading}
                    helptext="A two-letter code that specifies the continent where the data center maps to."
                    value={formState?.continent}
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
                </SelectRow>
                <TextInputRow
                    label="Country"
                    disabled={isLoading}
                    name="country"
                    helptext="A two-letter ISO 3166 country code that specifies the country where the data center maps to."
                    value={formState?.country}
                    onChange={handleChange}
                />
                <TextInputRow
                    label="State or Province"
                    disabled={isLoading}
                    name="state_or_province"
                    helptext="The name of the state or province where the data center is located."
                    value={formState.state_or_province}
                    onChange={handleChange}
                />
                <TextInputRow
                    label="City"
                    disabled={isLoading}
                    name="city"
                    helptext="The name of the city where the data center is located."
                    value={formState.city}
                    onChange={handleChange}
                />
                <Stack gap="2">
                    <TextInputRow
                        label="Longitude"
                        disabled={isLoading}
                        type="number"
                        helptext="Specifies the geographical longitude of the data center's position."
                        value={formState.longitude?.toString()}
                        onChange={(event) => setFormState({...formState, longitude: parseFloat(event.target.value)})}
                    />
                    <TextInputRow
                        label="Latitude"
                        disabled={isLoading}
                        type="number"
                        helptext="Specifies the geographic latitude of the data center's position."
                        value={formState.latitude?.toString()}
                        onChange={(event) => setFormState({...formState, latitude: parseFloat(event.target.value)})}
                    />
                </Stack>
            </PanelBody>
        </Form>
    )
}

export default EditDatacenterPanel
