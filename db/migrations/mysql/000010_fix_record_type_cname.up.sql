-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE domain MODIFY record_type VARCHAR(5) NOT NULL;
UPDATE domain SET record_type = 'CNAME' WHERE record_type = 'CNAM';
