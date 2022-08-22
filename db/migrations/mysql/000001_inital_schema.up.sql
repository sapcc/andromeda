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
    id                  VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL   DEFAULT 'PENDING_CREATE',
    status              VARCHAR(16)  NOT NULL   DEFAULT 'OFFLINE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    provider            VARCHAR(64)  NOT NULL,
    fqdn                VARCHAR(512) NOT NULL,
    record_type         VARCHAR(4)   NOT NULL,
    mode                VARCHAR(16)  NOT NULL,
    project_id          VARCHAR(36)  NOT NULL,
    UNIQUE KEY (fqdn, provider)
) ENGINE = InnoDB;

CREATE TABLE pool
(
    id                  VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL   DEFAULT 'PENDING_CREATE',
    status              VARCHAR(16)  NOT NULL   DEFAULT 'OFFLINE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    project_id          VARCHAR(36)  NOT NULL
) ENGINE = InnoDB;

CREATE TABLE domain_pool_relation
(
    domain_id VARCHAR(36) NOT NULL,
    pool_id   VARCHAR(36) NOT NULL,
    CONSTRAINT FOREIGN KEY (domain_id) REFERENCES domain (id) ON DELETE CASCADE,
    CONSTRAINT FOREIGN KEY (pool_id) REFERENCES pool (id) ON DELETE CASCADE,
    CONSTRAINT PRIMARY KEY (domain_id, pool_id)
) ENGINE = InnoDB;

CREATE TABLE datacenter
(
    id                  VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL   DEFAULT 'PENDING_CREATE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    state_or_province   VARCHAR(255) NOT NULL   DEFAULT '',
    city                VARCHAR(255) NOT NULL   DEFAULT '',
    continent           VARCHAR(2)   NOT NULL   DEFAULT '',
    country             VARCHAR(2)   NOT NULL   DEFAULT '',
    latitude            FLOAT        NOT NULL   DEFAULT 52.52,
    longitude           FLOAT        NOT NULL   DEFAULT 13.40,
    scope               VARCHAR(8)   NOT NULL   DEFAULT 'private' CHECK ( scope IN ('private', 'public')),
    project_id          VARCHAR(36)  NOT NULL,
    provider            VARCHAR(64)  NOT NULL,
    meta                INT          NOT NULL
) ENGINE = InnoDB;

CREATE TABLE member
(
    id                  VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL   DEFAULT 'PENDING_CREATE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    port                BIGINT       NOT NULL,
    status              VARCHAR(16)  NOT NULL   DEFAULT 'UNKNOWN',
    address             VARCHAR(255) NOT NULL,
    pool_id             VARCHAR(36)  NOT NULL,
    project_id          VARCHAR(36)  NOT NULL,
    datacenter_id       VARCHAR(36)  NULL,
    CONSTRAINT FOREIGN KEY (pool_id) REFERENCES pool (id) ON DELETE CASCADE,
    CONSTRAINT FOREIGN KEY (datacenter_id) REFERENCES datacenter (id),
    UNIQUE KEY (pool_id, address, port)
) ENGINE = InnoDB;

CREATE TABLE monitor
(
    id                  VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
    name                VARCHAR(255) NULL,
    provisioning_status VARCHAR(16)  NOT NULL   DEFAULT 'PENDING_CREATE',
    admin_state_up      BOOLEAN      NOT NULL,
    created_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    updated_at          TIMESTAMP    NOT NULL   DEFAULT now(),
    `interval`          BIGINT       NOT NULL,
    pool_id             VARCHAR(36)  NOT NULL,
    receive             VARCHAR(255) NULL,
    send                VARCHAR(255) NULL,
    timeout             BIGINT       NULL,
    type                VARCHAR(16)  NULL,
    project_id          VARCHAR(36)  NOT NULL,
    CONSTRAINT FOREIGN KEY (pool_id) REFERENCES pool (id) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE quota
(
    project_id        VARCHAR(36) PRIMARY KEY,
    domain            BIGINT NOT NULL,
    pool              BIGINT NOT NULL,
    member            BIGINT NOT NULL,
    monitor           BIGINT NOT NULL,
    datacenter        BIGINT NOT NULL,
    in_use_domain     BIGINT NOT NULL DEFAULT 0,
    in_use_pool       BIGINT NOT NULL DEFAULT 0,
    in_use_member     BIGINT NOT NULL DEFAULT 0,
    in_use_monitor    BIGINT NOT NULL DEFAULT 0,
    in_use_datacenter BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT in_use_domain CHECK (`in_use_domain` >= 0),
    CONSTRAINT in_use_pool CHECK (`in_use_pool` >= 0),
    CONSTRAINT in_use_member CHECK (`in_use_member` >= 0),
    CONSTRAINT in_use_monitor CHECK (`in_use_monitor` >= 0),
    CONSTRAINT in_use_datacenter CHECK (`in_use_datacenter` >= 0)
) ENGINE = InnoDB;

CREATE TABLE agent
(
    host           VARCHAR(36) PRIMARY KEY,
    admin_state_up BOOLEAN,
    heartbeat      TIMESTAMP NOT NULL,
    providers      JSON
) ENGINE = InnoDB;
