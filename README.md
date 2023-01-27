# Printer

Printer is a simple tool to print a table on the terminal.


## Installation

```bash
go get go.linka.cloud/printer
```

## Usage

The main exported function is `Print` which takes any value as input and prints it on the terminal.

The table is formatted using the struct's `printer` tags.

It is expected to have the following format: `printer:"<column name: string>(,<column order: int>)"`, or `printer:"-"` to ignore the field.

If no `printer` tag is found, the field name is used as the column name.

**Only exported fields are printed.**

## API

```go
package printer // import "go.linka.cloud/printer"


// FUNCTIONS

func Print(v any, opts ...Option) (err error)

// TYPES

type Encoder func(v any) ([]byte, error)

type Format int

const (
	Table Format = iota
	JSON
	YAML
)
func (f Format) String() string

type Option func(*printer)

func WithFields(fields ...string) Option

func WithFormat(format Format) Option

func WithFormatter(fieldName string, fn func(v any) string) Option

func WithJSON() Option

func WithJSONMarshaler(fn Encoder) Option

func WithLowerHeaders() Option

func WithLowerValues() Option

func WithMax(max int) Option

func WithNoHeaders() Option

func WithUpperHeaders() Option

func WithUpperValues() Option

func WithWriter(writer *tabwriter.Writer) Option

func WithYAML() Option

func WithYAMLMarshaler(fn Encoder) Option


```
