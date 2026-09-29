package uikit

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

func PrintTable(s Session, header string, rows [][]string) error {
	w := tabwriter.NewWriter(s.Out, tableMinWidth, tableTabWidth, tablePadding, tablePadChar, tableFlags)
	fmt.Fprintln(w, header)
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, columnSep))
	}
	return w.Flush()
}
