package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/daveonkels/clinkding/internal/config"
)

type Formatter struct {
	cfg    *config.Config
	writer io.Writer
}

func New(cfg *config.Config) *Formatter {
	return &Formatter{
		cfg:    cfg,
		writer: os.Stdout,
	}
}

func (f *Formatter) Print(format string, args ...interface{}) {
	if f.cfg.Quiet {
		return
	}
	fmt.Fprintf(f.writer, format, args...)
}

func (f *Formatter) Println(format string, args ...interface{}) {
	if f.cfg.Quiet {
		return
	}
	fmt.Fprintf(f.writer, format+"\n", args...)
}

func (f *Formatter) PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *Formatter) Success(format string, args ...interface{}) {
	if f.cfg.Quiet {
		return
	}
	if f.shouldUseColor() {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Fprintf(f.writer, green("✓")+" "+format+"\n", args...)
	} else {
		fmt.Fprintf(f.writer, format+"\n", args...)
	}
}

func (f *Formatter) Error(format string, args ...interface{}) {
	writer := os.Stderr
	if f.shouldUseColor() {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Fprintf(writer, red("Error:")+" "+format+"\n", args...)
	} else {
		fmt.Fprintf(writer, "Error: "+format+"\n", args...)
	}
}

func (f *Formatter) Warning(format string, args ...interface{}) {
	if f.cfg.Quiet {
		return
	}
	if f.shouldUseColor() {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Fprintf(f.writer, yellow("Warning:")+" "+format+"\n", args...)
	} else {
		fmt.Fprintf(f.writer, "Warning: "+format+"\n", args...)
	}
}

func (f *Formatter) Info(format string, args ...interface{}) {
	if f.cfg.Quiet {
		return
	}
	if f.shouldUseColor() {
		blue := color.New(color.FgBlue).SprintFunc()
		fmt.Fprintf(f.writer, blue("ℹ")+" "+format+"\n", args...)
	} else {
		fmt.Fprintf(f.writer, format+"\n", args...)
	}
}

func (f *Formatter) Bold(text string) string {
	if f.shouldUseColor() {
		return color.New(color.Bold).Sprint(text)
	}
	return text
}

func (f *Formatter) Dim(text string) string {
	if f.shouldUseColor() {
		return color.New(color.Faint).Sprint(text)
	}
	return text
}

func (f *Formatter) shouldUseColor() bool {
	if f.cfg.NoColor {
		return false
	}
	if f.cfg.OutputJSON || f.cfg.OutputPlain {
		return false
	}
	// Check if output is a terminal
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func (f *Formatter) IsTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
