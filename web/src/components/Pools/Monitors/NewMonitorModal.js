// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useMemo, useState} from "react"

import {authStore, urlStore} from "../../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../../actions"
import {Badge, Checkbox, Modal, Select, SelectOption, Stack, Textarea, TextInput} from "@cloudoperators/juno-ui-components"
import {Error} from "../../Components";

const NewMonitorModal = () => {
    const auth = authStore((state) => state.auth)
    const [closeModal, pool] = urlStore((state) => [state.closeModal, state.pool])
    const queryClient = useQueryClient()
    const [formState, setFormState] = useState({
        name: "",
        pool_id: pool,
        send: "/",
        receive: null,
        timeout: 10,
        type: "HTTP",
        interval: 60,
        admin_state_up: true,
    })

    const {error, mutate} = useMutation({mutationFn: createItem})
    const onSubmit = () => {
        mutate(
            {
                key: "monitors",
                endpoint: auth?.endpoint,
                token: auth?.token,
                formState: {"monitor": formState},
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["monitors", data.monitor.id], data)
                    queryClient
                        .invalidateQueries({queryKey: ["monitors"]})
                        .then(closeModal)
                }
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };
    const heading = useMemo(() => {
        return (
            <>
                Add new Monitor
                <p><small>Pool <Badge>{pool}</Badge></small></p>
            </>
        )
    }, [pool])
    const isHTTP = (type) => ["HTTP", "HTTPS"].includes(type)

    return (
        <Modal
            heading={heading}
            open
            onCancel={closeModal}
            confirmButtonLabel="Save new Monitor"
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
                    label="Interval"
                    name="interval"
                    type="number"
                    value={formState.interval.toString()}
                    onChange={(event) => setFormState({...formState, interval: parseInt(event.target.value)})}
                />
                <TextInput
                    label="Timeout"
                    name="timeout"
                    type="number"
                    value={formState.timeout.toString()}
                    onChange={(event) => setFormState({...formState, timeout: parseInt(event.target.value)})}
                />
                <Select
                    label="Type"
                    value={formState?.type}
                    onChange={(target) => setFormState({...formState, ...{type: target, send: isHTTP(target) ? "/" : null}})}
                >
                    <SelectOption key="icmp" value="ICMP" label="ICMP (Unsupported on Akamai)"/>
                    <SelectOption key="http" value="HTTP" label="HTTP"/>
                    <SelectOption key="https" value="HTTPS" label="HTTPS"/>
                    <SelectOption key="tcp" value="TCP" label="TCP"/>
                    <SelectOption key="udp" value="UDP" label="UDP"/>
                </Select>

                {formState.type !== "ICMP" && (
                    <Stack gap="2" distribution="between">
                        <Textarea
                            name="send"
                            className={"flex-auto"}
                            label={isHTTP(formState.type) ? "HTTP Path" : "send data"}
                            value={formState.send || ""}
                            onChange={handleChange}
                        />
                        {formState.type === "TCP" && (
                            <Textarea
                                name="receive"
                                className={"flex-auto"}
                                label="receive data"
                                value={formState.receive || ""}
                                onChange={handleChange}
                            />
                        )}
                    </Stack>
                )}
            </Stack>
        </Modal>
    )
}

export default NewMonitorModal
