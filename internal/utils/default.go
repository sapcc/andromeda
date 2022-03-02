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

package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/iancoleman/strcase"
)

var (
	SwaggerSpec *loads.Document
)

func SetModelDefaults(s interface{}) error {
	/*
		This is a workaround for default swagger values, since go-swagger currently doesn't populate default variables
		for nested definitions:
		https://github.com/go-swagger/go-swagger/issues/1393
	*/
	var instanceType string
	if _, err := fmt.Sscanf(fmt.Sprintf("%T", s), "*models.%s", &instanceType); err != nil {
		return err
	}
	instanceType = strings.ToLower(instanceType)
	for specDefinitionName, specDefinitionModel := range SwaggerSpec.Spec().Definitions {
		if specDefinitionName == instanceType {

			// Found the swagger model
			for propName, property := range specDefinitionModel.SchemaProps.Properties {

				// Check if model has default set
				if property.Default != nil {
					propertyField := reflect.ValueOf(s).Elem().FieldByName(strcase.ToCamel(propName))
					if propertyField.Kind() != reflect.Ptr && propertyField.Kind() != reflect.Uintptr {
						return fmt.Errorf("unexpected field %s for specDefinitionModel %s", propName, specDefinitionName)
					}

					if !propertyField.IsNil() {
						continue
					}

					// Generate correct Value
					vp := reflect.New(propertyField.Type())
					switch property.Default.(type) {
					case bool:
						val := property.Default.(bool)
						vp.Elem().Set(reflect.ValueOf(&val))
					case string:
						val := property.Default.(string)
						vp.Elem().Set(reflect.ValueOf(&val))
					case int64:
						val := property.Default.(int64)
						vp.Elem().Set(reflect.ValueOf(&val))
					case float64:
						val := property.Default.(float64)
						_tmp := int64(val)
						vp.Elem().Set(reflect.ValueOf(&_tmp))
					default:
						return fmt.Errorf("unexpected type %T for property %s", property.Default, propName)
					}
					propertyField.Set(vp.Elem())
				}
			}
			break
		}
	}
	return nil
}
