import React, {useState} from "react"

import "./styles.css"
import {authStore} from "./store"
import {AppShell, ContentHeading, PageHeader, PortalProvider, StyleProvider} from "juno-ui-components"
import {QueryCache, QueryClient, QueryClientProvider} from '@tanstack/react-query'
import AppContent from "./AppContent"
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";
import LogInModal from "./components/LogInModal";
import {HeaderUser} from "./components/Components";

const App = (props) => {
    const [theme, setTheme] = useState(props.theme)
    const [auth, setAuth] = authStore((state) => [state.auth, state.setAuth])

    const logout = () => {
        setAuth(undefined)
        queryClient.invalidateQueries().then()
    }

    const pageHeader = (
        <PageHeader heading="Converged Cloud | Andromeda" onClick={() => window.location.href = '/'}>
            {auth && (
                <HeaderUser auth={auth} logout={logout} theme={theme} setTheme={setTheme}/>
            )}
        </PageHeader>
    )

    // Create query client which it can be used from overall in the app
    const queryClient = new QueryClient({
        queryCache: new QueryCache({
            onError: (error) => {
                if (error?.statusCode === 401) {
                    // force re-authenticate
                    logout()
                }
            }
        }),
    })

    return (
        <StyleProvider stylesWrapper="head" theme={theme} key={theme}>
            <PortalProvider>
                <QueryClientProvider client={queryClient}>
                    <AppShell
                        pageHeader={pageHeader}
                        embedded={props.embedded === true}
                    >
                        <ContentHeading>Global Load Balancing as a Service</ContentHeading>
                        {auth ? (
                            <AppContent props={props}/>
                        ) : (
                            <LogInModal keystoneEndpoint={props.endpoint}
                                        overrideEndpoint={props.overrideAndromedaEndpoint}
                                        loginDomains={props?.loginDomains || []}
                                        loginProject={props?.loginProject}
                            />
                        )}
                    </AppShell>
                    <ReactQueryDevtools initialIsOpen={false}/>
                </QueryClientProvider>
            </PortalProvider>
        </StyleProvider>
    )
}

export default App
