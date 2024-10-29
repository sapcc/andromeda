/*
 *   Copyright 2022 SAP SE
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

package utils

import (
	"math"
	"reflect"
	"strings"

	"github.com/apex/log"
)

func getSubFields(subFieldName string, fieldsToCompare []string) []string {
	var subFields []string
	for _, fieldToCompare := range fieldsToCompare {
		if strings.HasPrefix(fieldToCompare, subFieldName+".") {
			subFields = append(subFields, fieldToCompare[len(subFieldName)+1:])
		}
	}
	return subFields
}

func in(field string, fieldsToCompare []string) bool {
	for _, fieldToCompare := range fieldsToCompare {
		if field == fieldToCompare {
			return true
		}
	}
	return false
}

func deepValueEqualField(v1, v2 reflect.Value, fieldsToCompare []string) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqualField(v1.Index(i), v2.Index(i), fieldsToCompare) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			found := false
			for j := 0; j < v2.Len(); j++ {
				if deepValueEqualField(v1.Index(i), v2.Index(j), fieldsToCompare) {
					found = true
				}
			}
			if !found {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqualField(v1.Elem(), v2.Elem(), fieldsToCompare)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			fieldName := reflect.Indirect(v1).Type().Field(i).Name
			if in(fieldName, fieldsToCompare) {
				subFields := getSubFields(fieldName, fieldsToCompare)
				if !deepValueEqualField(v1.Field(i), v2.Field(i), subFields) {
					log.Debugf("Field '%s': '%+v' != '%+v'", fieldName, v1.Field(i).Interface(), v2.Field(i).Interface())
					return false
				}
			}
		}
		return true
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return deepValueEqualField(v1.Elem(), v2.Elem(), fieldsToCompare)
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !deepValueEqualField(val1, val2, fieldsToCompare) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	case reflect.Float64:
		// This is likely a longitude/latitude field, so we need to compare them with a delta
		return math.Abs(v1.Float()-v2.Float()) < 0.0001
	default:
		return v1.Interface() == v2.Interface()
	}
}

// DeepEqualFields like reflect.DeepEqual, but with field array to compare
func DeepEqualFields(x, y interface{}, fieldsToCompare []string) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		log.Errorf("DeepEqualFields: Type mismatch: %s != %s", v1.Type(), v2.Type())
		return false
	}
	ret := deepValueEqualField(v1, v2, fieldsToCompare)
	log.Debugf("DeepEqualFields(%s): equal=%t", v1.Type(), ret)
	return ret
}
