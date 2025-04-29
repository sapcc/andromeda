-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

# trick to convert float -> double without conversion corruption
# https://stackoverflow.com/questions/30024990/mysql-convert-float-to-double

ALTER TABLE `datacenter` CHANGE COLUMN `latitude` `latitude` VARCHAR(255) NOT NULL DEFAULT '52.52';
ALTER TABLE `datacenter` CHANGE COLUMN `longitude` `longitude` VARCHAR(255) NOT NULL DEFAULT '13.40';

ALTER TABLE `datacenter` CHANGE COLUMN `latitude` `latitude` DOUBLE NOT NULL DEFAULT 52.52;
ALTER TABLE `datacenter` CHANGE COLUMN `longitude` `longitude` DOUBLE NOT NULL DEFAULT 13.40;
