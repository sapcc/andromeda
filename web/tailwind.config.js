// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

// This file is used to configure the tailwindcss library.
module.exports = {
  presets: [
    require("@cloudoperators/juno-ui-components/build/lib/tailwind.config"), // important, do not change
  ],
  prefix: "", // important, do not change
  content: ["./src/**/*.{js,jsx,ts,tsx}", "./public/index.html"],
  corePlugins: {
    preflight: false, // important, do not change
  },
  theme: {},
  plugins: [],
}
