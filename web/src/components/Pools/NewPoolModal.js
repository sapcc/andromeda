import React, {useState} from "react"

import {authStore, useStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../actions"
import {CheckboxRow, Message, Modal, TextInputRow} from "juno-ui-components"
import DomainMenu from "./DomainMenu";
import {currentState, push} from "url-state-provider";
import {Error, Loading} from "../Components";

const NewPoolModal = () => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        admin_state_up: true,
        domains: [],
    })

    const {mutate} = useMutation(createItem)

    const closeModal = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {...urlState, currentModal: ""})
    }

    const onSubmit = () => {
        mutate({
            key: "pools",
            formState: {"pool": formState},
            endpoint: auth?.endpoint,
            token: auth?.token,
        }, {
            onSuccess: (data) => {
                queryClient
                    .setQueryData(["pools", data.pool.id], data)
                queryClient
                    .invalidateQueries("pools")
                    .then(closeModal)
            },
            onError: setError,
        })
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Modal
            title="Add new Pool"
            open
            onCancel={closeModal}
            confirmButtonLabel="Save new Pool"
            onConfirm={onSubmit}
        >
            {/* Error Bar */}
            <Error error={error} />

            <CheckboxRow
                id="selectable"
                label="Enabled"
                checked={formState.admin_state_up}
                onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
            />
            <TextInputRow
                label="Name"
                name="name"
                value={formState.name}
                onChange={handleChange}
            />

            Associated Domain(s):
            <DomainMenu formState={formState} setFormState={setFormState} setError={setError}/>
        </Modal>
    )
}

export default NewPoolModal
