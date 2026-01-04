package output

import (
	"fmt"
	"os"
	"strings"
)

func PrintPlain(records [][]string) {
	for _, record := range records {
		fmt.Println(strings.Join(record, "\t"))
	}
}

func PrintPlainLine(fields ...string) {
	_, _ = fmt.Fprintln(os.Stdout, strings.Join(fields, "\t"))
}
