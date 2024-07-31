/*
 *   Copyright 2020 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License);
 *   you may NOT use this file except in compliance with the License.
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

CREATE TABLE domain
(
    id                  UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL DEFAULT 'PENDING_CREATE',
    status              VARCHAR(16)  NOT NULL DEFAULT 'OFFLINE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL DEFAULT now(),
    provider            VARCHAR(64)  NOT NULL,
    fqdn                VARCHAR(512) NOT NULL,
    record_type         VARCHAR(4)   NOT NULL,
    mode                VARCHAR(16)  NOT NULL,
    project_id          VARCHAR(36)  NOT NULL,
    UNIQUE (fqdn, provider)
);

CREATE TABLE pool
(
    id                  UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL DEFAULT 'PENDING_CREATE',
    status              VARCHAR(16)  NOT NULL DEFAULT 'OFFLINE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL DEFAULT now(),
    project_id          VARCHAR(36)  NOT NULL
);

CREATE TABLE domain_pool_relation
(
    domain_id UUID NOT NULL REFERENCES domain ON DELETE CASCADE,
    pool_id   UUID NOT NULL REFERENCES pool ON DELETE CASCADE,
    PRIMARY KEY (domain_id, pool_id)
);

CREATE TABLE datacenter
(
    id                  UUID PRIMARY KEY          DEFAULT gen_random_uuid(),
    name                VARCHAR(255)     NULL,
    provisioning_status VARCHAR(16)      NOT NULL DEFAULT 'PENDING_CREATE',
    admin_state_up      BOOLEAN          NOT NULL,
    created_at          TIMESTAMP        NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP        NOT NULL DEFAULT now(),
    state_or_province   VARCHAR(255)     NOT NULL DEFAULT '',
    city                VARCHAR(255)     NOT NULL DEFAULT '',
    continent           VARCHAR(2)       NOT NULL DEFAULT '',
    country             VARCHAR(2)       NOT NULL DEFAULT '',
    latitude            DOUBLE PRECISION NOT NULL DEFAULT 52.52,
    longitude           DOUBLE PRECISION NOT NULL DEFAULT 13.40,
    scope               VARCHAR(8)       NOT NULL DEFAULT 'private' CHECK ( scope IN ('private', 'public')),
    project_id          VARCHAR(36)      NOT NULL,
    provider            VARCHAR(64)      NOT NULL,
    meta                INTEGER          NOT NULL DEFAULT 0
);

CREATE TABLE member
(
    id                  UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL DEFAULT 'PENDING_CREATE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL DEFAULT now(),
    port                BIGINT       NOT NULL,
    status              VARCHAR(16)  NOT NULL DEFAULT 'UNKNOWN',
    address             VARCHAR(255) NOT NULL,
    pool_id             UUID         NOT NULL REFERENCES pool ON DELETE CASCADE,
    project_id          VARCHAR(36)  NOT NULL,
    datacenter_id       UUID         NULL REFERENCES datacenter,
    UNIQUE (pool_id, address, port)
);

CREATE TABLE monitor
(
    id                  UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL DEFAULT 'PENDING_CREATE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL DEFAULT now(),
    interval            BIGINT       NOT NULL,
    pool_id             UUID         NOT NULL REFERENCES pool ON DELETE CASCADE,
    receive             VARCHAR(255) NULL,
    send                VARCHAR(255) NULL,
    http_method         VARCHAR(16)  NOT NULL DEFAULT 'GET',
    timeout             BIGINT       NULL,
    type                VARCHAR(16)  NULL,
    domain_name         VARCHAR(255) NULL,
    project_id          VARCHAR(36)  NOT NULL
);

CREATE TABLE quota
(
    project_id VARCHAR(36) PRIMARY KEY,
    domain     BIGINT NOT NULL,
    pool       BIGINT NOT NULL,
    member     BIGINT NOT NULL,
    monitor    BIGINT NOT NULL,
    datacenter BIGINT NOT NULL
);

CREATE TABLE agent
(
    host           VARCHAR(36) PRIMARY KEY,
    admin_state_up BOOLEAN,
    heartbeat      TIMESTAMP NOT NULL,
    providers      JSON
);