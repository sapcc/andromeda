-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE domain ALTER COLUMN record_type TYPE VARCHAR(5);
UPDATE domain SET record_type = 'CNAME' WHERE record_type = 'CNAM';
