import React, {useState} from "react"

import {
    Button,
    CheckboxRow,
    Form,
    PanelBody,
    PanelFooter,
    SelectOption,
    SelectRow,
    Stack,
    TextareaRow,
    TextInputRow,
} from "juno-ui-components"
import {authStore, urlStore} from "../../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../../actions"
import {Error, Loading} from "../../Components";

const EditMonitorPanel = ({closeCallback}) => {
    const auth = authStore((state) => state.auth)
    const id = urlStore((state) => state.id)
    const queryClient = useQueryClient()
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        send: undefined,
        receive: undefined,
        timeout: undefined,
        type: undefined,
        interval: undefined,
        admin_state_up: undefined,
    })

    const {isLoading} = useQuery(
        ["monitors", id],
        fetchItem, {
            meta: auth,
            onError: setError,
            onSuccess: (data) => setFormState(updateAttributes(formState, data.monitor)),
            refetchOnWindowFocus: false,
        }
    )
    const mutation = useMutation(updateItem)

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "monitors",
                id: id,
                endpoint: auth?.endpoint,
                token: auth?.token,
                formState: {"monitor": formState},
            },
            {
                onSuccess: (data) => {
                    queryClient
                        .setQueryData(["monitors", data.monitor.id], data)
                    queryClient
                        .invalidateQueries("monitors")
                        .then(closeCallback)
                },
                onError: setError
            }
        )
    }

    const handleChange = (event) => {
        setFormState({...formState, [event.target.name]: event.target.value});
    };

    return (
        <Form>
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
                    name="name"
                    disabled={isLoading}
                    value={formState.name}
                    onChange={handleChange}
                />
                <TextInputRow
                    label="Interval"
                    name="interval"
                    disabled={isLoading}
                    type="number"
                    value={formState.interval?.toString()}
                    onChange={(event) => setFormState({...formState, interval: parseInt(event.target.value)})}
                />
                <TextInputRow
                    label="Timeout"
                    name="timeout"
                    disabled={isLoading}
                    type="number"
                    value={formState.timeout?.toString()}
                    onChange={(event) => setFormState({...formState, timeout: parseInt(event.target.value)})}
                />
                <SelectRow
                    label="Type"
                    disabled={isLoading}
                    value={formState.type}
                    onChange={(target) => setFormState({...formState, type: target})}
                >
                    <SelectOption key="icmp" value="ICMP" label="ICMP (Unsupported on Akamai)"/>
                    <SelectOption key="http" value="HTTP" label="HTTP"/>
                    <SelectOption key="https" value="HTTPS" label="HTTPS"/>
                    <SelectOption key="tcp" value="TCP" label="TCP"/>
                    <SelectOption key="udp" value="UDP" label="UDP"/>
                </SelectRow>

                {formState.type !== "ICMP" && (
                    <Stack gap="2" distribution="between">
                        <TextareaRow
                            label="Monitor send string"
                            className={"flex-auto"}
                            disabled={isLoading}
                            value={formState?.send}
                            onChange={handleChange}
                        />
                        <TextareaRow
                            label="Monitor expected receive string"
                            className={"flex-auto"}
                            disabled={isLoading}
                            value={formState?.receive}
                            onChange={handleChange}
                        />
                    </Stack>
                )}
            </PanelBody>
        </Form>
    )
}

export default EditMonitorPanel
