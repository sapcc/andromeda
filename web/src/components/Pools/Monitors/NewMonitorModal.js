import React, {useState} from "react"

import {authStore, useStore} from "../../../store"
import {useMutation, useQueryClient} from '@tanstack/react-query'
import {createItem} from "../../../actions"
import {
    Button,
    ButtonRow,
    CheckboxRow,
    Modal,
    SelectOption,
    SelectRow,
    Stack,
    TextareaRow,
    TextInputRow
} from "juno-ui-components"
import {currentState, push} from "url-state-provider"
import {Error} from "../../Components";

const NewMonitorModal = () => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const auth = authStore((state) => state.auth)
    const urlState = currentState(urlStateKey)
    const queryClient = useQueryClient()
    const [advancedSettings, setAdvancedSettings] = useState(false)
    const [formState, setFormState] = useState({
        name: "",
        pool_id: urlState?.pool,
        send: null,
        receive: null,
        timeout: 10,
        type: "ICMP",
        interval: 60,
        admin_state_up: true,
    })

    const {error, mutate} = useMutation(createItem)

    const closeNewMonitorModal = () => {
        const urlState = currentState(urlStateKey)
        push(urlStateKey, {...urlState, currentModal: ""})
    }

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
                        .invalidateQueries(`monitors`)
                        .then(closeNewMonitorModal)
                }
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Modal
            title="Add new Monitor"
            open
            onCancel={closeNewMonitorModal}
            confirmButtonLabel="Save new Monitor"
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
                label="Interval"
                name="interval"
                type="number"
                value={formState.interval.toString()}
                onChange={(event) => setFormState({...formState, interval: parseInt(event.target.value)})}
            />
            <TextInputRow
                label="Timeout"
                name="timeout"
                type="number"
                value={formState.timeout.toString()}
                onChange={(event) => setFormState({...formState, timeout: parseInt(event.target.value)})}
            />
            <SelectRow
                label="Type"
                name="type"
                value={formState?.type}
                onChange={handleChange}
            >
                <SelectOption key="icmp" value="ICMP" label="ICMP" />
                <SelectOption key="http" value="HTTP" label="HTTP" />
                <SelectOption key="https" value="HTTPS" label="HTTPS" />
                <SelectOption key="tcp" value="TCP" label="TCP" />
                <SelectOption key="udp" value="UDP" label="UDP" />
            </SelectRow>
            <ButtonRow>
                <Button
                    label="Show advanced settings"
                    variant="subdued"
                    onClick={() => setAdvancedSettings(!advancedSettings)}
                />
            </ButtonRow>

            {advancedSettings && (
                <Stack gap="2" distribution="between">
                    <TextareaRow
                        label="Monitor send message"
                        value={formState.send || ""}
                        onChange={handleChange}
                    />
                    <TextareaRow
                        label="Monitor expected receive message"
                        value={formState.receive || ""}
                        onChange={handleChange}
                    />
                </Stack>)}

        </Modal>
    )
}

export default NewMonitorModal
