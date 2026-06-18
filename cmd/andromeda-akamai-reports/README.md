<!--
SPDX-FileCopyrightText: Copyright 2022-2025 SAP SE or an SAP affiliate company and andromeda contributors

SPDX-License-Identifier: Apache-2.0
-->

# `andromeda-akamai-reports` app

This worker will bind a Prometheus HTTP server to `0.0.0.0:9090` (or the host and port provided by `PROMETHEUS_LISTEN`) and register a collector (`AndromedaAkamaiCollector`) that implements the `prometheus.Collector` interface. Whenever its `/metrics` endpoint is requested, Prometheus invokes `Collect()`, which in turn:

* Iterates over each property using the Akamai GTM API
* For each property...
  * Queries the Akamai GTM API for its [traffic report][traffic-report]
  * Retrieves its respective Keystone project ID from the Andromeda database
  * Writes the results of the API call to two gauge metrics:
    * `andromeda_akamai_requests_5m`
    * `andromeda_akamai_status_5m`

[traffic-report]: <https://techdocs.akamai.com/gtm-reporting/reference/get-traffic-property>

The Akamai metrics will be exposed at <http://localhost:9091/metrics>.

## Manually testing the main flow

Ensure the `akamai.yaml` file points to the desired Akamai domain (`akamai.domain` property).

Ensure that, for each property/domain/server combination that you'd like to test for, there's one respective `datacenter` table record in place. Navigate to the Akamai property details page and copy the the server UUID you'd like to test. Then, create a record in the local Andromeda `datacenter` table by providing some value (a real Keystone project ID is not necessary) for `datacenter.project_id` and  the Akamai property server UUID for `datacenter.id`:

```sql
INSERT INTO datacenter (id, project_id, provider, admin_state_up) VALUES
    ("...", "...", "akamai", 1);
```

On shell 1, start the RPC server.

```sh
go run cmd/andromemeda-server/main.go \
    --config-file ./tmp/andromeda.yaml \
    --config-file ./tmp/akamai.yaml
```

On Shell 2, start the Akamai Reports app:

```sh
export PROMETHEUS_LISTEN=0.0.0.0:9091
go run cmd/andromemeda-akamai-reports/main.go \
    --config-file ./tmp/andromeda.yaml \
    --config-file ./tmp/akamai.yaml
```

In order for Prometheus to invoke the Andromeda Akamai Collector, either browse or curl the `/metrics` endpoint:

```sh
curl http://0.0.0.0:9091/metrics | grep '_5m'
```

The output should include data for both `andromeda_akamai_requests_5m` and `andromeda_akamai_status_5m` custom metrics.
