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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ghodss/yaml"

	"go.linka.cloud/printer/internal/slices"
)

type T1 struct {
	AString     string        `json:"fieldA" yaml:"field_a" print:"STRING,3"`
	BInt        int           `json:"fieldB" yaml:"field_b" print:"INT,4"`
	CBool       bool          `json:"fieldC" yaml:"field_c" print:"BOOL,2"`
	Dur         time.Duration `json:"fieldD" yaml:"field_d" print:"DURATION,1"`
	Time        time.Time     `json:"fieldT" yaml:"field_t" print:"TIME,4"`
	Object      *T1           `json:"fieldO" yaml:"field_o" print:"OBJECT,5"`
	Noop        any           `json:"-" yaml:"-" print:"-"`
	noop        any
	OtherObject *T2
	Slice       []int
	Zero        int `json:"fieldZ" yaml:"field_z" print:"This one is always zero,42"`
}

func (t T1) String() string {
	return "struct"
}

type T2 struct {
	A string
}

func (t T2) String() string {
	return "struct 2"
}

func If[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

func P[T any](v T) *T {
	return &v
}

func RandomT1(depth int) T1 {
	return T1{
		AString: string(slices.Map(make([]int32, rand.Intn(8)+4), func(_ int32) rune {
			return []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")[rand.Intn(32)]
		})),
		BInt:        rand.Intn(42),
		CBool:       If(rand.Intn(2) == 0, true, false),
		Dur:         time.Duration(rand.Intn(int(time.Hour.Milliseconds()))),
		Time:        time.UnixMicro(rand.Int63n(time.Now().UnixMicro())),
		Object:      P(If(depth > 0, func() T1 { return RandomT1(depth - 1) }, func() T1 { return T1{} })()),
		Noop:        If(depth > 0, func() T1 { return RandomT1(depth - 1) }, func() T1 { return T1{} })(),
		noop:        If(depth > 0, func() T1 { return RandomT1(depth - 1) }, func() T1 { return T1{} })(),
		OtherObject: If(rand.Intn(4) > 1, &T2{}, nil),
		Slice:       slices.Map(make([]int, 4), func(_ int) int { return rand.Intn(42) }),
	}
}

func RandomT1Slice(n int, depth int) []T1 {
	s := make([]T1, n)
	for i := range s {
		s[i] = RandomT1(depth)
	}
	return s
}

type testCase struct {
	name  string
	opts  []Option
	count int
}

func TestPrint(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic: %v", r)
		}
	}()
	tcs := []testCase{
		{
			name:  "default",
			count: 10,
		},
		{
			name:  "max 3",
			opts:  []Option{WithMax(3)},
			count: 10,
		},
		{
			name:  "fields selection",
			opts:  []Option{WithFields("AString", "BInt", "CBool", "Dur", "Time")},
			count: 10,
		},
		{
			name:  "upper headers",
			opts:  []Option{WithUpperHeaders()},
			count: 10,
		},
		{
			name:  "no headers",
			opts:  []Option{WithNoHeaders()},
			count: 10,
		},
		{
			name:  "lower all",
			opts:  []Option{WithLowerValues(), WithLowerHeaders()},
			count: 10,
		},
		{
			name:  "json",
			opts:  []Option{WithJSON()},
			count: 3,
		},
		{
			name: "pretty json",
			opts: []Option{WithJSON(), WithJSONMarshaler(func(v interface{}) ([]byte, error) {
				return json.MarshalIndent(v, "", "  ")
			})},
			count: 2,
		},
		{
			name:  "yaml",
			opts:  []Option{WithYAML()},
			count: 2,
		},
		{
			name:  "ghodss yaml",
			opts:  []Option{WithYAML(), WithYAMLMarshaler(yaml.Marshal)},
			count: 2,
		},
		{
			name: "type formatter",
			opts: []Option{WithTypeFormatter(time.Time{}, func(v interface{}) string {
				return v.(time.Time).Format(time.RFC3339)
			})},
			count: 4,
		},
	}
	for _, v := range tcs {
		t.Run(v.name, func(t *testing.T) {
			fmt.Println()
			data := RandomT1Slice(v.count, 1)
			if err := Print(data, v.opts...); err != nil {
				t.Error(err)
			}
			fmt.Println()
		})
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
