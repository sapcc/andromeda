-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

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
