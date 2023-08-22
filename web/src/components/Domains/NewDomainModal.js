import React, {useState} from "react"

import {authStore, useStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../actions"
import {Checkbox, Modal, Select, SelectOption, Stack, TextInput} from "juno-ui-components"
import {currentState, push} from "url-state-provider"
import {Error} from "../Components";

const NewDomainModal = () => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const urlState = currentState(urlStateKey)
    const [formState, setFormState] = useState({
        name: "",
        provider: "akamai",
        fqdn: "",
        record_type: "A",
        admin_state_up: true,
    })
    const queryClient = useQueryClient()
    const {error, mutate} = useMutation(createItem)

    const closeNewDomainModal = () => {
        push(urlStateKey, {...urlState, currentModal: ""})
    }

    const onSubmit = () => {
        mutate(
            {
                key: "domains",
                formState: {"domain": formState},
                endpoint: auth?.endpoint,
                token: auth?.token,
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["domains", data.domain.id], data)
                    queryClient
                        .invalidateQueries(["domains"])
                        .then(closeNewDomainModal)
                }
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Modal
            title="Add new Domain"
            open
            onCancel={closeNewDomainModal}
            confirmButtonLabel="Save new Domain"
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
                <Select
                    label="Provider"
                    value={formState.provider}
                    onChange={(target) => setFormState({...formState, provider: target})}
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
                    value={formState.fqdn}
                    onChange={handleChange}
                    required
                />
                <Select
                    label="Record Type"
                    value={formState.record_type}
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
        </Modal>
    )
}

export default NewDomainModal
