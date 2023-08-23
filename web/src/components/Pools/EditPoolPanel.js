import React, {useState} from "react"

import {Button, Checkbox, Form, PanelBody, PanelFooter, Spinner, TextInput,} from "juno-ui-components"
import {authStore, urlStore} from "../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../actions"
import DomainMenu from "./DomainMenu";
import {Error, Loading} from "../Components";

const EditPoolPanel = ({closeCallback}) => {
    const auth = authStore((state) => state.auth)
    const id = urlStore((state) => state.id)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined, admin_state_up: true, domains: [],
    })

    const {isLoading} = useQuery(
        ["pools", id],
        fetchItem,
        {
            meta: auth,
            onError: setError,
            onSuccess: (data) => setFormState(updateAttributes(formState, data.pool)),
            refetchOnWindowFocus: false,
        })
    const mutation = useMutation(updateItem)

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "pools",
                id: id,
                endpoint: auth?.endpoint,
                token: auth?.token,
                formState: {pool: formState},
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["pools", data.pool.id], data)
                    queryClient
                        .setQueryDefaults([], {refetchInterval: 5000})
                    queryClient
                        .invalidateQueries("pools")
                        .then(closeCallback)
                },
                onError: setError
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    }

    return (
        <Form onSubmit={(e) => e.preventDefault()}>
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

                <Checkbox
                    id="selectable"
                    label="Enabled"
                    disabled={isLoading}
                    checked={formState.admin_state_up}
                    onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                />
                <TextInput
                    label="Name"
                    disabled={isLoading}
                    name="name"
                    value={formState.name}
                    onChange={handleChange}
                />
                Associated Domain(s):
                {isLoading && <Spinner/>}
                <DomainMenu formState={formState} setFormState={setFormState} setError={setError}/>
            </PanelBody>
        </Form>
    )
}

export default EditPoolPanel
