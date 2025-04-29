-- SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
--
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE `quota` ADD COLUMN `domain_f5` bigint(20) NOT NULL AFTER `domain_akamai`;
