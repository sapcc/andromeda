# Total DNS Requests CLI Tool

This tool calculates the total DNS requests for an Akamai GTM property over a specified time period. It provides detailed information about request distribution across datacenters.

## Description

The Total DNS Requests tool helps in monitoring and analyzing the DNS request traffic across different Akamai datacenters. This is useful for:

- Understanding traffic distribution patterns
- Identifying potential issues with specific datacenters
- Capacity planning and performance optimization
- Troubleshooting DNS-related problems

## Usage

```bash
andromeda-akamai-total-dns-requests --property=<property_id> [--domain=<domain>] [--start=<start_date>] [--end=<end_date>] [--output=<format>]
```

### Arguments

- `--property` - Name of the GTM property to analyze (required)
- `--domain` - Name of the GTM domain (default: andromeda.akadns.net)
- `--start` - Start date in RFC3339 format (default: 2 days before end date)
- `--end` - End date in RFC3339 format (default: 15 minutes ago)
- `--output` - Output format: json, csv, or text (default: text)

### Examples

```bash
# Get total DNS requests for a property with default settings
andromeda-akamai-total-dns-requests --property=andromeda-mena.ccee.sapcloud.io

# Get total DNS requests with a specific date range
andromeda-akamai-total-dns-requests --property=andromeda-mena.ccee.sapcloud.io --start=2025-03-17T00:00:00Z --end=2025-03-19T00:00:00Z

# Get total DNS requests in JSON format
andromeda-akamai-total-dns-requests --property=andromeda-mena.ccee.sapcloud.io --output=json

# Get total DNS requests in CSV format
andromeda-akamai-total-dns-requests --property=andromeda-mena.ccee.sapcloud.io --output=csv
```

## Output Formats

### Text (default)

```
Total DNS Requests for Property: andromeda-mena.ccee.sapcloud.io
Time Range: 2025-03-17T00:00:00Z to 2025-03-19T00:00:00Z
Total Requests: 16503
----------------------------------------------------------------------------------------------------
Datacenter                     | ID              | Traffic Target                | Requests        | Percentage
----------------------------------------------------------------------------------------------------
91c259c4-f838-11ef-8688-0af... | 11              | 91c259c4-f838-11ef-8688-0af... | 8223            | 49.83%
76484d18-f838-11ef-8688-0af... | 9               | 76484d18-f838-11ef-8688-0af... | 8280            | 50.17%
```

### JSON

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

### CSV

```csv
Property,Start Date,End Date,Total Requests
andromeda-mena.ccee.sapcloud.io,2025-03-17T00:00:00Z,2025-03-19T00:00:00Z,16503

Datacenter,Datacenter ID,Traffic Target,Requests,Percentage
91c259c4-f838-11ef-8688-0af1ec3a88ef,11,91c259c4-f838-11ef-8688-0af1ec3a88ef - 12.13.14.15:443,8223,49.83%
76484d18-f838-11ef-8688-0af1ec3a88ef,9,76484d18-f838-11ef-8688-0af1ec3a88ef - 23.24.25.26:443,8280,50.17%
```

## Building

To build just this tool, run:

```bash
make build-akamai-dns-requests
```

Alternatively, build all Andromeda components with:

```bash
make build-all
``` 