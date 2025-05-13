<!--
SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

# Metrics Commands

## Akamai DNS Requests

The `metrics akamai-dns-requests` command is used to retrieve and display DNS request metrics from Akamai GTM properties for a specific domain.

### Usage

```
m31ctl metrics akamai-dns-requests --domain-id DOMAIN_UUID [--start START_TIME] 
  [--end END_TIME] [--output FORMAT]
```

### Arguments

- `--domain-id` - Andromeda Domain UUID (required)
- `--start, -s` - Start date/time in RFC3339 format (e.g., 2025-05-01T00:00:00Z)
- `--end, -e` - End date/time in RFC3339 format (e.g., 2025-05-10T15:00:00Z)
- `--output, -o` - Output format: json, csv, or table (default: table)

### Examples

```bash
# Get total DNS requests for a domain with default settings
m31ctl metrics akamai-dns-requests --domain-id 12345678-1234-1234-1234-123456789012

# Get total DNS requests with a specific date range
m31ctl metrics akamai-dns-requests --domain-id 12345678-1234-1234-1234-123456789012 \
  --start 2025-03-17T00:00:00Z --end 2025-03-19T00:00:00Z

# Get total DNS requests in JSON format
m31ctl metrics akamai-dns-requests --domain-id 12345678-1234-1234-1234-123456789012 --output json

# Get total DNS requests in CSV format
m31ctl metrics akamai-dns-requests --domain-id 12345678-1234-1234-1234-123456789012 --output csv
```

### Output Examples

#### Table Format (default)

```
Total DNS Requests for Domain: 12345678-1234-1234-1234-123456789012
Time Range: 2025-03-17T00:00:00Z to 2025-03-19T00:00:00Z
Total Requests: 16503

+--------------------------------------+---------------+--------------------------------------+----------+-----------+
| DATACENTER                           | DATACENTER ID | TRAFFIC TARGET                       | REQUESTS | PERCENTAGE |
+--------------------------------------+---------------+--------------------------------------+----------+-----------+
| 91c259c4-f838-11ef-8688-0af1ec3a88ef | 11            | 91c259c4-f838-11ef-8688-0af1ec3a88ef | 8223     | 49.83%    |
| 76484d18-f838-11ef-8688-0af1ec3a88ef | 9             | 76484d18-f838-11ef-8688-0af1ec3a88ef | 8280     | 50.17%    |
+--------------------------------------+---------------+--------------------------------------+----------+-----------+
```

#### JSON Format

```json
{
  "property": "andromeda-mena.ccee.sapcloud.io",
  "start_date": "2025-03-17T00:00:00Z",
  "end_date": "2025-03-19T00:00:00Z",
  "total_requests": 16503,
  "datacenters": {
    "91c259c4-f838-11ef-8688-0af1ec3a88ef": {
      "datacenter_id": "11",
      "percentage": 0.4983,
      "total_requests": 8223,
      "traffic_target": "91c259c4-f838-11ef-8688-0af1ec3a88ef - 12.13.14.15:443"
    },
    "76484d18-f838-11ef-8688-0af1ec3a88ef": {
      "datacenter_id": "9",
      "percentage": 0.5017,
      "total_requests": 8280,
      "traffic_target": "76484d18-f838-11ef-8688-0af1ec3a88ef - 23.24.25.26:443"
    }
  }
}
```

#### CSV Format

```csv
Property,Start Date,End Date,Total Requests
andromeda-mena.ccee.sapcloud.io,2025-03-17T00:00:00Z,2025-03-19T00:00:00Z,16503

Datacenter,Datacenter ID,Traffic Target,Requests,Percentage
91c259c4-f838-11ef-8688-0af1ec3a88ef,11,91c259c4-f838-11ef-8688-0af1ec3a88ef - 12.13.14.15:443,8223,49.83%
76484d18-f838-11ef-8688-0af1ec3a88ef,9,76484d18-f838-11ef-8688-0af1ec3a88ef - 23.24.25.26:443,8280,50.17%
``` 