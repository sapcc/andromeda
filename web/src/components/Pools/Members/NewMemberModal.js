import React, {useState} from "react"

import {authStore, useStore} from "../../../store"
import {createItem} from "../../../actions"
import {CheckboxRow, Message, Modal, TextInputRow} from "juno-ui-components"
import {useMutation, useQueryClient} from "@tanstack/react-query";
import DatacenterMenu from "./DatacenterMenu";
import {currentState, push} from "url-state-provider";
import {Error, Loading} from "../../Components";

const NewMemberModal = () => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const urlState = currentState(urlStateKey)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        address: undefined,
        datacenter_id: undefined,
        port: undefined,
        admin_state_up: true,
        pool_id: urlState?.pool,
    })

    const {mutate} = useMutation(createItem)

    const closeModal = () => {
        push(urlStateKey, {...urlState, currentModal: ""})
    }

    const onSubmit = () => {
        mutate(
            {
                key: "members",
                endpoint: auth?.endpoint,
                token: auth?.token,
                formState: {"member": formState},
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["members", data.member.id], data)
                    queryClient
                        .invalidateQueries("members")
                        .then(closeModal)
                },
                onError: setError
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Modal
            title="Add new Member"
            open
            onCancel={closeModal}
            confirmButtonLabel="Save new Member"
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
            <TextInputRow
                label="IP address"
                name="address"
                value={formState.address}
                onChange={handleChange}
                required
            />
            <TextInputRow
                label="Port"
                type="number"
                value={formState.port?.toString()}
                onChange={(event) => setFormState({...formState, port: parseInt(event.target.value)})}
                required
            />
            Select a Datacenter
            <DatacenterMenu formState={formState} setFormState={setFormState} setError={setError}/>
        </Modal>
    )
}

export default NewMemberModal
