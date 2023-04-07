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
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
	"github.com/sapcc/andromeda/models"
)

var DefaultSortKeys = []string{"id", "created_at"}

var (
	ErrInvalidMarker   = errors.New("invalid marker")
	ErrPolicyForbidden = errors.New("forbidden by policy")
)

type Pagination struct {
	/*Sets the page size.
	  In: query
	*/
	limit *int64
	/*Pagination ID of the last item in the previous list.
	  In: query
	*/
	marker *strfmt.UUID
	/*Sets the page direction.
	  In: query
	*/
	pageReverse *bool
	/*Comma-separated list of sort keys, optinally prefix with - to reverse sort order.
	  In: query
	*/
	sort *string

	table string
	r     *regexp.Regexp
}

func NewPagination(Table string, Limit *int64, Marker *strfmt.UUID, Sort *string, pageReverse *bool) *Pagination {
	return &Pagination{
		limit:       Limit,
		marker:      Marker,
		sort:        Sort,
		pageReverse: pageReverse,
		table:       Table,
		r:           regexp.MustCompile("^[a-z0-9_]+$"),
	}
}

func stripDesc(sortDirKey string) (string, bool) {
	sortKey := strings.TrimPrefix(sortDirKey, "-")
	return sortKey, sortKey != sortDirKey
}

// Query pagination helper that also includes policy query filter
func (p *Pagination) Query(db *sqlx.DB, r *http.Request, filter []string) (*sqlx.Rows, error) {
	var sortDirKeys []string
	var whereClauses []string
	var orderBy string
	markerObj := make(map[string]interface{})

	projectID, err := auth.ProjectScopeForRequest(r)
	if err != nil {
		return nil, err
	}
	if policy.Engine.AuthorizeGetAllRequest(r, projectID) {
		// Allow fetch of all resources
		projectID = ""
	} else if !policy.Engine.AuthorizeRequest(r, projectID) {
		return nil, ErrPolicyForbidden
	}

	query := fmt.Sprintf(`SELECT * FROM %s`, p.table)

	//add filter
	whereClauses = append(whereClauses, filter...)

	//add sorting
	if !config.Global.ApiSettings.DisableSorting && p.sort != nil {
		sortDirKeys = strings.Split(*p.sort, ",")

		// Add default sort keys (if not existing)
		for _, defaultSortKey := range DefaultSortKeys {
			found := false
			for _, paramSortKey := range sortDirKeys {
				sortKey, _ := stripDesc(paramSortKey)
				if sortKey == defaultSortKey {
					found = true
					break
				}
			}

			if !found {
				sortDirKeys = append(sortDirKeys, defaultSortKey)
			}
		}
	} else {
		// Creates a copy
		sortDirKeys = append(sortDirKeys, DefaultSortKeys...)
	}

	//always order to ensure stable result
	orderBy += " ORDER BY "
	for i, sortDirKey := range sortDirKeys {
		// Input sanitation
		if !p.r.MatchString(sortDirKey) {
			continue
		}

		if sortKey, ok := stripDesc(sortDirKey); ok {
			orderBy += fmt.Sprintf("%s DESC", sortKey)
		} else {
			orderBy += sortDirKey
		}

		if i < len(sortDirKeys)-1 {
			orderBy += ", "
		}
	}

	if !config.Global.ApiSettings.DisablePagination && p.marker != nil {
		sql := db.Rebind(fmt.Sprintf(`SELECT * FROM %s WHERE id=?`, p.table))
		rows, err := db.Queryx(sql, p.marker)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			if err := rows.MapScan(markerObj); err != nil {
				return nil, err
			}
		}

		if len(markerObj) == 0 {
			return nil, ErrInvalidMarker
		}

		// Craft WHERE ... conditions
		var sortWhereClauses strings.Builder
		for i, sortDirKey := range sortDirKeys {
			var critAttrs []string = nil
			for j := range sortDirKeys[:i] {
				sortKey := strings.TrimPrefix(sortDirKeys[j], "-")
				critAttrs = append(critAttrs, fmt.Sprintf("%s = :%s", sortKey, sortKey))
			}

			if sortKey := strings.TrimPrefix(sortDirKey, "-"); sortKey != sortDirKey {
				critAttrs = append(critAttrs, fmt.Sprintf("%s < :%s", sortKey, sortKey))
			} else {
				critAttrs = append(critAttrs, fmt.Sprintf("%s > :%s", sortKey, sortKey))
			}

			sortWhereClauses.WriteString("( " + strings.Join(critAttrs, " AND ") + " )")

			if i < len(sortDirKeys)-1 {
				sortWhereClauses.WriteString(" OR ")
			}
		}
		whereClauses = append(whereClauses, sortWhereClauses.String())
	}

	//override project scope, ensures marker is not used for fetching others owner resources
	if projectID != "" {
		// hardcoded to datacenter, which allows a public scope for sharing datacenters
		if p.table == "datacenter" {
			whereClauses = append(whereClauses, "scope = 'public' OR project_id = :project_id")
		} else {
			whereClauses = append(whereClauses, "project_id = :project_id")
		}
		markerObj["project_id"] = projectID
	}

	//add WHERE
	if len(whereClauses) > 0 {
		query += " WHERE ( " + strings.Join(whereClauses, " ) AND ( ") + " )"
	}

	//add ORDER BY
	query += orderBy

	var limit = config.Global.ApiSettings.PaginationMaxLimit
	if p.limit != nil && *p.limit < config.Global.ApiSettings.PaginationMaxLimit {
		limit = *p.limit
	}
	query += fmt.Sprintf(" LIMIT %d", limit)

	return db.NamedQuery(query, markerObj)
}

func (p *Pagination) GetLinks(modelList interface{}, r *http.Request) []*models.Link {
	var links []*models.Link
	if reflect.TypeOf(modelList).Kind() != reflect.Slice {
		return nil
	}

	s := reflect.ValueOf(modelList)
	if s.Len() > 0 {
		var prevAttr, nextAttr []string
		first := s.Index(0).Elem().FieldByName("ID").String()
		last := s.Index(s.Len() - 1).Elem().FieldByName("ID").String()

		if p.sort != nil {
			prevAttr = append(prevAttr, fmt.Sprintf("sort=%s", *p.sort))
		}
		if p.limit != nil {
			prevAttr = append(prevAttr, fmt.Sprintf("limit=%d", *p.limit))
		}

		// Make a copy
		nextAttr = append(prevAttr[:0:0], prevAttr...)

		// Previous link of marker supplied
		if p.marker != nil {
			prevAttr = append(prevAttr, fmt.Sprintf("marker=%s", first), "page_reverse=True")
			prevUrl := fmt.Sprintf("%s%s?%s", config.Global.Default.ApiBaseURL, r.URL.Path, strings.Join(prevAttr, "&"))

			links = append(links, &models.Link{
				Href: strfmt.URI(prevUrl),
				Rel:  "previous",
			})
		}

		// Next link of limit < size(fetched items)
		if p.limit != nil && int64(s.Len()) >= *p.limit {
			nextAttr = append(nextAttr, fmt.Sprintf("marker=%s", last))
			nextUrl := fmt.Sprintf("%s%s?%s", config.Global.Default.ApiBaseURL, r.URL.Path, strings.Join(nextAttr, "&"))
			links = append(links, &models.Link{
				Href: strfmt.URI(nextUrl),
				Rel:  "next",
			})
		}
	}
	return links
}
