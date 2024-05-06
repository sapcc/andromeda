import React, {useMemo, useState} from "react";

import {
    Badge,
    Button,
    Icon,
    JsonViewer,
    Message,
    Modal,
    Spinner,
    Stack,
    Toast,
    Tooltip,
    TooltipContent,
    TooltipTrigger
} from "juno-ui-components";

export const Error = ({error}) => {
    if (error) {
        return (
            <Message variant="error" className="jn-mb-4">
                {error.statusCode && `[${error.statusCode}] `}
                {error.message}
            </Message>
        )
    }
}

export const Loading = ({isLoading}) => {
    if (isLoading) {
        return (
            <Message variant="info" className="jn-mb-4">
                <Stack>
                    Loading...
                </Stack>
            </Message>
        )
    }
}

const variantClass = (variant) => {
    switch (variant) {
        case "PENDING_CREATE":
            return "primary"
        case "PENDING_DELETE":
            return "danger"
        case "PENDING_UPDATE":
            return "warning"
    }
}
const backgroundClass = (variant) => {
    switch (variant) {
        case "PENDING_CREATE":
            return "animate-pulse"
        case "PENDING_DELETE":
            return "animate-pulse"
        case "PENDING_UPDATE":
            return "animate-pulse"
    }
}

export async function copyTextToClipboard(text) {
    if ('clipboard' in navigator) {
        return await navigator.clipboard.writeText(text);
    } else {
        return document.execCommand('copy', true, text);
    }
}

export const ListItemSpinner = ({data, onClick, className, maxLength=15}) => {
    const name = data.name || data.id
    const [showToast, setToast] = useState(false)

    return (
        <Stack
            alignment="center"
            className={`${className} ${backgroundClass(data.provisioning_status)} jn-font-bold`}
            onClick={onClick}
            gap="1.5"
        >
            {["ACTIVE", "DELETED", "ERROR"].includes(data.provisioning_status) || <Spinner
                variant={variantClass(data.provisioning_status)} size="small"/>}
                <Tooltip triggerEvent="hover">
                    <TooltipTrigger asChild>
                        {name.substring(0, maxLength)}{name > maxLength && "..."}
                    </TooltipTrigger>
                    <TooltipContent>{data.id}</TooltipContent>
                </Tooltip>
            <div>
                {showToast && <Toast text="ID copied to clipboard" className="absolute"/>}
                <Icon size="16" icon="contentCopy" onClick={() => {
                    copyTextToClipboard(data.id).then(() => {
                        setToast(true)
                    setTimeout(function () { //Start the timer
                        setToast(false)
                    }.bind(this), 1000)
                })
                }} />
            </div>
        </Stack>
    )
}

export const ListItemStatus = ({data}) => {
    if (data.provisioning_status === "ACTIVE") {
        if ("status" in data) {
            return <Badge
                text={data.status}
                variant={data.status === "ONLINE" ? "info" : (data.status === "UNKNOWN" ? "warning" : "danger") }
            />
        } else {
            return <Badge
                text={data.provisioning_status}
                variant={data.provisioning_status === "ACTIVE" ? "info" : "danger"}
            />
        }
    } else {
        return <Badge
            variant={["PENDING_DELETE", "DELETED", "ERROR"].includes(
                data.provisioning_status) ? "danger" : "warning"}
                text={data.provisioning_status}
            />
    }
}

const avatarCss = `
h-8
w-8
bg-theme-background-lvl-2
rounded-full
mr-2
bg-cover 
`

export const HeaderUser = ({ auth, logout, theme, setTheme }) => {
    const sapID = useMemo(() => auth?.user.name || "", [auth])

    return (
        <Stack alignment="center" className="ml-auto" distribution="end" gap="4">
            <Tooltip triggerEvent="hover">
                <TooltipTrigger>
                    <Stack alignment="center">
                        <div
                            style={{
                                background: `url(https://avatars.wdf.sap.corp/avatar/${sapID}?size=24x24) no-repeat`,
                                backgroundSize: `cover`,
                            }}
                            className={avatarCss}
                        />
                        {<span>{sapID}</span>}
                    </Stack>
                </TooltipTrigger>
                <TooltipContent>
                    {auth.project.name}@{auth.project.domain.name}
                </TooltipContent>
            </Tooltip>

            <Stack gap="2">
                <Button
                    onClick={() => {setTheme(theme === "theme-light" ? "theme-dark" : "theme-light")
                        console.log(`${theme} => setTheme(${theme === "theme-light" ? "theme-dark" : "theme-light"})`)
                    }}
                    icon="danger"
                    variant="subdued"
                    size="small"
                />

                <Button
                    href={`${auth.endpoint}/docs`}
                    target="_blank"
                    icon="openInNew"
                    label="API"
                    variant="subdued"
                    size="small"
                />

                <Button
                    label="Logout"
                    variant="primary-danger"
                    icon="exitToApp"
                    size="small"
                    onClick={logout}
                />
            </Stack>
        </Stack>
    )
}

export const JsonModal = (data) => {
    return (
        <Modal size="large" open >
            <JsonViewer data={data} toolbar />
        </Modal>
    )
}