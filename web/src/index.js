import { createRoot } from "react-dom/client"
import React from "react"

var NODE_ENV = process.env.NODE_ENV
if (NODE_ENV === 'development') {
  new EventSource('/esbuild').addEventListener('change', () => location.reload());
}

// export mount and unmount functions
export const mount = (container, options = {}) => {
  import("./App").then((App) => {
    mount.root = createRoot(container)
    mount.root.render(React.createElement(App.default, options?.props))
  })
}

export const unmount = () => mount.root && mount.root.unmount()
