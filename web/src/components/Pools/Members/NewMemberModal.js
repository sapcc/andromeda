import React, {useMemo, useState} from "react"

import {authStore, urlStore} from "../../../store"
import {createItem} from "../../../actions"
import {Badge, Checkbox, Modal, Stack, TextInput} from "@cloudoperators/juno-ui-components"
import {useMutation, useQueryClient} from "@tanstack/react-query";
import DatacenterMenu from "./DatacenterMenu";
import {Error} from "../../Components";

const NewMemberModal = () => {
    const auth = authStore((state) => state.auth)
    const [closeModal, pool] = urlStore((state) => [state.closeModal, state.pool])
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: "",
        address: "",
        datacenter_id: "",
        port: 0,
        admin_state_up: true,
        pool_id: pool,
    })

    const {mutate} = useMutation({mutationFn: createItem})
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
                        .invalidateQueries({queryKey: ["members"]})
                        .then(closeModal)
                },
                onError: setError
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };
    const heading = useMemo(() => {
        return (
            <>
                Add new Member
                <p><small>Pool <Badge>{pool}</Badge></small></p>
            </>
        )
    }, [pool])

    return (
        <Modal
            heading={heading}
            open
            onCancel={closeModal}
            confirmButtonLabel="Save new Member"
            onConfirm={onSubmit}
        >
            <Stack direction="vertical" gap="2">
                {/* Error Bar */}
                <Error error={error}/>

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
                <TextInput
                    label="IP address"
                    name="address"
                    value={formState.address}
                    onChange={handleChange}
                    required
                />
                <TextInput
                    label="Port"
                    type="number"
                    value={formState.port?.toString()}
                    onChange={(event) => setFormState({...formState, port: parseInt(event.target.value)})}
                    required
                />
                Select a Datacenter
                <DatacenterMenu formState={formState} setFormState={setFormState} setError={setError}/>
            </Stack>
        </Modal>
    )
}

export default NewMemberModal
