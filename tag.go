// Copyright 2023 Linka Cloud  All rights reserved.
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

package printer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parseTag(index int, fieldName, tag string) (column, error) {
	parts := strings.Split(tag, ",")
	if len(parts) > 2 {
		return column{}, fmt.Errorf("invalid tag for field %q: %q, tag should be `print:\"(name),(order)\"` or `print:\"-\"", fieldName, tag)
	}
	c := column{name: fieldName, header: fieldName, order: index}
	for _, v := range parts {
		v = strings.TrimSpace(v)
		if v == "-" {
			c.disabled = true
			return c, nil
		}
		if n, err := strconv.Atoi(v); err == nil {
			c.order = n
			continue
		}
		if v != "" {
			c.header = v
		}
	}
	return c, nil
}

func arr(v any) []any {
	r := reflect.ValueOf(v)
	if r.Kind() != reflect.Slice {
		return []any{v}
	}
	var out []any
	for i := 0; i < r.Len(); i++ {
		out = append(out, r.Index(i).Interface())
	}
	return out
}

func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func derefType(v reflect.Type) reflect.Type {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
