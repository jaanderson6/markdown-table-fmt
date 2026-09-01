// Command mdfmt reformats the markdown tables in a file or on stdin.
package main

import (
	"fmt"
	"io"
	"os"

	mdtable "github.com/jaanderson6/markdown-table-fmt"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "mdfmt:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var input []byte
	var err error

	switch len(args) {
	case 0:
		input, err = io.ReadAll(os.Stdin)
	case 1:
		input, err = os.ReadFile(args[0])
	default:
		return fmt.Errorf("usage: mdfmt [file]")
	}
	if err != nil {
		return err
	}

	_, err = os.Stdout.WriteString(mdtable.Format(string(input)))
	return err
}
