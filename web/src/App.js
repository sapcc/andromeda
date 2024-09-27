import React, {useState} from "react"

import "./styles.css"
import {authStore} from "./store"
import {AppShell, ContentHeading, PageHeader, PortalProvider, StyleProvider} from "@cloudoperators/juno-ui-components"
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

    const pageHeaderStyles = `jn-min-h-[3.25rem] jn-bg-juno-grey-blue-11 jn-sticky jn-top-0 jn-px-6 jn-py-3 jn-z-50`
    const pageHeaderInnerStyles = `jn-grid jn-grid-cols-[minmax(0,max-content),1fr] jn-gap-3 jn-h-8 jn-w-full jn-overflow-hidden jn-items-center`
    const headingStyles =  `jn-text-lg jn-text-theme-high jn-cursor-pointer`
    const pageHeader = (
        <div className={`juno-pageheader theme-dark ${pageHeaderStyles}`} role="banner">
            <div className={`juno-pageheader-inner ${pageHeaderInnerStyles}`}>
                <div>
                    <div className={headingStyles} onClick={() => {window.location.href = '/'}}>
                        Converged Cloud | Andromeda
                    </div>
                </div>
                {auth && (
                    <HeaderUser auth={auth} logout={logout} theme={theme} setTheme={setTheme}/>
                )}
            </div>
        </div>
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
                        <ContentHeading heading="Global Load Balancing as a Service" className="jn-p-2"/>
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
