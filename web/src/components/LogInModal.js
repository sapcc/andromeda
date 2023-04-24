import React, {useState} from "react"

import {
    Button,
    ButtonRow,
    Form,
    IntroBox,
    LoadingIndicator,
    Modal,
    SelectOption,
    SelectRow,
    Stack,
    TextInputRow
} from "juno-ui-components"
import {authStore, useStore} from "../store"
import {currentState, push} from "url-state-provider";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {login} from "../actions";
import heroImage from "../img/app_bg_example.svg?url"
import {Error, Loading} from "./Components";

const LogInModal = ({keystoneEndpoint, overrideEndpoint}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const setAuth = authStore((state) => state.setAuth)
    const queryClient = useQueryClient()
    const {mutate, isLoading, error} = useMutation(login)
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
                credentials
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
                    setAuth(auth)
                    queryClient.invalidateQueries().then()
                    const urlState = currentState(urlStateKey)
                    push(urlStateKey, {...urlState, currentModal: ""})
                }
            }
        )
    }

    const handleChange = (event) => {
        setCredentials({...credentials, [event.target.name]: event.target.value});
    };

    return (
        <Modal
            title="Andromeda"
            closeable={false}
            open
        >
            <IntroBox
                heroImage={`url(${heroImage})`}
                text="Log in using your OpenStack Credentials."
                variant="hero"
            />

            {/* Error Bar */}
            <Error error={error}/>

            {/* Loading indicator for page content */}
            <Loading isLoading={isLoading} />

            <Form onSubmit={onSubmit}>
                <Stack distribution="between">
                    <SelectRow
                        label="Domain"
                        name="domain"
                        disabled={isLoading}
                        onChange={handleChange}
                        value={credentials.domain}
                    >
                        <SelectOption key="monsoon3" value="monsoon3" label="monsoon3" />
                        <SelectOption key="ccadmin" value="ccadmin" label="ccadmin" />
                    </SelectRow>
                    <TextInputRow
                        label="Project"
                        name="project"
                        value={credentials.project}
                        disabled={isLoading}
                        onChange={handleChange}
                        required
                    />
                </Stack>
                <TextInputRow
                    label="User Name"
                    name="username"
                    value={credentials.username}
                    disabled={isLoading}
                    onChange={handleChange}
                    required
                />
                <TextInputRow
                    label="Password"
                    name="password"
                    type="password"
                    value={credentials.password}
                    disabled={isLoading}
                    onChange={handleChange}
                    required
                />
                <div className="jn-py-2">
                    <ButtonRow name="Default ButtonRow" className="jn-justify-end">
                        <Button
                            label="Connect"
                            title="Connect"
                            variant="primary"
                            type="submit"
                            onClick={onSubmit}
                        />
                    </ButtonRow>
                </div>
            </Form>
        </Modal>
    )
}

export default LogInModal
