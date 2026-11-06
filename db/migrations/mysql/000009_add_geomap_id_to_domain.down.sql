-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE domain DROP FOREIGN KEY fk_domain_geomap_id, DROP COLUMN geomap_id;
