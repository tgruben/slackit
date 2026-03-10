package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// OutputFormat represents the output format type.
type OutputFormat string

const (
	FormatTable OutputFormat = "table"
	FormatJSON  OutputFormat = "json"
	FormatNames OutputFormat = "names"
)

// ParseFormat parses a format string, returning an error for unknown formats.
func ParseFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "table", "":
		return FormatTable, nil
	case "json":
		return FormatJSON, nil
	case "names":
		return FormatNames, nil
	default:
		return "", fmt.Errorf("unknown format %q (use table, json, or names)", s)
	}
}

// Table writes rows in a tab-aligned table.
func Table(w io.Writer, headers []string, rows [][]string) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	tw.Flush()
}

// JSON writes data as indented JSON.
func JSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// Names writes a single column of names, one per line.
func Names(w io.Writer, names []string) {
	for _, name := range names {
		fmt.Fprintln(w, name)
	}
}
