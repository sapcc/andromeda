// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import {createRoot} from "react-dom/client"
import React from "react"
import _ from "lodash-es"

var NODE_ENV = process.env.NODE_ENV
if (NODE_ENV === 'development') {
  new EventSource('/esbuild').addEventListener('change', () => location.reload());
}

let APP_PROPS = process.env.APP_PROPS;
if (APP_PROPS) {
    try { APP_PROPS = JSON.parse(APP_PROPS); } catch (e) { }
}

// export mount and unmount functions
export const mount = (container, options = {props: {}}) => {
  import("./App").then((App) => {
    if (APP_PROPS) { _.extend(options.props, APP_PROPS); }
    console.log(options.props);

    mount.root = createRoot(container)
    mount.root.render(React.createElement(App.default, options.props))
  })
}

export const unmount = () => mount.root && mount.root.unmount()
