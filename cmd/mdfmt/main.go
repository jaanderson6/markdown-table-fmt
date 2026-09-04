// Command mdfmt reformats the markdown tables in a file or on stdin.
package main

import (
	"flag"
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
	fs := flag.NewFlagSet("mdfmt", flag.ContinueOnError)
	write := fs.Bool("w", false, "rewrite the file in place instead of printing to stdout")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: mdfmt [-w] [file]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) > 1 {
		return fmt.Errorf("usage: mdfmt [-w] [file]")
	}
	if *write && len(rest) == 0 {
		return fmt.Errorf("-w requires a file argument, not stdin")
	}

	var input []byte
	var err error
	if len(rest) == 1 {
		input, err = os.ReadFile(rest[0])
	} else {
		input, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return err
	}

	formatted := mdtable.Format(string(input))

	if *write {
		info, err := os.Stat(rest[0])
		if err != nil {
			return err
		}
		return os.WriteFile(rest[0], []byte(formatted), info.Mode().Perm())
	}

	_, err = os.Stdout.WriteString(formatted)
	return err
}
