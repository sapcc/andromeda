/*
 *   Copyright 2020 SAP SE
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

package db

import (
	"fmt"
	"strings"
)

func GetNamedUpdateSets(values map[string]interface{}) string {
	var sql strings.Builder
	for key := range values {
		sql.WriteString(fmt.Sprintf("\"%s\"=:%s, ", key, key))
	}
	return sql.String()
}

func GetNamedWhere(values map[string]interface{}) string {
	var sql strings.Builder
	var size = len(values)
	var i = 0

	for key := range values {
		sql.WriteString(fmt.Sprintf("\"%s\" = :%s", key, key))

		i += 1
		if i < size {
			sql.WriteString(" AND ")
		}
	}
	return sql.String()
}
