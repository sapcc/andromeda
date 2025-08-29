<!--
SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

# Andromeda Web GUI

This is the [Juno](https://github.com/sapcc/juno) based Web GUI. 

![](https://github.com/sapcc/andromeda/blob/main/web/screenshot.png?raw=true)

## Development

Install dependencies by running

```sh
npm install
```

### Running dev server

```sh
export APP_PORT=8000; npm run start
```

### Build static minified release

```sh
npm run build
```

### Keystone Server requirements

The Keystone API Server must allow CORS (Cross-Origin Resource Sharing) on `/v3/auth/tokens` so the Andromeda Web GUI can create a Keystone token. The following is a sample minimal CORS configuration for a local Keystone API Server fronted by the Apache 2 HTTP Server:

```
<Location /v3/auth/tokens>
    Header set Access-Control-Allow-Origin "http://localhost:8000"
    Header set Access-Control-Allow-Methods "POST"
    Header set Access-Control-Allow-Headers "Content-Type, X-User-Domain-Name"
    Header set Access-Control-Expose-Headers "X-Subject-Token"
</Location>
```