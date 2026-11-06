-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE domain ADD COLUMN geomap_id UUID NULL DEFAULT NULL REFERENCES geographic_map (id) ON DELETE SET NULL;
