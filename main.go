package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.starlark.net/starlark"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s filename\n", os.Args[0])
		os.Exit(1)
	}

	fn := os.Args[1]
	key := "MANIFEST"

	// Check that input file exists. (We will use the info later when re-writing
	// the file, to preserve its permissions.)
	info, err := os.Stat(fn)
	if err != nil {
		fatal(err)
	}

	// Execute the input file
	thread := &starlark.Thread{Name: "whatever"}
	globals, err := starlark.ExecFile(thread, fn, nil, nil)
	if err != nil {
		fatal(err)
	}

	// Sanity check
	if len(globals) != 1 {
		fatal(fmt.Errorf("expected input file to define exactly one global"))
	}

	manifest, ok := globals[key]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: expected input file to define %s\n", key)
		os.Exit(1)
	}

	// Build something which looks... kind of like the file.

	unformatted := fmt.Sprintf("MANIFEST=%v", manifest.String())

	// Format output with Black. Can't use Buildifier because it doesn't expand
	// the (very compact) code at all; just leaves it all on one line.

	// TODO: Do something less dumb than piping to Python (!!) here. Maybe make
	//       our own low-tech pretty Starlark renderer.

	cmd := exec.Command("black", "-q", "--line-length=80", "-")
	cmd.Stdin = strings.NewReader(unformatted)

	var formatted bytes.Buffer
	cmd.Stdout = &formatted
	err = cmd.Run()
	if err != nil {
		fatal(err)
	}

	// Write formatted output back to input file.

	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fatal(err)
		}
	}()

	fmt.Fprint(f, formatted.String())
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
