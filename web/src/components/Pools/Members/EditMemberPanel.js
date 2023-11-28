import React, {useEffect, useState} from "react"

import {Button, Checkbox, Form, PanelBody, PanelFooter, Stack, TextInput,} from "juno-ui-components"
import {authStore, urlStore} from "../../../store"
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {fetchItem, updateAttributes, updateItem} from "../../../actions"
import DatacenterMenu from "./DatacenterMenu";
import {Error, Loading} from "../../Components";

const EditMemberPanel = ({closeCallback}) => {
    const auth = authStore((state) => state.auth)
    const id = urlStore((state) => state.id)
    const [error, setError] = useState()
    const [formState, setFormState] = useState({
        name: undefined,
        address: undefined,
        datacenter_id: undefined,
        port: undefined,
        admin_state_up: undefined,
    })

    const queryClient = useQueryClient()
    const {data, isLoading} = useQuery({
        queryKey: ["members", id],
        queryFn: fetchItem,
        meta: auth,
        refetchOnWindowFocus: false,
    })
    const mutation = useMutation({mutationFn: updateItem})

    // update formState when data is fetched
    useEffect(() => {
        if (data) {
            setFormState(updateAttributes(formState, data.member))
        }
    }, [data]);

    const onSubmit = () => {
        mutation.mutate(
            {
                key: "members",
                id: id,
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
                <Error error={error}/>

                {/* Loading indicator for page content */}
                <Loading isLoading={isLoading || mutation.isLoading}/>

                <Stack direction="vertical" gap="2">
                    <Checkbox
                        id="selectable"
                        label="Enabled"
                        checked={formState.admin_state_up}
                        disabled={isLoading}
                        onChange={(event) => setFormState({...formState, admin_state_up: event.target.checked})}
                    />
                    <TextInput
                        label="Name"
                        name="name"
                        value={formState.name}
                        disabled={isLoading}
                        onChange={handleChange}
                    />
                    <TextInput
                        label="IP address"
                        name="address"
                        value={formState.address}
                        disabled={isLoading}
                        onChange={handleChange}
                        required
                    />
                    <TextInput
                        label="Port"
                        type="number"
                        disabled={isLoading}
                        value={formState.port?.toString()}
                        onChange={(event) => setFormState({...formState, port: parseInt(event.target.value)})}
                        required
                    />
                    Select a Datacenter
                    <DatacenterMenu formState={formState} setFormState={setFormState} setError={setError}/>
                </Stack>
            </PanelBody>
        </Form>
    )
}

export default EditMemberPanel
