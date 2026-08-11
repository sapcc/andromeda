<!--
SPDX-FileCopyrightText: Copyright 2022-2025 SAP SE or an SAP affiliate company and andromeda contributors

SPDX-License-Identifier: Apache-2.0
-->

# `andromeda-akamai-agent` app

This worker will continuously run an Andromeda / Akamai reconciliation loop. The Andromeda database is considered the "source of truth", meaning that any configuration drift caused by manual updates to Akamai API resources controlled by respective Andromeda resources will be overwritten by the reconciliation loop that next time it runs.

## Manually testing the reconciliation loop

Pre-requisites for testing each flow:
* Ensure the `akamai.yaml` file points to the desired Akamai domain (`akamai.domain` property).
* Ensure that the local database is in a fresh state: run `andromeda-migrate` on an empty Andromeda database.

### Flow: creating a property and server

On shell 1, start the Andromeda API Server:

```sh
go run cmd/andromemeda-server/main.go --config-file ./tmp/andromeda.yaml --config-file ./tmp/akamai.yaml
```

On shell 2, ensure the appropriate `OS_` environment variables are in place, then perform the following operations in turn, adjusting values to your actual needs, if needed.

```sh
m31ctl datacenter create --name eu-de-1 --provider akamai
datacenter_id="$(m31ctl datacenter list -f json | jq -r '.[] | select(.name == "eu-de-1") | .id')"

m31ctl domain create --name my-test-domain.de --fqdn my-test-domain.de --provider akamai
domain_id="$(m31ctl domain list -f json | jq -r '.[] | select(.name == "my-test-domain.de") | .id')"

m31ctl pool create --name my-test-domain.de-pool-1 --domain "$domain_id"
pool_id="$(m31ctl pool list -f json | jq -r '.[] | select(.name == "my-test-domain.de-pool-1") | .id')"

m31ctl member create --name member-1 --address '200.100.0.1' --port 8080 --datacenter "$datacenter_id" "$pool_id"

# Start the Andromeda / Akamai reconciliation loop
PROMETHEUS_LISTEN=0.0.0.0:9091 go run cmd/andromemeda-akamai-agent/main.go \
    --config-file ./tmp/andromeda.yaml \
    --config-file ./tmp/akamai.yaml
```

If the reconciliation loop succeeds, you should be able to navigate to the domain page and see that `my-test-domain.de` property was added under the domain configured in your local `akamai.yaml` file. The property should have one server with IP `200.100.0.1`.

### Flow: deleting a property (by deleting all its servers)

When the reconciliation loop detects that a property in the local database has no associated servers, it deletes it from Akamai. For testing this flow, please first follow the instructions for flow [creating a property and server](#creating-a-property-and-server).

On shell 1, start the Andromeda API Server:

```sh
go run cmd/andromemeda-server/main.go --config-file ./tmp/andromeda.yaml --config-file ./tmp/akamai.yaml
```

On shell 2, ensure the appropriate `OS_` environment variables are in place, then:

```sh
member_id="$(m31ctl member list -f json | jq -r '.[] | .id')"
m31ctl member delete "$member_id"

# Start the Andromeda / Akamai reconciliation loop
PROMETHEUS_LISTEN=0.0.0.0:9091 go run cmd/andromemeda-akamai-agent/main.go \
    --config-file ./tmp/andromeda.yaml \
    --config-file ./tmp/akamai.yaml
```

If the reconciliation loop succeeds, you should be able to navigate to the domain page and see that `my-test-domain.de` property under the domain configured in your local `akamai.yaml` file is no longer there.
