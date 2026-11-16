-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

UPDATE domain SET record_type = 'CNAM' WHERE record_type = 'CNAME';
ALTER TABLE domain MODIFY record_type VARCHAR(4) NOT NULL;
