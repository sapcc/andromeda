// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

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

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/models"
)

var (
	sortDirKeyRegex  = regexp.MustCompile("^[a-z0-9_]+$")
	defaultSortKeys  = []string{"id", "created_at"}
	ErrInvalidMarker = errors.New("invalid marker")
)

type Pagination struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Sets the page size.
	  In: query
	*/
	Limit *int64
	/*Pagination ID of the last item in the previous list.
	  In: query
	*/
	Marker *strfmt.UUID
	/*Sets the page direction.
	  In: query
	*/
	PageReverse *bool
	/*Comma-separated list of sort keys, optionally prefix with - to reverse sort order.
	  In: query
	*/
	Sort *string
}

func stripDesc(sortDirKey string) (string, bool) {
	sortKey := strings.TrimPrefix(sortDirKey, "-")
	return sortKey, sortKey != sortDirKey
}

// Query pagination helper that also includes policy query filter
func (p *Pagination) Query(db *sqlx.DB, query string, filter map[string]any) (*sqlx.Rows, error) {
	var sortDirKeys []string
	var whereClauses []string
	var orderBy string
	var pageReverse bool

	// add filter
	for key := range filter {
		if strings.HasSuffix(query, "datacenter") {
			whereClauses = append(whereClauses, fmt.Sprintf("( %s = :%s OR scope = 'public')", key, key))
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = :%s", key, key))
		}
	}

	// page reverse
	if p.PageReverse != nil {
		pageReverse = *p.PageReverse
	}

	//add sorting
	if !config.Global.ApiSettings.DisableSorting && p.Sort != nil {
		sortDirKeys = strings.Split(*p.Sort, ",")

		// Add default sort keys (if not existing)
		for _, defaultSortKey := range defaultSortKeys {
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
		sortDirKeys = append(sortDirKeys, defaultSortKeys...)
	}

	// always order to ensure stable result
	orderBy += " ORDER BY "
	for i, sortDirKey := range sortDirKeys {
		// Input sanitation
		if !sortDirKeyRegex.MatchString(sortDirKey) {
			continue
		}

		sortKey, desc := stripDesc(sortDirKey)
		orderBy += sortKey
		if (desc && !pageReverse) || (!desc && pageReverse) {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}

		if i < len(sortDirKeys)-1 {
			orderBy += ", "
		}
	}

	if !config.Global.ApiSettings.DisablePagination && p.Marker != nil {
		sql := db.Rebind(fmt.Sprintf(`%s WHERE id = ?`, query))
		if err := db.Get(&filter, sql, p.Marker); err != nil {
			return nil, err
		}

		if len(filter) == 0 {
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

	//add WHERE
	if len(whereClauses) > 0 {
		query += " WHERE ( " + strings.Join(whereClauses, " ) AND ( ") + " )"
	}

	//add ORDER BY
	query += orderBy

	// maximum limit
	var maxLimit = config.Global.ApiSettings.PaginationMaxLimit
	if p.Limit == nil || (p.Limit != nil && *p.Limit > maxLimit) {
		p.Limit = &maxLimit
	}
	query += fmt.Sprint(" LIMIT ", *p.Limit)

	return db.NamedQuery(query, filter)
}

func (p *Pagination) GetLinks(modelList any) []*models.Link {
	var links []*models.Link
	if reflect.TypeOf(modelList).Kind() != reflect.Slice {
		return nil
	}

	s := reflect.ValueOf(modelList)
	if s.Len() > 0 {
		var prevAttr, nextAttr []string
		first := s.Index(0).Elem().FieldByName("ID").String()
		last := s.Index(s.Len() - 1).Elem().FieldByName("ID").String()

		if p.HTTPRequest != nil {
			for key, val := range p.HTTPRequest.URL.Query() {
				if key == "marker" || key == "page_reverse" {
					continue
				}
				prevAttr = append(prevAttr, fmt.Sprint(key, "=", val[0]))
			}
		}

		// Make a shallow copy
		nextAttr = append(prevAttr[:0:0], prevAttr...)

		// Previous link of marker supplied
		if p.Marker != nil {
			prevAttr = append(prevAttr, fmt.Sprintf("marker=%s", first), "page_reverse=True")
			prevUrl := fmt.Sprint(config.GetApiBaseUrl(p.HTTPRequest), p.HTTPRequest.URL.Path,
				"?", strings.Join(prevAttr, "&"))

			links = append(links, &models.Link{
				Href: strfmt.URI(prevUrl),
				Rel:  "previous",
			})
		}

		// Next link of limit < size(fetched items)
		if p.Limit != nil && int64(s.Len()) >= *p.Limit {
			nextAttr = append(nextAttr, fmt.Sprintf("marker=%s", last))
			nextUrl := fmt.Sprint(config.GetApiBaseUrl(p.HTTPRequest), p.HTTPRequest.URL.Path,
				"?", strings.Join(nextAttr, "&"))
			links = append(links, &models.Link{
				Href: strfmt.URI(nextUrl),
				Rel:  "next",
			})
		}
	}
	return links
}
