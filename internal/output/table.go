package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Table struct {
	headers []string
	rows    [][]string
}

func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

func (t *Table) Append(row []string) {
	t.rows = append(t.rows, row)
}

func (t *Table) Render() {
	if len(t.headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(t.headers))
	for i, h := range t.headers {
		widths[i] = len(h)
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Check if we should use colors
	useColor := shouldUseColors()

	// Print header with bold/cyan color
	headerRow := make([]string, len(t.headers))
	for i, h := range t.headers {
		padded := padRight(h, widths[i])
		if useColor {
			headerRow[i] = color.New(color.FgCyan, color.Bold).Sprint(padded)
		} else {
			headerRow[i] = padded
		}
	}
	fmt.Println(strings.Join(headerRow, "  "))

	// Print separator with dim color
	separators := make([]string, len(t.headers))
	for i, w := range widths {
		sep := strings.Repeat("─", w)
		if useColor {
			separators[i] = color.New(color.Faint).Sprint(sep)
		} else {
			separators[i] = strings.Repeat("-", w)
		}
	}
	fmt.Println(strings.Join(separators, "  "))

	// Print rows with subtle coloring on first column (ID)
	for _, row := range t.rows {
		rowCells := make([]string, len(t.headers))
		for i := 0; i < len(t.headers); i++ {
			if i < len(row) {
				padded := padRight(row[i], widths[i])
				// Color the first column (usually ID) in green
				if i == 0 && useColor {
					rowCells[i] = color.New(color.FgGreen).Sprint(padded)
				} else {
					rowCells[i] = padded
				}
			} else {
				rowCells[i] = padRight("", widths[i])
			}
		}
		fmt.Println(strings.Join(rowCells, "  "))
	}
}

func shouldUseColors() bool {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// Check if stdout is a terminal
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func FormatTags(tags []string, maxLen int) string {
	if len(tags) == 0 {
		return "-"
	}
	joined := strings.Join(tags, ", ")
	return TruncateString(joined, maxLen)
}
