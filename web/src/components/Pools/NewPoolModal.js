import React, {useState} from "react"

import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../actions"
import {Checkbox, Modal, TextInput} from "juno-ui-components"
import DomainMenu from "./DomainMenu";
import {Error} from "../Components";

const NewPoolModal = () => {
    const auth = authStore((state) => state.auth)
    const setModal = urlStore((state) => state.openModal)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        admin_state_up: true,
        domains: [],
    })

    const {mutate} = useMutation(createItem)
    const closeModal = () => setModal(null)

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
                    .invalidateQueries({queryKey: ["pools"]})
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

            Associated Domain(s):
            <DomainMenu formState={formState} setFormState={setFormState} setError={setError}/>
        </Modal>
    )
}

export default NewPoolModal
