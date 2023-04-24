import React, {useState} from "react"

import {Button, CheckboxRow, Form, PanelBody, PanelFooter, Spinner, TextInputRow,} from "juno-ui-components"
import {authStore, useStore} from "../../store"
import {currentState} from "url-state-provider"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../actions"
import DomainMenu from "./DomainMenu";
import {Error, Loading} from "../Components";

const EditPoolPanel = ({closeCallback}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const urlState = currentState(urlStateKey)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined, admin_state_up: true, domains: [],
    })

    const {isLoading} = useQuery(
        ["pools", urlState.id],
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
                id: urlState.id,
                endpoint: auth?.endpoint,
                token: auth?.token,
                formState: {pool: formState},
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["pools", data.pool.id, endpoint], data)
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

                <CheckboxRow
                    id="selectable"
                    label="Enabled"
                    disabled={isLoading}
                    checked={formState.admin_state_up}
                    onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                />
                <TextInputRow
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
