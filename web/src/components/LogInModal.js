import React, {useState} from "react"

import {
    Button,
    ButtonRow,
    Form,
    FormRow,
    IntroBox,
    Modal,
    Select,
    SelectOption,
    Stack,
    TextInput
} from "juno-ui-components"
import {authStore, urlStore} from "../store"
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {login} from "../actions";
import {Error} from "./Components";

const LogInModal = ({keystoneEndpoint, overrideEndpoint}) => {
    const setAuth = authStore((state) => state.setAuth)
    const setModal = urlStore((state) => state.openModal)
    const queryClient = useQueryClient()
    const {mutate, isLoading, error} = useMutation(login)
    const [showCredentials, setShowCredentials] = useState(false)
    const [mounted, setMounted] = useState(false)
    const [credentials, setCredentials] = useState({
        username: undefined,
        password: undefined,
        domain: "monsoon3",
        project: "cc-demo",
    })

    const onSubmit = (event) => {
        event.preventDefault();
        mutate({
                endpoint: keystoneEndpoint,
                ...credentials
            },
            {
                onSuccess: ([token, data]) => {
                    const auth = {
                        token: token,
                        ...data
                    }

                    // set application endpoint
                    if (overrideEndpoint) {
                        auth.endpoint = overrideEndpoint
                    } else {
                        auth.endpoint = data
                            .catalog
                            .find(endpoints => endpoints.type === "gtm")
                            .endpoints
                            .find(endpoint => endpoint.interface === "public")
                            .url
                    }
                    setMounted(false)
                    setTimeout(function() { //Start the timer
                        setAuth(auth)
                        queryClient.invalidateQueries().then()
                    }.bind(this), 700)
                    setModal(null)
                }
            }
        )
    }

    const handleChange = (event) => {
        setCredentials({...credentials, [event.target.name]: event.target.value});
    }

    React.useEffect(() => {
        // used for animation
        setMounted(true)
    }, []);

    return (
        <Modal
            title="Andromeda"
            className={`transition-opacity delay-150 duration-700 ${mounted?"opacity-100":"opacity-0"}`}
            closeable={false}
            open
        >
            {/* Some nice animations */}
            <div className="galaxy">
                <div className="stars" />
            </div>

            <IntroBox
                variant="hero"
                title="Andromeda"
                text="Global Load Balancing as a Service."
            />

            {/* Error Bar */}
            <Error error={error}/>

            <Form onSubmit={onSubmit}>
                <Stack distribution="between" gap="2" className="pt-2">
                    <Select
                        label="Domain"
                        disabled={isLoading}
                        onChange={(target) => setCredentials({...credentials, domain: target})}
                        value={credentials.domain}
                    >
                        <SelectOption key="monsoon3" value="monsoon3" label="monsoon3" />
                        <SelectOption key="ccadmin" value="ccadmin" label="ccadmin" />
                    </Select>
                    <TextInput
                        label="Project"
                        name="project"
                        value={credentials.project}
                        disabled={isLoading}
                        onChange={handleChange}
                    />
                    <Button
                        className="flex-none jn-relative jn-mb-2"
                        icon="manageAccounts"
                        onClick={() => setShowCredentials(!showCredentials)}
                    />
                </Stack>
                {showCredentials && (
                    <div>
                        <FormRow>
                            <TextInput
                                label="User Name"
                                name="username"
                                value={credentials.username}
                                disabled={isLoading}
                                onChange={handleChange}
                                required
                            />
                        </FormRow>
                        <FormRow>
                            <TextInput
                                label="Password"
                                name="password"
                                type="password"
                                value={credentials.password}
                                disabled={isLoading}
                                onChange={handleChange}
                                required
                            />
                        </FormRow>
                    </div>
            )}
                <ButtonRow name="Default ButtonRow" className="jn-justify-end">
                    <Button
                        icon="openInBrowser"
                        label={`Enter ${credentials.domain}`}
                        variant="primary"
                        type="submit"
                        disabled={isLoading || !mounted}
                        progress={isLoading || !mounted}
                        onClick={onSubmit}
                    />
                </ButtonRow>
            </Form>
        </Modal>
    )
}

export default LogInModal
