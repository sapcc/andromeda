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
} from "@cloudoperators/juno-ui-components";

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

export const ListItemSpinner = ({data, onClick, className, maxLength=25}) => {
    const name = data.name || data.id
    const [showToast, setToast] = useState(false)

    return (
        <Stack alignment="center">
            {["ACTIVE", "DELETED", "ERROR"].includes(data.provisioning_status) && <></> || <Spinner variant={variantClass(data.provisioning_status)} size="small"/>}
            <Stack
                direction="vertical"
                className={`${backgroundClass(data.provisioning_status)}`}
                onClick={onClick}
            >
                <div>
                    {name === data.id && <b>{name}</b> || <small className={className}>{data.id} </small>}
                    {showToast && <Toast text="ID copied to clipboard" className="absolute"/>}
                    &nbsp;
                    <Icon size={name === data.id && "16" || "14"} icon="contentCopy" title="Copy" onClick={() => {
                        copyTextToClipboard(data.id).then(() => {
                            setToast(true)
                            setTimeout(function () { //Start the timer
                                setToast(false)
                            }.bind(this), 1000)
                        })
                    }}/>
                </div>
                {name !== data.id && <p className="jn-font-bold">{name.substring(0, maxLength)}{name.length >= maxLength && "..."}</p> || ""}
            </Stack>
        </Stack>
    )
}

export const ListItemStatus = ({data}) => {
    if (data.provisioning_status === "ACTIVE") {
        if ("status" in data) {
            let icon = "help"
            let variant = "warning"

            switch (data.status) {
                case "ACTIVE":
                case "ONLINE":
                    icon = "wbSunny"
                    variant = "info"
                    break
                case "OFFLINE":
                    icon = "danger"
                    variant = "danger"
                    break
                case "ERROR":
                    icon = "error"
                    variant = "danger"
                    break
                case "NO_MONITOR":
                    return (
                        <Tooltip triggerEvent="hover">
                            <TooltipTrigger><Badge icon="warning" text={data.status} variant="warning" /></TooltipTrigger>
                            <TooltipContent>Member is handed out, but Monitor failed or doesn't exist</TooltipContent>
                        </Tooltip>
                        )
            }
            return <Badge text={data.status} variant={variant} icon={icon} />
        } else {
            return <Badge
                text={data.provisioning_status}
                variant={data.provisioning_status === "ACTIVE" ? "info" : "danger"}
                icon={data.provisioning_status === "ACTIVE" ? "check" : "error"}
            />
        }
    } else {
        let variant = ["PENDING_DELETE", "DELETED", "ERROR"].includes(data.provisioning_status) ? "danger" : "warning"
        return <Badge
            variant={variant}
            text={data.provisioning_status}
            icon={variant}
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
        <Stack className="ml-auto" distribution="end" gap="4">
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
                    href={`https://github.com/sapcc/andromeda/releases`}
                    target="_blank"
                    icon="download"
                    label="CLI Client"
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