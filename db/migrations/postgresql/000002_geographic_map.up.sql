/*
 *   Copyright 2020 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

CREATE TABLE geographic_map
(
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL DEFAULT 'PENDING_CREATE',
    created_at          TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL DEFAULT now(),
    default_datacenter  UUID         NOT NULL REFERENCES datacenter(id),
    scope               VARCHAR(8)   NOT NULL DEFAULT 'private' CHECK ( scope IN ('private', 'public')),
    project_id          VARCHAR(36)  NOT NULL,
    provider            VARCHAR(64)  NOT NULL
);

CREATE TABLE geographic_map_assignment
(
    geographic_map_id   UUID         NOT NULL REFERENCES geographic_map(id) ON DELETE CASCADE,
    datacenter          UUID         NOT NULL REFERENCES datacenter(id),
    country             VARCHAR(2)   NOT NULL
);
