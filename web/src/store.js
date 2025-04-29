// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import { createWithEqualityFn as create } from 'zustand/traditional'
import { createJSONStorage, persist } from "zustand/middleware"
import { querystring } from "zustand-querystring";
import { mountStoreDevtool } from 'simple-zustand-devtools';

export const authStore = create(
    persist((set, get) => ({
            auth: undefined,
            setAuth: (newAuth) => set(() => ({auth: newAuth})),
        }),
        {
            name: 'auth',
            storage: createJSONStorage(() => sessionStorage),
        }
    )
)

export const urlStore = create(
    querystring((set) => ({
            t: 0,         // tab index
            m: null,      // modal
            p: null,      // panel
            id: null,     // id
            pool: null,   // pool
            domain: null, // domain
            openModal: (newModal) => set(() => ({m: newModal})),
            openModalWithId: (newModal, newId) => set(() => ({m: newModal, id: newId})),
            closeModal: () => set(() => ({m: null})),
            openPanel: (newPanel, newId) => set(() => ({p: newPanel, id: newId})),
            setTab: (newTab) => set(() => ({t: newTab})),
            setPool: (newPool) => set(() => ({pool: newPool})),
            clearPool: () => set(() => ({pool: null})),
            setDomain: (newDomain) => set(() => ({domain: newDomain})),
        })
    )
)

if (process.env.NODE_ENV === 'development') {
    mountStoreDevtool('AutoStore', authStore);
    mountStoreDevtool('URLStore', urlStore);
}
