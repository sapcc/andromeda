-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

# domain
ALTER TABLE domain ADD COLUMN created_at_datetime DATETIME NOT NULL DEFAULT now();
ALTER TABLE domain ADD COLUMN updated_at_datetime DATETIME NOT NULL DEFAULT now();

UPDATE domain SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE domain SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE domain DROP COLUMN created_at;
ALTER TABLE domain DROP COLUMN updated_at;

ALTER TABLE domain RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE domain RENAME COLUMN updated_at_datetime TO updated_at;

# pool
ALTER TABLE pool ADD COLUMN created_at_datetime DATETIME NOT NULL DEFAULT now();
ALTER TABLE pool ADD COLUMN updated_at_datetime DATETIME NOT NULL DEFAULT now();

UPDATE pool SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE pool SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE pool DROP COLUMN created_at;
ALTER TABLE pool DROP COLUMN updated_at;

ALTER TABLE pool RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE pool RENAME COLUMN updated_at_datetime TO updated_at;

# datacenter
ALTER TABLE datacenter ADD COLUMN created_at_datetime DATETIME NOT NULL DEFAULT now();
ALTER TABLE datacenter ADD COLUMN updated_at_datetime DATETIME NOT NULL DEFAULT now();

UPDATE datacenter SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE datacenter SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE datacenter DROP COLUMN created_at;
ALTER TABLE datacenter DROP COLUMN updated_at;

ALTER TABLE datacenter RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE datacenter RENAME COLUMN updated_at_datetime TO updated_at;

# member
ALTER TABLE member ADD COLUMN created_at_datetime DATETIME NOT NULL DEFAULT now();
ALTER TABLE member ADD COLUMN updated_at_datetime DATETIME NOT NULL DEFAULT now();

UPDATE member SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE member SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE member DROP COLUMN created_at;
ALTER TABLE member DROP COLUMN updated_at;

ALTER TABLE member RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE member RENAME COLUMN updated_at_datetime TO updated_at;

# monitor
ALTER TABLE monitor ADD COLUMN created_at_datetime DATETIME NOT NULL DEFAULT now();
ALTER TABLE monitor ADD COLUMN updated_at_datetime DATETIME NOT NULL DEFAULT now();

UPDATE monitor SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE monitor SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE monitor DROP COLUMN created_at;
ALTER TABLE monitor DROP COLUMN updated_at;

ALTER TABLE monitor RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE monitor RENAME COLUMN updated_at_datetime TO updated_at;

# geographic_map
ALTER TABLE geographic_map ADD COLUMN created_at_datetime DATETIME NOT NULL DEFAULT now();
ALTER TABLE geographic_map ADD COLUMN updated_at_datetime DATETIME NOT NULL DEFAULT now();

UPDATE geographic_map SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE geographic_map SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE geographic_map DROP COLUMN created_at;
ALTER TABLE geographic_map DROP COLUMN updated_at;

ALTER TABLE geographic_map RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE geographic_map RENAME COLUMN updated_at_datetime TO updated_at;
