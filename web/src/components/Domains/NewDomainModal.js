// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useState} from "react"

import {authStore, urlStore} from "../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../actions"
import {Checkbox, Modal, Select, SelectOption, Stack, TextInput} from "@cloudoperators/juno-ui-components"
import {Error, Loading} from "../Components";

const NewDomainModal = () => {
    const auth = authStore((state) => state.auth)
    const closeModal = urlStore((state) => state.closeModal)
    const [formState, setFormState] = useState({
        name: "",
        provider: "akamai",
        fqdn: "",
        record_type: "A",
        admin_state_up: true,
    })
    const queryClient = useQueryClient()
    const {error, mutate, isPending} = useMutation({mutationFn: createItem})

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
                    queryClient.invalidateQueries({
                        queryKey: ["domains"]
                    })
                        .then(closeModal)
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
            onCancel={closeModal}
            confirmButtonLabel="Save new Domain"
            onConfirm={onSubmit}
        >
            <Stack direction="vertical" gap="2">
                {/* Error Bar */}
                <Error error={error}/>

                {/* Loading indicator for page content */}
                <Loading isLoading={isPending}/>

                <Checkbox
                    id="selectable"
                    disabled={isPending}
                    label="Enabled"
                    checked={formState.admin_state_up}
                    onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                />
                <TextInput
                    label="Name"
                    disabled={isPending}
                    helptext="The name of the domain, e.g. example-com"
                    name="name"
                    value={formState.name}
                    onChange={handleChange}
                />
                <Select
                    label="Provider"
                    disabled={isPending}
                    value={formState.provider}
                    helptext="Currently, only Akamai is supported."
                    onChange={(target) => setFormState({...formState, provider: target})}
                    required
                >
                    <SelectOption
                        label="Akamai"
                        value="akamai"
                    />
                </Select>
                <TextInput
                    label="Fully Qualified Domain Name"
                    disabled={isPending}
                    helptext="The fully qualified domain name of the domain, e.g. example.com"
                    name="fqdn"
                    value={formState.fqdn}
                    onChange={handleChange}
                    required
                />
                <Select
                    label="Record Type"
                    disabled={isPending}
                    helptext="The type of DNS record to create."
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
