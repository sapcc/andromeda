// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useState} from "react"

import {
    Button,
    ButtonRow,
    Form,
    FormRow,
    IntroBox,
    Message,
    Modal,
    Select,
    SelectOption,
    Spinner,
    Stack,
    TextInput
} from "@cloudoperators/juno-ui-components"
import {authStore, urlStore} from "../store"
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {login} from "../actions";
import {Error} from "./Components";

const LogInModal = ({keystoneEndpoint, overrideEndpoint, loginDomains, loginProject, message}) => {
    const setAuth = authStore((state) => state.setAuth)
    const setModal = urlStore((state) => state.openModal)
    const queryClient = useQueryClient()
    const {mutate, error} = useMutation({mutationFn: login})
    const [isLoading, setIsLoading] = useState(false)
    const [showCredentials, setShowCredentials] = useState(false)
    const [mounted, setMounted] = useState(false)
    const [credentials, setCredentials] = useState({
        username: "",
        password: "",
        domain: loginDomains[0] || "",
        project: loginProject,
    })

    const onSubmit = (event) => {
        setIsLoading(true)
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
                },
                onSettled: setIsLoading(false),
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
            {isLoading ? <Spinner variant="primary"/> : null}

            {/* Warning Message */}
            {message && (
                <Message variant="warning" className="mb-4">
                    <Stack>
                        {message}
                    </Stack>
                </Message>
            )}

            {/* Form */}


            <Form onSubmit={onSubmit}>
                <Stack distribution="between" gap="2">
                    <Select
                        helptext="Domain"
                        disabled={isLoading}
                        onChange={(target) => setCredentials({...credentials, domain: target})}
                        value={credentials.domain}
                    >
                        {loginDomains.map(domain =>
                            <SelectOption key={domain} value={domain} label={domain} />
                        )}
                    </Select>
                    <TextInput
                        helptext="Project"
                        name="project"
                        value={credentials.project}
                        disabled={isLoading}
                        onChange={handleChange}
                        autoComplete="on"
                    />
                    <Button
                        className="h-textinput"
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
                <ButtonRow name="Default ButtonRow" className="mt-2 justify-end">
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
