<!--
SPDX-FileCopyrightText: Copyright 2022-2025 SAP SE or an SAP affiliate company and andromeda contributors

SPDX-License-Identifier: Apache-2.0
-->

# `andromeda-f5-status-agent` app

Note: this documentation refers to the current "MVP" implementation of F5 status agent.

## Manually testing the main flow

On shell 1, start both the Andromeda API Server and the RPC server:

```sh
go run cmd/andromemeda-server/main.go --config-file ./tmp/andromeda.yaml
```

On shell 2, run the Andromeda / F5 status reconciliation loop:

```sh
# Start the Andromeda / F5 status reconciliation loop:
PROMETHEUS_LISTEN=0.0.0.9091 go run cmd/andromeda-f5-status-agent/main.go \
    --config-file ./tmp/andromeda.yaml \
    --config-file ./tmp/f5.yaml
```
