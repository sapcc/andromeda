/*
 *   Copyright 2024 SAP SE
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

# trick to convert float -> double without conversion corruption
# https://stackoverflow.com/questions/30024990/mysql-convert-float-to-double

ALTER TABLE `datacenter` CHANGE COLUMN `latitude` `latitude` VARCHAR(255) NOT NULL DEFAULT '52.52';
ALTER TABLE `datacenter` CHANGE COLUMN `longitude` `longitude` VARCHAR(255) NOT NULL DEFAULT '13.40';

ALTER TABLE `datacenter` CHANGE COLUMN `latitude` `latitude` DOUBLE NOT NULL DEFAULT 52.52;
ALTER TABLE `datacenter` CHANGE COLUMN `longitude` `longitude` DOUBLE NOT NULL DEFAULT 13.40;
