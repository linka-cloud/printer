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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

type Format int

const (
	Table Format = iota
	JSON
	YAML
)

func (f Format) String() string {
	switch f {
	case Table:
		return "table"
	case JSON:
		return "json"
	case YAML:
		return "yaml"
	default:
		return "unknown"
	}
}

func Print(v any, opts ...Option) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("panic recovered: %v: \n%v", r, string(debug.Stack())))
		}
	}()
	p := printer{
		writer:     tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0),
		format:     Table,
		max:        math.MaxInt,
		json:       json.Marshal,
		yaml:       yaml.Marshal,
		formatters: make(map[string]func(v any) string),
	}
	for _, v := range opts {
		v(&p)
	}
	switch p.format {
	case JSON:
		var (
			b   []byte
			err error
		)
		if p.json != nil {
			b, err = p.json(v)
		} else {
			b, err = json.Marshal(v)
		}
		if err != nil {
			return err
		}
		if _, err = fmt.Fprintln(p.writer, string(b)); err != nil {
			return err
		}
		return p.writer.Flush()
	case YAML:
		var (
			b   []byte
			err error
		)
		if p.yaml != nil {
			b, err = p.yaml(v)
		} else {
			b, err = yaml.Marshal(v)
		}
		if err != nil {
			return err
		}
		if _, err = fmt.Fprintln(p.writer, string(b)); err != nil {
			return err
		}
		return p.writer.Flush()
	case Table:
		cols, err := makeColumns(v)
		if err != nil {
			return err
		}
		cols = cols.Filter(p.fields...)
		if p.max > len(cols.Headers()) {
			p.max = len(cols.Headers())
		}
		if !p.noHeaders {
			head := strings.Join(cols.Headers()[:p.max], "\t")
			if p.LowerHeaders {
				head = strings.ToLower(head)
			}
			if p.UpperHeaders {
				head = strings.ToUpper(head)
			}
			if _, err = fmt.Fprintln(p.writer, head); err != nil {
				return err
			}
		}
		for _, v := range arr(v) {
			row := strings.Join(cols.Values(v, p.formatters)[:p.max], "\t")
			if p.LowerValues {
				row = strings.ToLower(row)
			}
			if p.UpperValues {
				row = strings.ToUpper(row)
			}
			if _, err := fmt.Fprintln(p.writer, row); err != nil {
				return err
			}
		}
		return p.writer.Flush()
	default:
		return fmt.Errorf("unknown format %q", p.format)
	}

}
