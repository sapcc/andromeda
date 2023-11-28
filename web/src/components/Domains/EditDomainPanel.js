import React, {useEffect, useState} from "react"

import {
    Button,
    Checkbox,
    Form,
    PanelBody,
    PanelFooter,
    Select,
    SelectOption,
    Stack,
    TextInput,
} from "juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../actions"
import {Error, Loading} from "../Components";

const EditDomainPanel = ({closeCallback}) => {
    const id = urlStore((state) => state.id)
    const auth = authStore((state) => state.auth)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        provider: undefined,
        fqdn: undefined,
        record_type: undefined,
        admin_state_up: undefined,
    })

    const {data, isLoading} = useQuery({
        queryKey: ["domains", id],
        queryFn: fetchItem,
        meta: auth,
        refetchOnWindowFocus: false,
    })
    const mutation = useMutation({mutationFn: updateItem})

    // update formState when data is fetched
    useEffect(() => {
        if (data) {
            setFormState(updateAttributes(formState, data.domain))
        }
    }, [data]);

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "domains",
                id: id,
                formState: {"domain": formState},
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["domains", data.domain.id], data)
                    queryClient
                        .invalidateQueries({queryKey: ["domains"]})
                        .then(closeCallback)
                },
                onError: setError
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <PanelBody
            footer={
                <PanelFooter>
                    <Button label="Cancel" variant="subdued" onClick={closeCallback}/>
                    <Button label="Save" variant="primary" onClick={onSubmit}/>
                </PanelFooter>
            }
        >
            {/* Error Bar */}
            <Error error={error}/>

            {/* Loading indicator for page content */}
            <Loading isLoading={isLoading || mutation.isLoading}/>

            <Form>
                <Stack direction="vertical" gap="2">
                    <Checkbox
                        id="selectable"
                        label="Enabled"
                        checked={formState.admin_state_up}
                        disabled={isLoading}
                        onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                    />
                    <TextInput
                        label="Name"
                        name="name"
                        value={formState?.name}
                        disabled={isLoading}
                        onChange={handleChange}
                    />
                    <Select
                        label="Provider"
                        name="provider"
                        value={formState?.provider}
                        disabled={isLoading}
                        onChange={handleChange}
                        required
                    >
                        <SelectOption
                            key="akamai"
                            label="Akamai"
                            value="akamai"
                        />
                    </Select>
                    <TextInput
                        label="Fully Qualified Domain Name"
                        name="fqdn"
                        value={formState?.fqdn}
                        disabled={isLoading}
                        onChange={handleChange}
                        required
                    />
                    <Select
                        label="Record Type"
                        value={formState?.record_type}
                        disabled={isLoading}
                        onChange={(target) => setFormState({...formState, record_type: target})}
                    >
                        <SelectOption
                            label="A"
                            value="A"
                        />
                        <SelectOption
                            label="AAAA"
                            value="AAAA"
                        />
                        <SelectOption
                            label="CNAME"
                            value="CNAME"
                        />
                        <SelectOption
                            label="MX"
                            value="MX"
                        />
                    </Select>
                </Stack>
            </Form>
        </PanelBody>
    )
}

export default EditDomainPanel
