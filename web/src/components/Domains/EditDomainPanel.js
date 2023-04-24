import React, {useState} from "react"

import {
    Button,
    Checkbox, CheckboxRow,
    Form,
    Message,
    PanelBody,
    PanelFooter,
    SelectOption,
    SelectRow,
    Spinner,
    Stack,
    TextInputRow,
} from "juno-ui-components"
import {authStore, useStore} from "../../store"
import {currentState} from "url-state-provider"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../actions"
import {Error, Loading} from "../Components";

const EditDomainPanel = ({closeCallback}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const urlState = currentState(urlStateKey)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        provider: undefined,
        fqdn: undefined,
        record_type: undefined,
        admin_state_up: undefined,
    })

    const {isLoading} = useQuery(["domains", urlState.id],
        fetchItem,
        {
            meta: auth,
            onError: setError,
            onSuccess: (data) => setFormState(updateAttributes(formState, data.domain)),
            refetchOnWindowFocus: false,
        })
    const mutation = useMutation(updateItem)

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "domains",
                id: urlState.id,
                formState: {"domain": formState},
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["domains", data.domain.id], data)
                    queryClient
                        .invalidateQueries("domains")
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
            <Error error={error} />

            {/* Loading indicator for page content */}
            <Loading isLoading={isLoading || mutation.isLoading} />

            <Form>
                <CheckboxRow
                    id="selectable"
                    label="Enabled"
                    checked={formState.admin_state_up}
                    disabled={isLoading}
                    onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                />
                <TextInputRow
                    label="Name"
                    name="name"
                    value={formState?.name}
                    disabled={isLoading}
                    onChange={handleChange}
                />
                <SelectRow
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
                </SelectRow>
                <TextInputRow
                    label="Fully Qualified Domain Name"
                    name="fqdn"
                    value={formState?.fqdn}
                    disabled={isLoading}
                    onChange={handleChange}
                    required
                />
                <SelectRow
                    label="Record Type"
                    name="record_type"
                    value={formState?.record_type}
                    disabled={isLoading}
                    onChange={handleChange}
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
                </SelectRow>
            </Form>
        </PanelBody>
    )
}

export default EditDomainPanel
