// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package view

import "fmt"

type JsonMap map[string]interface{}

func (j JsonMap) GetValueAsString(key string) string {
	if _, isObj := j[key].(map[string]interface{}); isObj {
		return ""
	}
	if _, isArr := j[key].([]interface{}); isArr {
		return ""
	}
	if val, ok := j[key]; ok {
		return fmt.Sprint(val)
	}

	return ""
}

func (j JsonMap) GetObject(key string) JsonMap {
	if obj, isObj := j[key].(map[string]interface{}); isObj {
		return obj
	}
	return JsonMap{}
}

func (j JsonMap) GetObjectsArray(key string) []JsonMap {
	if array, ok := j[key].([]interface{}); ok {
		objectsArray := make([]JsonMap, 0)
		for _, el := range array {
			if obj, ok := el.(map[string]interface{}); ok {
				objectsArray = append(objectsArray, obj)
			}
		}
		return objectsArray
	}
	return []JsonMap{}
}

func ConvertYamlToJsonMap(yaml map[interface{}]interface{}) JsonMap {
	mapStringInterface := convertMapI2MapS(yaml)
	if result, ok := mapStringInterface.(map[string]interface{}); ok {
		return result
	}
	return nil
}

func convertMapI2MapS(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string:
				m[k2] = convertMapI2MapS(v2)
			default:
				m[fmt.Sprint(k)] = convertMapI2MapS(v2)
			}
		}
		v = m

	case []interface{}:
		for i, v2 := range x {
			x[i] = convertMapI2MapS(v2)
		}

	case map[string]interface{}:
		for k, v2 := range x {
			x[k] = convertMapI2MapS(v2)
		}
	}
	return v
}
