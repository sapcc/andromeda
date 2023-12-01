/*
 *   Copyright 2023 SAP SE
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

# domain
ALTER TABLE domain ADD COLUMN created_at_datetime DATETIME;
ALTER TABLE domain ADD COLUMN updated_at_datetime DATETIME;

UPDATE domain SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE domain SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE domain DROP COLUMN created_at;
ALTER TABLE domain DROP COLUMN updated_at;

ALTER TABLE domain RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE domain RENAME COLUMN updated_at_datetime TO updated_at_datetime;

# pool
ALTER TABLE pool ADD COLUMN created_at_datetime DATETIME;
ALTER TABLE pool ADD COLUMN updated_at_datetime DATETIME;

UPDATE pool SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE pool SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE pool DROP COLUMN created_at;
ALTER TABLE pool DROP COLUMN updated_at;

ALTER TABLE pool RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE pool RENAME COLUMN updated_at_datetime TO updated_at_datetime;

# datacenter
ALTER TABLE datacenter ADD COLUMN created_at_datetime DATETIME;
ALTER TABLE datacenter ADD COLUMN updated_at_datetime DATETIME;

UPDATE datacenter SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE datacenter SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE datacenter DROP COLUMN created_at;
ALTER TABLE datacenter DROP COLUMN updated_at;

ALTER TABLE datacenter RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE datacenter RENAME COLUMN updated_at_datetime TO updated_at_datetime;

# member
ALTER TABLE member ADD COLUMN created_at_datetime DATETIME;
ALTER TABLE member ADD COLUMN updated_at_datetime DATETIME;

UPDATE member SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE member SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE member DROP COLUMN created_at;
ALTER TABLE member DROP COLUMN updated_at;

ALTER TABLE member RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE member RENAME COLUMN updated_at_datetime TO updated_at_datetime;

# monitor
ALTER TABLE monitor ADD COLUMN created_at_datetime DATETIME;
ALTER TABLE monitor ADD COLUMN updated_at_datetime DATETIME;

UPDATE monitor SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE monitor SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE monitor DROP COLUMN created_at;
ALTER TABLE monitor DROP COLUMN updated_at;

ALTER TABLE monitor RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE monitor RENAME COLUMN updated_at_datetime TO updated_at_datetime;

# geographic_map
ALTER TABLE geographic_map ADD COLUMN created_at_datetime DATETIME;
ALTER TABLE geographic_map ADD COLUMN updated_at_datetime DATETIME;

UPDATE geographic_map SET created_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(created_at));
UPDATE geographic_map SET updated_at_datetime=FROM_UNIXTIME(UNIX_TIMESTAMP(updated_at));

ALTER TABLE geographic_map DROP COLUMN created_at;
ALTER TABLE geographic_map DROP COLUMN updated_at;

ALTER TABLE geographic_map RENAME COLUMN created_at_datetime TO created_at;
ALTER TABLE geographic_map RENAME COLUMN updated_at_datetime TO updated_at_datetime;
