


# Andromeda
Platform agnostic GSLB frontend
  

## Informations

### Version

1.0.0

## Tags

  ### <span id="tag-domains"></span>Domains

Domains are the highest level of entities defining DNS endpoints to be load balanced.

  ### <span id="tag-pools"></span>Pools

Pools are logical collections of datacenters hosting target applications.

  ### <span id="tag-datacenters"></span>Datacenters

Datacenter are collections of Members that share the same geographical location.

  ### <span id="tag-members"></span>Members

Members are IP/Port endpoints of the applications to be load balanced.

  ### <span id="tag-monitors"></span>Monitors

Monitors are health checks that influce load balancing decisions.

  ### <span id="tag-administrative"></span>Administrative

Administrative API

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  administrative

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/quotas/{project_id} | [delete quotas project ID](#delete-quotas-project-id) | Reset all Quota of a project |
| GET | /v1/quotas | [get quotas](#get-quotas) | List Quotas |
| GET | /v1/quotas/defaults | [get quotas defaults](#get-quotas-defaults) | Show Quota Defaults |
| GET | /v1/quotas/{project_id} | [get quotas project ID](#get-quotas-project-id) | Show Quota detail |
| GET | /v1/services | [get services](#get-services) | List Services |
| POST | /v1/sync | [post sync](#post-sync) | Enqueue a full sync |
| PUT | /v1/quotas/{project_id} | [put quotas project ID](#put-quotas-project-id) | Update Quota |
  


###  datacenters

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/datacenters/{datacenter_id} | [delete datacenters datacenter ID](#delete-datacenters-datacenter-id) | Delete a datacenter |
| GET | /v1/datacenters | [get datacenters](#get-datacenters) | List datacenters |
| GET | /v1/datacenters/{datacenter_id} | [get datacenters datacenter ID](#get-datacenters-datacenter-id) | Show datacenter detail |
| POST | /v1/datacenters | [post datacenters](#post-datacenters) | Create new datacenter |
| PUT | /v1/datacenters/{datacenter_id} | [put datacenters datacenter ID](#put-datacenters-datacenter-id) | Update a datacenter |
  


###  domains

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/domains/{domain_id} | [delete domains domain ID](#delete-domains-domain-id) | Delete a domain |
| GET | /v1/domains | [get domains](#get-domains) | List domains |
| GET | /v1/domains/{domain_id} | [get domains domain ID](#get-domains-domain-id) | Show domain detail |
| POST | /v1/domains | [post domains](#post-domains) | Create new domain |
| PUT | /v1/domains/{domain_id} | [put domains domain ID](#put-domains-domain-id) | Update a domain |
  


###  geographic_maps

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/geomaps/{geomap_id} | [delete geomaps geomap ID](#delete-geomaps-geomap-id) | Delete a geographic map |
| GET | /v1/geomaps | [get geomaps](#get-geomaps) | List geographic maps |
| GET | /v1/geomaps/{geomap_id} | [get geomaps geomap ID](#get-geomaps-geomap-id) | Show geographic map detail |
| POST | /v1/geomaps | [post geomaps](#post-geomaps) | Create new geographic map |
| PUT | /v1/geomaps/{geomap_id} | [put geomaps geomap ID](#put-geomaps-geomap-id) | Update a geographic map |
  


###  members

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/members/{member_id} | [delete members member ID](#delete-members-member-id) | Delete a member |
| GET | /v1/members | [get members](#get-members) | List members |
| GET | /v1/members/{member_id} | [get members member ID](#get-members-member-id) | Show member detail |
| POST | /v1/members | [post members](#post-members) | Create new member |
| PUT | /v1/members/{member_id} | [put members member ID](#put-members-member-id) | Update a member |
  


###  monitors

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/monitors/{monitor_id} | [delete monitors monitor ID](#delete-monitors-monitor-id) | Delete a monitor |
| GET | /v1/monitors | [get monitors](#get-monitors) | List monitors |
| GET | /v1/monitors/{monitor_id} | [get monitors monitor ID](#get-monitors-monitor-id) | Show monitor detail |
| POST | /v1/monitors | [post monitors](#post-monitors) | Create new monitor |
| PUT | /v1/monitors/{monitor_id} | [put monitors monitor ID](#put-monitors-monitor-id) | Update a monitor |
  


###  pools

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /v1/pools/{pool_id} | [delete pools pool ID](#delete-pools-pool-id) | Delete a pool |
| GET | /v1/pools | [get pools](#get-pools) | List pools |
| GET | /v1/pools/{pool_id} | [get pools pool ID](#get-pools-pool-id) | Show pool detail |
| POST | /v1/pools | [post pools](#post-pools) | Create new pool |
| PUT | /v1/pools/{pool_id} | [put pools pool ID](#put-pools-pool-id) | Update a pool |
  


## Paths

### <span id="delete-datacenters-datacenter-id"></span> Delete a datacenter (*DeleteDatacentersDatacenterID*)

```
DELETE /v1/datacenters/{datacenter_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| datacenter_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the datacenter |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-datacenters-datacenter-id-204) | No Content | Resource successfully deleted. |  | [schema](#delete-datacenters-datacenter-id-204-schema) |
| [404](#delete-datacenters-datacenter-id-404) | Not Found | Not Found |  | [schema](#delete-datacenters-datacenter-id-404-schema) |
| [default](#delete-datacenters-datacenter-id-default) | | Unexpected Error |  | [schema](#delete-datacenters-datacenter-id-default-schema) |

#### Responses


##### <span id="delete-datacenters-datacenter-id-204"></span> 204 - Resource successfully deleted.
Status: No Content

###### <span id="delete-datacenters-datacenter-id-204-schema"></span> Schema

##### <span id="delete-datacenters-datacenter-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-datacenters-datacenter-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-datacenters-datacenter-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-datacenters-datacenter-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="delete-domains-domain-id"></span> Delete a domain (*DeleteDomainsDomainID*)

```
DELETE /v1/domains/{domain_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| domain_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the domain |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-domains-domain-id-204) | No Content | Resource successfully deleted. |  | [schema](#delete-domains-domain-id-204-schema) |
| [404](#delete-domains-domain-id-404) | Not Found | Not Found |  | [schema](#delete-domains-domain-id-404-schema) |
| [default](#delete-domains-domain-id-default) | | Unexpected Error |  | [schema](#delete-domains-domain-id-default-schema) |

#### Responses


##### <span id="delete-domains-domain-id-204"></span> 204 - Resource successfully deleted.
Status: No Content

###### <span id="delete-domains-domain-id-204-schema"></span> Schema

##### <span id="delete-domains-domain-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-domains-domain-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-domains-domain-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-domains-domain-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="delete-geomaps-geomap-id"></span> Delete a geographic map (*DeleteGeomapsGeomapID*)

```
DELETE /v1/geomaps/{geomap_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| geomap_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the geomap |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-geomaps-geomap-id-204) | No Content | Resource successfully deleted. |  | [schema](#delete-geomaps-geomap-id-204-schema) |
| [404](#delete-geomaps-geomap-id-404) | Not Found | Not Found |  | [schema](#delete-geomaps-geomap-id-404-schema) |
| [default](#delete-geomaps-geomap-id-default) | | Unexpected Error |  | [schema](#delete-geomaps-geomap-id-default-schema) |

#### Responses


##### <span id="delete-geomaps-geomap-id-204"></span> 204 - Resource successfully deleted.
Status: No Content

###### <span id="delete-geomaps-geomap-id-204-schema"></span> Schema

##### <span id="delete-geomaps-geomap-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-geomaps-geomap-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-geomaps-geomap-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-geomaps-geomap-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="delete-members-member-id"></span> Delete a member (*DeleteMembersMemberID*)

```
DELETE /v1/members/{member_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| member_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the member |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-members-member-id-204) | No Content | Resource successfully deleted. |  | [schema](#delete-members-member-id-204-schema) |
| [404](#delete-members-member-id-404) | Not Found | Not Found |  | [schema](#delete-members-member-id-404-schema) |
| [default](#delete-members-member-id-default) | | Unexpected Error |  | [schema](#delete-members-member-id-default-schema) |

#### Responses


##### <span id="delete-members-member-id-204"></span> 204 - Resource successfully deleted.
Status: No Content

###### <span id="delete-members-member-id-204-schema"></span> Schema

##### <span id="delete-members-member-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-members-member-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-members-member-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-members-member-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="delete-monitors-monitor-id"></span> Delete a monitor (*DeleteMonitorsMonitorID*)

```
DELETE /v1/monitors/{monitor_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| monitor_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the monitor |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-monitors-monitor-id-204) | No Content | Resource successfully deleted. |  | [schema](#delete-monitors-monitor-id-204-schema) |
| [404](#delete-monitors-monitor-id-404) | Not Found | Not Found |  | [schema](#delete-monitors-monitor-id-404-schema) |
| [default](#delete-monitors-monitor-id-default) | | Unexpected Error |  | [schema](#delete-monitors-monitor-id-default-schema) |

#### Responses


##### <span id="delete-monitors-monitor-id-204"></span> 204 - Resource successfully deleted.
Status: No Content

###### <span id="delete-monitors-monitor-id-204-schema"></span> Schema

##### <span id="delete-monitors-monitor-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-monitors-monitor-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-monitors-monitor-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-monitors-monitor-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="delete-pools-pool-id"></span> Delete a pool (*DeletePoolsPoolID*)

```
DELETE /v1/pools/{pool_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| pool_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the pool |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-pools-pool-id-204) | No Content | Resource successfully deleted. |  | [schema](#delete-pools-pool-id-204-schema) |
| [404](#delete-pools-pool-id-404) | Not Found | Not Found |  | [schema](#delete-pools-pool-id-404-schema) |
| [default](#delete-pools-pool-id-default) | | Unexpected Error |  | [schema](#delete-pools-pool-id-default-schema) |

#### Responses


##### <span id="delete-pools-pool-id-204"></span> 204 - Resource successfully deleted.
Status: No Content

###### <span id="delete-pools-pool-id-204-schema"></span> Schema

##### <span id="delete-pools-pool-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-pools-pool-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-pools-pool-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-pools-pool-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="delete-quotas-project-id"></span> Reset all Quota of a project (*DeleteQuotasProjectID*)

```
DELETE /v1/quotas/{project_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | string | `string` |  | ✓ |  | The ID of the project to query. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#delete-quotas-project-id-204) | No Content | Resource successfully reseted. |  | [schema](#delete-quotas-project-id-204-schema) |
| [404](#delete-quotas-project-id-404) | Not Found | Not Found |  | [schema](#delete-quotas-project-id-404-schema) |
| [default](#delete-quotas-project-id-default) | | Unexpected Error |  | [schema](#delete-quotas-project-id-default-schema) |

#### Responses


##### <span id="delete-quotas-project-id-204"></span> 204 - Resource successfully reseted.
Status: No Content

###### <span id="delete-quotas-project-id-204-schema"></span> Schema

##### <span id="delete-quotas-project-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-quotas-project-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="delete-quotas-project-id-default"></span> Default Response
Unexpected Error

###### <span id="delete-quotas-project-id-default-schema"></span> Schema

  

[Error](#error)

### <span id="get-datacenters"></span> List datacenters (*GetDatacenters*)

```
GET /v1/datacenters
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| limit | `query` | integer | `int64` |  |  |  | Sets the page size. |
| marker | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pagination ID of the last item in the previous list. |
| not-tags | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple not-tags are considered as logical AND.
Should be provided in a comma separated list. |
| not-tags-any | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |
| page_reverse | `query` | boolean | `bool` |  |  |  | Sets the page direction. |
| sort | `query` | string | `string` |  |  |  | Comma-separated list of sort keys, optinally prefix with - to reverse sort order. |
| tags | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical AND. 
Should be provided in a comma separated list. |
| tags-any | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-datacenters-200) | OK | A JSON array of datacenters |  | [schema](#get-datacenters-200-schema) |
| [400](#get-datacenters-400) | Bad Request | Bad request |  | [schema](#get-datacenters-400-schema) |
| [default](#get-datacenters-default) | | Unexpected Error |  | [schema](#get-datacenters-default-schema) |

#### Responses


##### <span id="get-datacenters-200"></span> 200 - A JSON array of datacenters
Status: OK

###### <span id="get-datacenters-200-schema"></span> Schema
   
  

[GetDatacentersOKBody](#get-datacenters-o-k-body)

##### <span id="get-datacenters-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="get-datacenters-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-datacenters-default"></span> Default Response
Unexpected Error

###### <span id="get-datacenters-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-datacenters-o-k-body"></span> GetDatacentersOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenters | [][Datacenter](#datacenter)| `[]*models.Datacenter` |  | |  |  |
| links | [][Link](#link)| `[]*models.Link` |  | |  |  |



### <span id="get-datacenters-datacenter-id"></span> Show datacenter detail (*GetDatacentersDatacenterID*)

```
GET /v1/datacenters/{datacenter_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| datacenter_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the datacenter |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-datacenters-datacenter-id-200) | OK | Shows the details of a specific datacenter. |  | [schema](#get-datacenters-datacenter-id-200-schema) |
| [404](#get-datacenters-datacenter-id-404) | Not Found | Not Found |  | [schema](#get-datacenters-datacenter-id-404-schema) |
| [default](#get-datacenters-datacenter-id-default) | | Unexpected Error |  | [schema](#get-datacenters-datacenter-id-default-schema) |

#### Responses


##### <span id="get-datacenters-datacenter-id-200"></span> 200 - Shows the details of a specific datacenter.
Status: OK

###### <span id="get-datacenters-datacenter-id-200-schema"></span> Schema
   
  

[GetDatacentersDatacenterIDOKBody](#get-datacenters-datacenter-id-o-k-body)

##### <span id="get-datacenters-datacenter-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-datacenters-datacenter-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-datacenters-datacenter-id-default"></span> Default Response
Unexpected Error

###### <span id="get-datacenters-datacenter-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-datacenters-datacenter-id-o-k-body"></span> GetDatacentersDatacenterIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenter | [Datacenter](#datacenter)| `models.Datacenter` |  | |  |  |



### <span id="get-domains"></span> List domains (*GetDomains*)

```
GET /v1/domains
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| limit | `query` | integer | `int64` |  |  |  | Sets the page size. |
| marker | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pagination ID of the last item in the previous list. |
| not-tags | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple not-tags are considered as logical AND.
Should be provided in a comma separated list. |
| not-tags-any | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |
| page_reverse | `query` | boolean | `bool` |  |  |  | Sets the page direction. |
| sort | `query` | string | `string` |  |  |  | Comma-separated list of sort keys, optinally prefix with - to reverse sort order. |
| tags | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical AND. 
Should be provided in a comma separated list. |
| tags-any | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-domains-200) | OK | A JSON array of domains |  | [schema](#get-domains-200-schema) |
| [400](#get-domains-400) | Bad Request | Bad request |  | [schema](#get-domains-400-schema) |
| [default](#get-domains-default) | | Unexpected Error |  | [schema](#get-domains-default-schema) |

#### Responses


##### <span id="get-domains-200"></span> 200 - A JSON array of domains
Status: OK

###### <span id="get-domains-200-schema"></span> Schema
   
  

[GetDomainsOKBody](#get-domains-o-k-body)

##### <span id="get-domains-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="get-domains-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-domains-default"></span> Default Response
Unexpected Error

###### <span id="get-domains-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-domains-o-k-body"></span> GetDomainsOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domains | [][Domain](#domain)| `[]*models.Domain` |  | |  |  |
| links | [][Link](#link)| `[]*models.Link` |  | |  |  |



### <span id="get-domains-domain-id"></span> Show domain detail (*GetDomainsDomainID*)

```
GET /v1/domains/{domain_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| domain_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the domain |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-domains-domain-id-200) | OK | Shows the details of a specific domain. |  | [schema](#get-domains-domain-id-200-schema) |
| [404](#get-domains-domain-id-404) | Not Found | Not Found |  | [schema](#get-domains-domain-id-404-schema) |
| [default](#get-domains-domain-id-default) | | Unexpected Error |  | [schema](#get-domains-domain-id-default-schema) |

#### Responses


##### <span id="get-domains-domain-id-200"></span> 200 - Shows the details of a specific domain.
Status: OK

###### <span id="get-domains-domain-id-200-schema"></span> Schema
   
  

[GetDomainsDomainIDOKBody](#get-domains-domain-id-o-k-body)

##### <span id="get-domains-domain-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-domains-domain-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-domains-domain-id-default"></span> Default Response
Unexpected Error

###### <span id="get-domains-domain-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-domains-domain-id-o-k-body"></span> GetDomainsDomainIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domain | [Domain](#domain)| `models.Domain` |  | |  |  |



### <span id="get-geomaps"></span> List geographic maps (*GetGeomaps*)

```
GET /v1/geomaps
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| limit | `query` | integer | `int64` |  |  |  | Sets the page size. |
| marker | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pagination ID of the last item in the previous list. |
| not-tags | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple not-tags are considered as logical AND.
Should be provided in a comma separated list. |
| not-tags-any | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |
| page_reverse | `query` | boolean | `bool` |  |  |  | Sets the page direction. |
| sort | `query` | string | `string` |  |  |  | Comma-separated list of sort keys, optinally prefix with - to reverse sort order. |
| tags | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical AND. 
Should be provided in a comma separated list. |
| tags-any | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-geomaps-200) | OK | A JSON array of geographic maps |  | [schema](#get-geomaps-200-schema) |
| [400](#get-geomaps-400) | Bad Request | Bad request |  | [schema](#get-geomaps-400-schema) |
| [default](#get-geomaps-default) | | Unexpected Error |  | [schema](#get-geomaps-default-schema) |

#### Responses


##### <span id="get-geomaps-200"></span> 200 - A JSON array of geographic maps
Status: OK

###### <span id="get-geomaps-200-schema"></span> Schema
   
  

[GetGeomapsOKBody](#get-geomaps-o-k-body)

##### <span id="get-geomaps-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="get-geomaps-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-geomaps-default"></span> Default Response
Unexpected Error

###### <span id="get-geomaps-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-geomaps-o-k-body"></span> GetGeomapsOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| geomaps | [][Geomap](#geomap)| `[]*models.Geomap` |  | |  |  |
| links | [][Link](#link)| `[]*models.Link` |  | |  |  |



### <span id="get-geomaps-geomap-id"></span> Show geographic map detail (*GetGeomapsGeomapID*)

```
GET /v1/geomaps/{geomap_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| geomap_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the geomap |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-geomaps-geomap-id-200) | OK | Shows the details of a specific geomap. |  | [schema](#get-geomaps-geomap-id-200-schema) |
| [404](#get-geomaps-geomap-id-404) | Not Found | Not Found |  | [schema](#get-geomaps-geomap-id-404-schema) |
| [default](#get-geomaps-geomap-id-default) | | Unexpected Error |  | [schema](#get-geomaps-geomap-id-default-schema) |

#### Responses


##### <span id="get-geomaps-geomap-id-200"></span> 200 - Shows the details of a specific geomap.
Status: OK

###### <span id="get-geomaps-geomap-id-200-schema"></span> Schema
   
  

[GetGeomapsGeomapIDOKBody](#get-geomaps-geomap-id-o-k-body)

##### <span id="get-geomaps-geomap-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-geomaps-geomap-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-geomaps-geomap-id-default"></span> Default Response
Unexpected Error

###### <span id="get-geomaps-geomap-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-geomaps-geomap-id-o-k-body"></span> GetGeomapsGeomapIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| geomap | [Geomap](#geomap)| `models.Geomap` |  | |  |  |



### <span id="get-members"></span> List members (*GetMembers*)

```
GET /v1/members
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| limit | `query` | integer | `int64` |  |  |  | Sets the page size. |
| marker | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pagination ID of the last item in the previous list. |
| not-tags | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple not-tags are considered as logical AND.
Should be provided in a comma separated list. |
| not-tags-any | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |
| page_reverse | `query` | boolean | `bool` |  |  |  | Sets the page direction. |
| pool_id | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pool ID of the members to fetch |
| sort | `query` | string | `string` |  |  |  | Comma-separated list of sort keys, optinally prefix with - to reverse sort order. |
| tags | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical AND. 
Should be provided in a comma separated list. |
| tags-any | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-members-200) | OK | A JSON array of members |  | [schema](#get-members-200-schema) |
| [400](#get-members-400) | Bad Request | Bad request |  | [schema](#get-members-400-schema) |
| [default](#get-members-default) | | Unexpected Error |  | [schema](#get-members-default-schema) |

#### Responses


##### <span id="get-members-200"></span> 200 - A JSON array of members
Status: OK

###### <span id="get-members-200-schema"></span> Schema
   
  

[GetMembersOKBody](#get-members-o-k-body)

##### <span id="get-members-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="get-members-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-members-default"></span> Default Response
Unexpected Error

###### <span id="get-members-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-members-o-k-body"></span> GetMembersOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| links | [][Link](#link)| `[]*models.Link` |  | |  |  |
| members | [][Member](#member)| `[]*models.Member` |  | |  |  |



### <span id="get-members-member-id"></span> Show member detail (*GetMembersMemberID*)

```
GET /v1/members/{member_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| member_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the member |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-members-member-id-200) | OK | Shows the details of a specific member. |  | [schema](#get-members-member-id-200-schema) |
| [404](#get-members-member-id-404) | Not Found | Not Found |  | [schema](#get-members-member-id-404-schema) |
| [default](#get-members-member-id-default) | | Unexpected Error |  | [schema](#get-members-member-id-default-schema) |

#### Responses


##### <span id="get-members-member-id-200"></span> 200 - Shows the details of a specific member.
Status: OK

###### <span id="get-members-member-id-200-schema"></span> Schema
   
  

[GetMembersMemberIDOKBody](#get-members-member-id-o-k-body)

##### <span id="get-members-member-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-members-member-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-members-member-id-default"></span> Default Response
Unexpected Error

###### <span id="get-members-member-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-members-member-id-o-k-body"></span> GetMembersMemberIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| member | [Member](#member)| `models.Member` |  | |  |  |



### <span id="get-monitors"></span> List monitors (*GetMonitors*)

```
GET /v1/monitors
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| limit | `query` | integer | `int64` |  |  |  | Sets the page size. |
| marker | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pagination ID of the last item in the previous list. |
| not-tags | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple not-tags are considered as logical AND.
Should be provided in a comma separated list. |
| not-tags-any | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |
| page_reverse | `query` | boolean | `bool` |  |  |  | Sets the page direction. |
| pool_id | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pool ID of the monitors to fetch |
| sort | `query` | string | `string` |  |  |  | Comma-separated list of sort keys, optinally prefix with - to reverse sort order. |
| tags | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical AND. 
Should be provided in a comma separated list. |
| tags-any | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-monitors-200) | OK | A JSON array of monitors |  | [schema](#get-monitors-200-schema) |
| [400](#get-monitors-400) | Bad Request | Bad request |  | [schema](#get-monitors-400-schema) |
| [default](#get-monitors-default) | | Unexpected Error |  | [schema](#get-monitors-default-schema) |

#### Responses


##### <span id="get-monitors-200"></span> 200 - A JSON array of monitors
Status: OK

###### <span id="get-monitors-200-schema"></span> Schema
   
  

[GetMonitorsOKBody](#get-monitors-o-k-body)

##### <span id="get-monitors-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="get-monitors-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-monitors-default"></span> Default Response
Unexpected Error

###### <span id="get-monitors-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-monitors-o-k-body"></span> GetMonitorsOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| links | [][Link](#link)| `[]*models.Link` |  | |  |  |
| monitors | [][Monitor](#monitor)| `[]*models.Monitor` |  | |  |  |



### <span id="get-monitors-monitor-id"></span> Show monitor detail (*GetMonitorsMonitorID*)

```
GET /v1/monitors/{monitor_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| monitor_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the monitor |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-monitors-monitor-id-200) | OK | Shows the details of a specific monitor. |  | [schema](#get-monitors-monitor-id-200-schema) |
| [404](#get-monitors-monitor-id-404) | Not Found | Not Found |  | [schema](#get-monitors-monitor-id-404-schema) |
| [default](#get-monitors-monitor-id-default) | | Unexpected Error |  | [schema](#get-monitors-monitor-id-default-schema) |

#### Responses


##### <span id="get-monitors-monitor-id-200"></span> 200 - Shows the details of a specific monitor.
Status: OK

###### <span id="get-monitors-monitor-id-200-schema"></span> Schema
   
  

[GetMonitorsMonitorIDOKBody](#get-monitors-monitor-id-o-k-body)

##### <span id="get-monitors-monitor-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-monitors-monitor-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-monitors-monitor-id-default"></span> Default Response
Unexpected Error

###### <span id="get-monitors-monitor-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-monitors-monitor-id-o-k-body"></span> GetMonitorsMonitorIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| monitor | [Monitor](#monitor)| `models.Monitor` |  | |  |  |



### <span id="get-pools"></span> List pools (*GetPools*)

```
GET /v1/pools
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| limit | `query` | integer | `int64` |  |  |  | Sets the page size. |
| marker | `query` | uuid (formatted string) | `strfmt.UUID` |  |  |  | Pagination ID of the last item in the previous list. |
| not-tags | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple not-tags are considered as logical AND.
Should be provided in a comma separated list. |
| not-tags-any | `query` | []string | `[]string` |  |  |  | Filter for resources not having tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |
| page_reverse | `query` | boolean | `bool` |  |  |  | Sets the page direction. |
| sort | `query` | string | `string` |  |  |  | Comma-separated list of sort keys, optinally prefix with - to reverse sort order. |
| tags | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical AND. 
Should be provided in a comma separated list. |
| tags-any | `query` | []string | `[]string` |  |  |  | Filter for tags, multiple tags are considered as logical OR.
Should be provided in a comma separated list. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-pools-200) | OK | A JSON array of pools |  | [schema](#get-pools-200-schema) |
| [400](#get-pools-400) | Bad Request | Bad request |  | [schema](#get-pools-400-schema) |
| [default](#get-pools-default) | | Unexpected Error |  | [schema](#get-pools-default-schema) |

#### Responses


##### <span id="get-pools-200"></span> 200 - A JSON array of pools
Status: OK

###### <span id="get-pools-200-schema"></span> Schema
   
  

[GetPoolsOKBody](#get-pools-o-k-body)

##### <span id="get-pools-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="get-pools-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-pools-default"></span> Default Response
Unexpected Error

###### <span id="get-pools-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-pools-o-k-body"></span> GetPoolsOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| links | [][Link](#link)| `[]*models.Link` |  | |  |  |
| pools | [][Pool](#pool)| `[]*models.Pool` |  | |  |  |



### <span id="get-pools-pool-id"></span> Show pool detail (*GetPoolsPoolID*)

```
GET /v1/pools/{pool_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| pool_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the pool |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-pools-pool-id-200) | OK | Shows the details of a specific pool. |  | [schema](#get-pools-pool-id-200-schema) |
| [404](#get-pools-pool-id-404) | Not Found | Not Found |  | [schema](#get-pools-pool-id-404-schema) |
| [default](#get-pools-pool-id-default) | | Unexpected Error |  | [schema](#get-pools-pool-id-default-schema) |

#### Responses


##### <span id="get-pools-pool-id-200"></span> 200 - Shows the details of a specific pool.
Status: OK

###### <span id="get-pools-pool-id-200-schema"></span> Schema
   
  

[GetPoolsPoolIDOKBody](#get-pools-pool-id-o-k-body)

##### <span id="get-pools-pool-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-pools-pool-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-pools-pool-id-default"></span> Default Response
Unexpected Error

###### <span id="get-pools-pool-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-pools-pool-id-o-k-body"></span> GetPoolsPoolIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| pool | [Pool](#pool)| `models.Pool` |  | |  |  |



### <span id="get-quotas"></span> List Quotas (*GetQuotas*)

```
GET /v1/quotas
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `query` | string | `string` |  |  |  | The ID of the project to query. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-quotas-200) | OK | A JSON array of quotas |  | [schema](#get-quotas-200-schema) |
| [default](#get-quotas-default) | | Unexpected Error |  | [schema](#get-quotas-default-schema) |

#### Responses


##### <span id="get-quotas-200"></span> 200 - A JSON array of quotas
Status: OK

###### <span id="get-quotas-200-schema"></span> Schema
   
  

[GetQuotasOKBody](#get-quotas-o-k-body)

##### <span id="get-quotas-default"></span> Default Response
Unexpected Error

###### <span id="get-quotas-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-quotas-o-k-body"></span> GetQuotasOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| quotas | [][GetQuotasOKBodyQuotasItems0](#get-quotas-o-k-body-quotas-items0)| `[]*GetQuotasOKBodyQuotasItems0` |  | |  |  |



**<span id="get-quotas-o-k-body-quotas-items0"></span> GetQuotasOKBodyQuotasItems0**


  


* composed type [Quota](#quota)
* inlined member (*AO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |



### <span id="get-quotas-defaults"></span> Show Quota Defaults (*GetQuotasDefaults*)

```
GET /v1/quotas/defaults
```

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-quotas-defaults-200) | OK | Show the quota defaults configured for new projects. |  | [schema](#get-quotas-defaults-200-schema) |
| [default](#get-quotas-defaults-default) | | Unexpected Error |  | [schema](#get-quotas-defaults-default-schema) |

#### Responses


##### <span id="get-quotas-defaults-200"></span> 200 - Show the quota defaults configured for new projects.
Status: OK

###### <span id="get-quotas-defaults-200-schema"></span> Schema
   
  

[GetQuotasDefaultsOKBody](#get-quotas-defaults-o-k-body)

##### <span id="get-quotas-defaults-default"></span> Default Response
Unexpected Error

###### <span id="get-quotas-defaults-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-quotas-defaults-o-k-body"></span> GetQuotasDefaultsOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| quota | [Quota](#quota)| `models.Quota` |  | |  |  |



### <span id="get-quotas-project-id"></span> Show Quota detail (*GetQuotasProjectID*)

```
GET /v1/quotas/{project_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | string | `string` |  | ✓ |  | The ID of the project to query. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-quotas-project-id-200) | OK | Shows the details of a specific monitor. |  | [schema](#get-quotas-project-id-200-schema) |
| [404](#get-quotas-project-id-404) | Not Found | Not Found |  | [schema](#get-quotas-project-id-404-schema) |
| [default](#get-quotas-project-id-default) | | Unexpected Error |  | [schema](#get-quotas-project-id-default-schema) |

#### Responses


##### <span id="get-quotas-project-id-200"></span> 200 - Shows the details of a specific monitor.
Status: OK

###### <span id="get-quotas-project-id-200-schema"></span> Schema
   
  

[GetQuotasProjectIDOKBody](#get-quotas-project-id-o-k-body)

##### <span id="get-quotas-project-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-quotas-project-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="get-quotas-project-id-default"></span> Default Response
Unexpected Error

###### <span id="get-quotas-project-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-quotas-project-id-o-k-body"></span> GetQuotasProjectIDOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| quota | [](#)| `` |  | |  |  |



### <span id="get-services"></span> List Services (*GetServices*)

```
GET /v1/services
```

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-services-200) | OK | A JSON array of services |  | [schema](#get-services-200-schema) |
| [default](#get-services-default) | | Unexpected Error |  | [schema](#get-services-default-schema) |

#### Responses


##### <span id="get-services-200"></span> 200 - A JSON array of services
Status: OK

###### <span id="get-services-200-schema"></span> Schema
   
  

[GetServicesOKBody](#get-services-o-k-body)

##### <span id="get-services-default"></span> Default Response
Unexpected Error

###### <span id="get-services-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="get-services-o-k-body"></span> GetServicesOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| services | [][Service](#service)| `[]*models.Service` |  | |  |  |



### <span id="post-datacenters"></span> Create new datacenter (*PostDatacenters*)

```
POST /v1/datacenters
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| datacenter | `body` | [PostDatacentersBody](#post-datacenters-body) | `PostDatacentersBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-datacenters-201) | Created | Created datacenter. |  | [schema](#post-datacenters-201-schema) |
| [404](#post-datacenters-404) | Not Found | Not Found |  | [schema](#post-datacenters-404-schema) |
| [default](#post-datacenters-default) | | Unexpected Error |  | [schema](#post-datacenters-default-schema) |

#### Responses


##### <span id="post-datacenters-201"></span> 201 - Created datacenter.
Status: Created

###### <span id="post-datacenters-201-schema"></span> Schema
   
  

[PostDatacentersCreatedBody](#post-datacenters-created-body)

##### <span id="post-datacenters-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-datacenters-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-datacenters-default"></span> Default Response
Unexpected Error

###### <span id="post-datacenters-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-datacenters-body"></span> PostDatacentersBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenter | [Datacenter](#datacenter)| `models.Datacenter` | ✓ | |  |  |



**<span id="post-datacenters-created-body"></span> PostDatacentersCreatedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenter | [Datacenter](#datacenter)| `models.Datacenter` |  | |  |  |



### <span id="post-domains"></span> Create new domain (*PostDomains*)

```
POST /v1/domains
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| domain | `body` | [PostDomainsBody](#post-domains-body) | `PostDomainsBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-domains-201) | Created | Created domain. |  | [schema](#post-domains-201-schema) |
| [400](#post-domains-400) | Bad Request | Bad request |  | [schema](#post-domains-400-schema) |
| [default](#post-domains-default) | | Unexpected Error |  | [schema](#post-domains-default-schema) |

#### Responses


##### <span id="post-domains-201"></span> 201 - Created domain.
Status: Created

###### <span id="post-domains-201-schema"></span> Schema
   
  

[PostDomainsCreatedBody](#post-domains-created-body)

##### <span id="post-domains-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="post-domains-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-domains-default"></span> Default Response
Unexpected Error

###### <span id="post-domains-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-domains-body"></span> PostDomainsBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domain | [Domain](#domain)| `models.Domain` | ✓ | |  |  |



**<span id="post-domains-created-body"></span> PostDomainsCreatedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domain | [Domain](#domain)| `models.Domain` |  | |  |  |



### <span id="post-geomaps"></span> Create new geographic map (*PostGeomaps*)

```
POST /v1/geomaps
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| geomap | `body` | [PostGeomapsBody](#post-geomaps-body) | `PostGeomapsBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-geomaps-201) | Created | Created geomap. |  | [schema](#post-geomaps-201-schema) |
| [400](#post-geomaps-400) | Bad Request | Bad request |  | [schema](#post-geomaps-400-schema) |
| [404](#post-geomaps-404) | Not Found | Not Found |  | [schema](#post-geomaps-404-schema) |
| [default](#post-geomaps-default) | | Unexpected Error |  | [schema](#post-geomaps-default-schema) |

#### Responses


##### <span id="post-geomaps-201"></span> 201 - Created geomap.
Status: Created

###### <span id="post-geomaps-201-schema"></span> Schema
   
  

[PostGeomapsCreatedBody](#post-geomaps-created-body)

##### <span id="post-geomaps-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="post-geomaps-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-geomaps-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-geomaps-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-geomaps-default"></span> Default Response
Unexpected Error

###### <span id="post-geomaps-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-geomaps-body"></span> PostGeomapsBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| geomap | [Geomap](#geomap)| `models.Geomap` | ✓ | |  |  |



**<span id="post-geomaps-created-body"></span> PostGeomapsCreatedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| geomap | [Geomap](#geomap)| `models.Geomap` |  | |  |  |



### <span id="post-members"></span> Create new member (*PostMembers*)

```
POST /v1/members
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| member | `body` | [PostMembersBody](#post-members-body) | `PostMembersBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-members-201) | Created | Created member. |  | [schema](#post-members-201-schema) |
| [400](#post-members-400) | Bad Request | Bad request |  | [schema](#post-members-400-schema) |
| [404](#post-members-404) | Not Found | Not Found |  | [schema](#post-members-404-schema) |
| [default](#post-members-default) | | Unexpected Error |  | [schema](#post-members-default-schema) |

#### Responses


##### <span id="post-members-201"></span> 201 - Created member.
Status: Created

###### <span id="post-members-201-schema"></span> Schema
   
  

[PostMembersCreatedBody](#post-members-created-body)

##### <span id="post-members-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="post-members-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-members-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-members-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-members-default"></span> Default Response
Unexpected Error

###### <span id="post-members-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-members-body"></span> PostMembersBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| member | [Member](#member)| `models.Member` | ✓ | |  |  |



**<span id="post-members-created-body"></span> PostMembersCreatedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| member | [Member](#member)| `models.Member` |  | |  |  |



### <span id="post-monitors"></span> Create new monitor (*PostMonitors*)

```
POST /v1/monitors
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| monitor | `body` | [PostMonitorsBody](#post-monitors-body) | `PostMonitorsBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-monitors-201) | Created | Created monitor. |  | [schema](#post-monitors-201-schema) |
| [400](#post-monitors-400) | Bad Request | Bad request |  | [schema](#post-monitors-400-schema) |
| [404](#post-monitors-404) | Not Found | Not Found |  | [schema](#post-monitors-404-schema) |
| [default](#post-monitors-default) | | Unexpected Error |  | [schema](#post-monitors-default-schema) |

#### Responses


##### <span id="post-monitors-201"></span> 201 - Created monitor.
Status: Created

###### <span id="post-monitors-201-schema"></span> Schema
   
  

[PostMonitorsCreatedBody](#post-monitors-created-body)

##### <span id="post-monitors-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="post-monitors-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-monitors-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-monitors-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-monitors-default"></span> Default Response
Unexpected Error

###### <span id="post-monitors-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-monitors-body"></span> PostMonitorsBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| monitor | [Monitor](#monitor)| `models.Monitor` | ✓ | |  |  |



**<span id="post-monitors-created-body"></span> PostMonitorsCreatedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| monitor | [Monitor](#monitor)| `models.Monitor` |  | |  |  |



### <span id="post-pools"></span> Create new pool (*PostPools*)

```
POST /v1/pools
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| pool | `body` | [PostPoolsBody](#post-pools-body) | `PostPoolsBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-pools-201) | Created | Created pool. |  | [schema](#post-pools-201-schema) |
| [400](#post-pools-400) | Bad Request | Bad request |  | [schema](#post-pools-400-schema) |
| [default](#post-pools-default) | | Unexpected Error |  | [schema](#post-pools-default-schema) |

#### Responses


##### <span id="post-pools-201"></span> 201 - Created pool.
Status: Created

###### <span id="post-pools-201-schema"></span> Schema
   
  

[PostPoolsCreatedBody](#post-pools-created-body)

##### <span id="post-pools-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="post-pools-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="post-pools-default"></span> Default Response
Unexpected Error

###### <span id="post-pools-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-pools-body"></span> PostPoolsBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| pool | [Pool](#pool)| `models.Pool` | ✓ | |  |  |



**<span id="post-pools-created-body"></span> PostPoolsCreatedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| pool | [Pool](#pool)| `models.Pool` |  | |  |  |



### <span id="post-sync"></span> Enqueue a full sync (*PostSync*)

```
POST /v1/sync
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| domains | `body` | [PostSyncBody](#post-sync-body) | `PostSyncBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#post-sync-202) | Accepted | Full sync has been enqueued. |  | [schema](#post-sync-202-schema) |
| [default](#post-sync-default) | | Unexpected Error |  | [schema](#post-sync-default-schema) |

#### Responses


##### <span id="post-sync-202"></span> 202 - Full sync has been enqueued.
Status: Accepted

###### <span id="post-sync-202-schema"></span> Schema

##### <span id="post-sync-default"></span> Default Response
Unexpected Error

###### <span id="post-sync-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="post-sync-body"></span> PostSyncBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domains | []uuid (formatted string)| `[]strfmt.UUID` | ✓ | |  |  |



### <span id="put-datacenters-datacenter-id"></span> Update a datacenter (*PutDatacentersDatacenterID*)

```
PUT /v1/datacenters/{datacenter_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| datacenter_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the datacenter |
| datacenter | `body` | [PutDatacentersDatacenterIDBody](#put-datacenters-datacenter-id-body) | `PutDatacentersDatacenterIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-datacenters-datacenter-id-202) | Accepted | Updated datacenter. |  | [schema](#put-datacenters-datacenter-id-202-schema) |
| [404](#put-datacenters-datacenter-id-404) | Not Found | Not Found |  | [schema](#put-datacenters-datacenter-id-404-schema) |
| [default](#put-datacenters-datacenter-id-default) | | Unexpected Error |  | [schema](#put-datacenters-datacenter-id-default-schema) |

#### Responses


##### <span id="put-datacenters-datacenter-id-202"></span> 202 - Updated datacenter.
Status: Accepted

###### <span id="put-datacenters-datacenter-id-202-schema"></span> Schema
   
  

[PutDatacentersDatacenterIDAcceptedBody](#put-datacenters-datacenter-id-accepted-body)

##### <span id="put-datacenters-datacenter-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-datacenters-datacenter-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-datacenters-datacenter-id-default"></span> Default Response
Unexpected Error

###### <span id="put-datacenters-datacenter-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-datacenters-datacenter-id-accepted-body"></span> PutDatacentersDatacenterIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenter | [Datacenter](#datacenter)| `models.Datacenter` |  | |  |  |



**<span id="put-datacenters-datacenter-id-body"></span> PutDatacentersDatacenterIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenter | [Datacenter](#datacenter)| `models.Datacenter` | ✓ | |  |  |



### <span id="put-domains-domain-id"></span> Update a domain (*PutDomainsDomainID*)

```
PUT /v1/domains/{domain_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| domain_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the domain |
| domain | `body` | [PutDomainsDomainIDBody](#put-domains-domain-id-body) | `PutDomainsDomainIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-domains-domain-id-202) | Accepted | Updated domain. |  | [schema](#put-domains-domain-id-202-schema) |
| [400](#put-domains-domain-id-400) | Bad Request | Bad request |  | [schema](#put-domains-domain-id-400-schema) |
| [404](#put-domains-domain-id-404) | Not Found | Not Found |  | [schema](#put-domains-domain-id-404-schema) |
| [default](#put-domains-domain-id-default) | | Unexpected Error |  | [schema](#put-domains-domain-id-default-schema) |

#### Responses


##### <span id="put-domains-domain-id-202"></span> 202 - Updated domain.
Status: Accepted

###### <span id="put-domains-domain-id-202-schema"></span> Schema
   
  

[PutDomainsDomainIDAcceptedBody](#put-domains-domain-id-accepted-body)

##### <span id="put-domains-domain-id-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="put-domains-domain-id-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-domains-domain-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-domains-domain-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-domains-domain-id-default"></span> Default Response
Unexpected Error

###### <span id="put-domains-domain-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-domains-domain-id-accepted-body"></span> PutDomainsDomainIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domain | [Domain](#domain)| `models.Domain` |  | |  |  |



**<span id="put-domains-domain-id-body"></span> PutDomainsDomainIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| domain | [Domain](#domain)| `models.Domain` | ✓ | |  |  |



### <span id="put-geomaps-geomap-id"></span> Update a geographic map (*PutGeomapsGeomapID*)

```
PUT /v1/geomaps/{geomap_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| geomap_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the geomap |
| geomap | `body` | [PutGeomapsGeomapIDBody](#put-geomaps-geomap-id-body) | `PutGeomapsGeomapIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-geomaps-geomap-id-202) | Accepted | Updated geomap. |  | [schema](#put-geomaps-geomap-id-202-schema) |
| [400](#put-geomaps-geomap-id-400) | Bad Request | Bad request |  | [schema](#put-geomaps-geomap-id-400-schema) |
| [404](#put-geomaps-geomap-id-404) | Not Found | Not Found |  | [schema](#put-geomaps-geomap-id-404-schema) |
| [default](#put-geomaps-geomap-id-default) | | Unexpected Error |  | [schema](#put-geomaps-geomap-id-default-schema) |

#### Responses


##### <span id="put-geomaps-geomap-id-202"></span> 202 - Updated geomap.
Status: Accepted

###### <span id="put-geomaps-geomap-id-202-schema"></span> Schema
   
  

[PutGeomapsGeomapIDAcceptedBody](#put-geomaps-geomap-id-accepted-body)

##### <span id="put-geomaps-geomap-id-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="put-geomaps-geomap-id-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-geomaps-geomap-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-geomaps-geomap-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-geomaps-geomap-id-default"></span> Default Response
Unexpected Error

###### <span id="put-geomaps-geomap-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-geomaps-geomap-id-accepted-body"></span> PutGeomapsGeomapIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| geomap | [Geomap](#geomap)| `models.Geomap` |  | |  |  |



**<span id="put-geomaps-geomap-id-body"></span> PutGeomapsGeomapIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| geomap | [Geomap](#geomap)| `models.Geomap` | ✓ | |  |  |



### <span id="put-members-member-id"></span> Update a member (*PutMembersMemberID*)

```
PUT /v1/members/{member_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| member_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the member |
| member | `body` | [PutMembersMemberIDBody](#put-members-member-id-body) | `PutMembersMemberIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-members-member-id-202) | Accepted | Updated member. |  | [schema](#put-members-member-id-202-schema) |
| [400](#put-members-member-id-400) | Bad Request | Bad request |  | [schema](#put-members-member-id-400-schema) |
| [404](#put-members-member-id-404) | Not Found | Not Found |  | [schema](#put-members-member-id-404-schema) |
| [default](#put-members-member-id-default) | | Unexpected Error |  | [schema](#put-members-member-id-default-schema) |

#### Responses


##### <span id="put-members-member-id-202"></span> 202 - Updated member.
Status: Accepted

###### <span id="put-members-member-id-202-schema"></span> Schema
   
  

[PutMembersMemberIDAcceptedBody](#put-members-member-id-accepted-body)

##### <span id="put-members-member-id-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="put-members-member-id-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-members-member-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-members-member-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-members-member-id-default"></span> Default Response
Unexpected Error

###### <span id="put-members-member-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-members-member-id-accepted-body"></span> PutMembersMemberIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| member | [Member](#member)| `models.Member` |  | |  |  |



**<span id="put-members-member-id-body"></span> PutMembersMemberIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| member | [Member](#member)| `models.Member` | ✓ | |  |  |



### <span id="put-monitors-monitor-id"></span> Update a monitor (*PutMonitorsMonitorID*)

```
PUT /v1/monitors/{monitor_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| monitor_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the monitor |
| monitor | `body` | [PutMonitorsMonitorIDBody](#put-monitors-monitor-id-body) | `PutMonitorsMonitorIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-monitors-monitor-id-202) | Accepted | Updated monitor. |  | [schema](#put-monitors-monitor-id-202-schema) |
| [400](#put-monitors-monitor-id-400) | Bad Request | Bad request |  | [schema](#put-monitors-monitor-id-400-schema) |
| [404](#put-monitors-monitor-id-404) | Not Found | Not Found |  | [schema](#put-monitors-monitor-id-404-schema) |
| [default](#put-monitors-monitor-id-default) | | Unexpected Error |  | [schema](#put-monitors-monitor-id-default-schema) |

#### Responses


##### <span id="put-monitors-monitor-id-202"></span> 202 - Updated monitor.
Status: Accepted

###### <span id="put-monitors-monitor-id-202-schema"></span> Schema
   
  

[PutMonitorsMonitorIDAcceptedBody](#put-monitors-monitor-id-accepted-body)

##### <span id="put-monitors-monitor-id-400"></span> 400 - Bad request
Status: Bad Request

###### <span id="put-monitors-monitor-id-400-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-monitors-monitor-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-monitors-monitor-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-monitors-monitor-id-default"></span> Default Response
Unexpected Error

###### <span id="put-monitors-monitor-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-monitors-monitor-id-accepted-body"></span> PutMonitorsMonitorIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| monitor | [Monitor](#monitor)| `models.Monitor` |  | |  |  |



**<span id="put-monitors-monitor-id-body"></span> PutMonitorsMonitorIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| monitor | [Monitor](#monitor)| `models.Monitor` | ✓ | |  |  |



### <span id="put-pools-pool-id"></span> Update a pool (*PutPoolsPoolID*)

```
PUT /v1/pools/{pool_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| pool_id | `path` | uuid (formatted string) | `strfmt.UUID` |  | ✓ |  | The UUID of the pool |
| pool | `body` | [PutPoolsPoolIDBody](#put-pools-pool-id-body) | `PutPoolsPoolIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-pools-pool-id-202) | Accepted | Updated pool. |  | [schema](#put-pools-pool-id-202-schema) |
| [404](#put-pools-pool-id-404) | Not Found | Not Found |  | [schema](#put-pools-pool-id-404-schema) |
| [default](#put-pools-pool-id-default) | | Unexpected Error |  | [schema](#put-pools-pool-id-default-schema) |

#### Responses


##### <span id="put-pools-pool-id-202"></span> 202 - Updated pool.
Status: Accepted

###### <span id="put-pools-pool-id-202-schema"></span> Schema
   
  

[PutPoolsPoolIDAcceptedBody](#put-pools-pool-id-accepted-body)

##### <span id="put-pools-pool-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-pools-pool-id-404-schema"></span> Schema
   
  

[Error](#error)

##### <span id="put-pools-pool-id-default"></span> Default Response
Unexpected Error

###### <span id="put-pools-pool-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-pools-pool-id-accepted-body"></span> PutPoolsPoolIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| pool | [Pool](#pool)| `models.Pool` |  | |  |  |



**<span id="put-pools-pool-id-body"></span> PutPoolsPoolIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| pool | [Pool](#pool)| `models.Pool` | ✓ | |  |  |



### <span id="put-quotas-project-id"></span> Update Quota (*PutQuotasProjectID*)

```
PUT /v1/quotas/{project_id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | string | `string` |  | ✓ |  | The ID of the project to query. |
| quota | `body` | [PutQuotasProjectIDBody](#put-quotas-project-id-body) | `PutQuotasProjectIDBody` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [202](#put-quotas-project-id-202) | Accepted | Updated quota for a project. |  | [schema](#put-quotas-project-id-202-schema) |
| [default](#put-quotas-project-id-default) | | Unexpected Error |  | [schema](#put-quotas-project-id-default-schema) |

#### Responses


##### <span id="put-quotas-project-id-202"></span> 202 - Updated quota for a project.
Status: Accepted

###### <span id="put-quotas-project-id-202-schema"></span> Schema
   
  

[PutQuotasProjectIDAcceptedBody](#put-quotas-project-id-accepted-body)

##### <span id="put-quotas-project-id-default"></span> Default Response
Unexpected Error

###### <span id="put-quotas-project-id-default-schema"></span> Schema

  

[Error](#error)

###### Inlined models

**<span id="put-quotas-project-id-accepted-body"></span> PutQuotasProjectIDAcceptedBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| quota | [Quota](#quota)| `models.Quota` |  | |  |  |



**<span id="put-quotas-project-id-body"></span> PutQuotasProjectIDBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| quota | [Quota](#quota)| `models.Quota` | ✓ | |  |  |



## Models

### <span id="datacenter"></span> datacenter


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| admin_state_up | boolean| `bool` |  | `true`| The administrative state of the resource, which is up (true) or down (false). Default is true. |  |
| city | string| `string` |  | |  | `Berlin` |
| continent | string| `string` |  | |  | `EU` |
| country | string| `string` |  | |  | `DE` |
| created_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-05-11T17:21:34` |
| id | uuid (formatted string)| `strfmt.UUID` |  | | The id of the resource. |  |
| latitude | float (formatted number)| `float32` |  | `52.52`|  | `52.526055` |
| longitude | float (formatted number)| `float32` |  | `13.4`|  | `13.403454` |
| meta | integer| `int64` |  | |  |  |
| name | string| `string` |  | | Human-readable name of the resource. |  |
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |
| provider | string| `string` |  | | Provider driver for the backend solution | `akamai` |
| provisioning_status | string| `string` |  | |  |  |
| scope | string| `string` |  | `"private"`| Visibility of datacenter between different projects |  |
| state_or_province | string| `string` |  | |  | `Berlin` |
| updated_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-09-09T14:52:15` |



### <span id="domain"></span> domain


> A representation of a domain
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| admin_state_up | boolean| `bool` |  | `true`| The administrative state of the resource, which is up (true) or down (false). Default is true. |  |
| aliases | []string (formatted string)| `[]string` |  | |  |  |
| cname_target | hostname (formatted string)| `strfmt.Hostname` |  | | If not empty, the backend created a CNAME target to be used for the FQDN. | `example.org.production.gtm.com` |
| created_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-05-11T17:21:34` |
| fqdn | hostname (formatted string)| `strfmt.Hostname` | ✓ | | Desired Fully-Qualified Host Name. | `example.org` |
| id | uuid (formatted string)| `strfmt.UUID` |  | | The id of the resource. |  |
| mode | string| `string` |  | `"ROUND_ROBIN"`| Load balancing method to use for the references pools. |  |
| name | string| `string` |  | | Human-readable name of the resource. |  |
| pools | []uuid (formatted string)| `[]strfmt.UUID` |  | |  |  |
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |
| provider | string| `string` | ✓ | | Supported provider drivers | `akamai` |
| provisioning_status | string| `string` |  | |  |  |
| record_type | string| `string` |  | `"A"`| DNS Record type to use. |  |
| status | string| `string` |  | |  |  |
| updated_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-09-09T14:52:15` |



### <span id="error"></span> error


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | | HTTP Error code | `404` |
| message | string| `string` |  | |  | `An example error message` |



### <span id="geomap"></span> geomap


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| assignments | [][GeomapAssignmentsItems0](#geomap-assignments-items0)| `[]*GeomapAssignmentsItems0` |  | | Country to datacenter assignments. |  |
| created_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-05-11T17:21:34` |
| default_datacenter | uuid (formatted string)| `strfmt.UUID` | ✓ | | Datacenter ID |  |
| id | uuid (formatted string)| `strfmt.UUID` |  | | The id of the resource. |  |
| name | string| `string` |  | | Human-readable name of the resource. |  |
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |
| provider | string| `string` |  | | Provider driver for the backend solution | `akamai` |
| provisioning_status | string| `string` |  | |  |  |
| scope | string| `string` |  | `"private"`| Visibility of datacenter between different projects |  |
| updated_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-09-09T14:52:15` |



#### Inlined models

**<span id="geomap-assignments-items0"></span> GeomapAssignmentsItems0**


> Assignment.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| country | string| `string` |  | | ISO 3166 2-Letter Country code. |  |
| datacenter | uuid (formatted string)| `strfmt.UUID` |  | | Datacenter ID |  |



### <span id="link"></span> link


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| href | uri (formatted string)| `strfmt.URI` |  | |  |  |
| rel | string| `string` |  | |  |  |



### <span id="member"></span> member


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| address | ipv4 (formatted string)| `strfmt.IPv4` | ✓ | | Address to use. | `1.2.3.4` |
| admin_state_up | boolean| `bool` |  | `true`| The administrative state of the resource, which is up (true) or down (false). Default is true. |  |
| created_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-05-11T17:21:34` |
| datacenter_id | uuid (formatted string)| `strfmt.UUID` |  | | Datacenter assigned for this member. |  |
| id | uuid (formatted string)| `strfmt.UUID` |  | | The id of the resource. |  |
| name | string| `string` |  | | Human-readable name of the resource. |  |
| pool_id | uuid (formatted string)| `strfmt.UUID` |  | | pool id. |  |
| port | integer| `int64` | ✓ | | Port to use for monitor checks. | `80` |
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |
| provisioning_status | string| `string` |  | |  |  |
| status | string| `string` |  | |  |  |
| updated_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-09-09T14:52:15` |



### <span id="monitor"></span> monitor


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| admin_state_up | boolean| `bool` |  | `true`| The administrative state of the resource, which is up (true) or down (false). Default is true. |  |
| created_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-05-11T17:21:34` |
| id | uuid (formatted string)| `strfmt.UUID` |  | | The id of the resource. |  |
| interval | integer| `int64` |  | `60`| The interval, in seconds, between health checks. | `10` |
| name | string| `string` |  | | Human-readable name of the resource. |  |
| pool_id | uuid (formatted string)| `strfmt.UUID` |  | | ID of the pool to check members |  |
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |
| provisioning_status | string| `string` |  | |  |  |
| receive | string| `string` |  | | Specifies the text string that the monitor expects to receive from the target member. | `HTTP/1.` |
| send | string| `string` |  | | Specifies the text string that the monitor sends to the target member. | `HEAD / HTTP/1.0\\r\\n\\r\\n` |
| timeout | integer| `int64` |  | `10`| The time in total, in seconds, after which a health check times out. | `30` |
| type | string| `string` |  | `"ICMP"`| Type of the health check monitor. | `HTTP` |
| updated_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-09-09T14:52:15` |



### <span id="pool"></span> pool


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| admin_state_up | boolean| `bool` |  | `true`| The administrative state of the resource, which is up (true) or down (false). Default is true. |  |
| created_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-05-11 17:21:34` |
| domains | []uuid (formatted string)| `[]strfmt.UUID` |  | | Array of domains assigned to this pool |  |
| id | uuid (formatted string)| `strfmt.UUID` |  | | The id of the resource. |  |
| members | []uuid (formatted string)| `[]strfmt.UUID` |  | | Array of member ids that this pool uses for load balancing. |  |
| monitors | []uuid (formatted string)| `[]strfmt.UUID` |  | | Array of monitor ids that this pool uses health checks. |  |
| name | string| `string` |  | | Human-readable name of the resource. |  |
| project_id | string| `string` |  | | The ID of the project owning this resource. | `fa84c217f361441986a220edf9b1e337` |
| provisioning_status | string| `string` |  | |  |  |
| status | string| `string` |  | |  |  |
| updated_at | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when the resource was created. | `2020-09-09 14:52:15` |



### <span id="quota"></span> quota


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| datacenter | integer| `int64` |  | | The configured datacenter quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited. | `5` |
| domain | integer| `int64` |  | | The configured domain quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited. | `5` |
| member | integer| `int64` |  | | The configured member quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited. | `5` |
| monitor | integer| `int64` |  | | The configured monitor quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited. | `5` |
| pool | integer| `int64` |  | | The configured pool quota limit. A setting of null means it is using the deployment default quota. A setting of -1 means unlimited. | `5` |



### <span id="quota-usage"></span> quota_usage


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| in_use_datacenter | integer| `int64` |  | | The current quota usage of datacenter. | `5` |
| in_use_domain | integer| `int64` |  | | The current quota usage of domain. | `5` |
| in_use_member | integer| `int64` |  | | The current quota usage of member. | `5` |
| in_use_monitor | integer| `int64` |  | | The current quota usage of monitor. | `5` |
| in_use_pool | integer| `int64` |  | | The current quota usage of pool. | `5` |



### <span id="service"></span> service


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| heartbeat | dateTime (formatted string)| `string` |  | | The UTC date and timestamp when had the last heartbeat. | `2020-05-11 17:21:34` |
| host | hostname (formatted string)| `strfmt.Hostname` |  | | Hostname of the computer the service is running. | `example.host` |
| id | string| `string` |  | | ID of the RPC service. | `andromeda-agent-fbb49979-03f5-4a97-a334-1fd2c9f61e7e` |
| metadata | [interface{}](#interface)| `interface{}` |  | |  |  |
| provider | string| `string` |  | | Provider this service supports. | `akamai` |
| rpc_address | string| `string` |  | | RPC Endpoint Address. | `_INBOX.VEfFxcAzZQ9iM9vwGH49It` |
| type | string| `string` |  | | Type of service. | `healthcheck` |
| version | string| `string` |  | | Version of the service. | `1.2.3` |


