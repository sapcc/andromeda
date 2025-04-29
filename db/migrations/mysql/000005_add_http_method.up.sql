-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE monitor ADD COLUMN http_method VARCHAR(16) NOT NULL DEFAULT 'GET';
