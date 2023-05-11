import React, {useState} from "react"

import "./styles.css"
import {authStore, useStore} from "./store"
import {AppShell, PageHeader, PortalProvider, StyleProvider} from "juno-ui-components"
import {QueryCache, QueryClient, QueryClientProvider} from '@tanstack/react-query'
import AppContent from "./AppContent"
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";
import LogInModal from "./components/LogInModal";
import {HeaderUser} from "./components/Components";

const URL_STATE_KEY = "andromeda"

const App = (props) => {
    const [theme, setTheme] = useState(props.theme)
    const setUrlStateKey = useStore((state) => state.setUrlStateKey)
    const [auth, setAuth] = authStore((state) => [state.auth, state.setAuth])

    // on app initial load save Endpoint and URL_STATE_KEY so it can be
    // used from overall in the application
    React.useEffect(() => {
        setUrlStateKey(URL_STATE_KEY)
    }, [])

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
                        contentHeading="Global Load Balancing as a Service"
                        embedded={props.embedded === true}
                    >
                        {auth ? (
                            <AppContent props={props}/>
                        ) : (
                            <LogInModal keystoneEndpoint={props.endpoint}
                                        overrideEndpoint={props.overrideAndromedaEndpoint}/>
                        )}
                    </AppShell>
                    <ReactQueryDevtools initialIsOpen={false}/>
                </QueryClientProvider>
            </PortalProvider>
        </StyleProvider>
    )
}

export default App
