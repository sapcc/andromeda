# Metrics Commands

## Akamai DNS Requests

The `metrics akamai-dns-requests` command is used to retrieve and display DNS request metrics from Akamai GTM properties.

### Usage

```
m31ctl metrics akamai-dns-requests --property PROPERTY_NAME [--domain DOMAIN_NAME] 
  [--start START_TIME] [--end END_TIME] [--output FORMAT] [--edgerc EDGERC_PATH]
```

### Arguments

- `--property, -p` - Name of the GTM property to analyze (required)
- `--domain, -d` - Name of the GTM domain (default: andromeda.akadns.net)
- `--start, -s` - Start date in RFC3339 format (default: 2 days before end date)
- `--end, -e` - End date in RFC3339 format (default: 15 minutes ago)
- `--output, -o` - Output format: json, csv, or table (default: table)
- `--edgerc` - Path to .edgerc file for Akamai authentication (optional, will look in common locations by default)

### Authentication

This command requires an Akamai .edgerc file for authentication. By default, it will look for this file in:

1. Current directory (.edgerc)
2. User's home directory (~/.edgerc)
3. Default location from configuration

You can also specify a custom location using the `--edgerc` flag:

```bash
m31ctl metrics akamai-dns-requests --property PROPERTY --edgerc /path/to/.edgerc
```

The .edgerc file should contain your Akamai API credentials in this format:

```
[default]
client_token = your-client-token
client_secret = your-client-secret
access_token = your-access-token
host = your-akamai-host
```

### Examples

```bash
# Get total DNS requests for a property with default settings
m31ctl metrics akamai-dns-requests --property andromeda-mena.ccee.sapcloud.io

# Get total DNS requests with a specific date range
m31ctl metrics akamai-dns-requests --property andromeda-mena.ccee.sapcloud.io \
  --start 2025-03-17T00:00:00Z --end 2025-03-19T00:00:00Z

# Get total DNS requests in JSON format
m31ctl metrics akamai-dns-requests --property andromeda-mena.ccee.sapcloud.io --output json

# Get total DNS requests in CSV format
m31ctl metrics akamai-dns-requests --property andromeda-mena.ccee.sapcloud.io --output csv

# Get total DNS requests with a specific .edgerc file
m31ctl metrics akamai-dns-requests --property andromeda-mena.ccee.sapcloud.io \
  --edgerc /path/to/custom/.edgerc
```

### Output Examples

#### Table Format (default)

```
Total DNS Requests for Property: andromeda-mena.ccee.sapcloud.io
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