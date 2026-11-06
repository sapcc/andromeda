-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE domain ADD COLUMN geomap_id VARCHAR(36) NULL DEFAULT NULL,
    ADD CONSTRAINT fk_domain_geomap_id FOREIGN KEY (geomap_id) REFERENCES geographic_map (id) ON DELETE SET NULL;
