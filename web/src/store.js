import React from "react"
import {create} from 'zustand'
import {createJSONStorage, devtools, persist} from "zustand/middleware"

// global zustand store. See how this works here: https://github.com/pmndrs/zustand
export const useStore = create(
    devtools((set) => ({
        // basics
        urlStateKey: undefined,
        setUrlStateKey: (newUrlStateKey) => set((state) => ({urlStateKey: newUrlStateKey})),
    }))
)

export const authStore = create(
    devtools(
        persist(
            (set, get) => ({
                auth: undefined,
                setAuth: (newAuth) => set(() => ({auth: newAuth})),
            }),
            {
                name: 'auth',
                storage: createJSONStorage(() => sessionStorage),
            }
        )
    )
)
