package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s filename\n", os.Args[0])
		os.Exit(1)
	}

	inner(func(data interface{}) {
		d := data.(map[string]interface{})
		d["prefix"] = "Mutated!"
	})
}

// TODO: Move this to a separate file.

func inner(mutator func(interface{})) {
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
		fatal(fmt.Errorf("expected input file to define %s", key))
	}

	prim, err := toGolang(manifest)
	if err != nil {
		fatal(err)
	}

	if mutator != nil {
		// TODO: Catch panics so the mutator doesn't have to bother.
		mutator(prim)
	}

	star, err := toStarlark(prim)
	if err != nil {
		fatal(err)
	}

	// Build something which looks... kind of like the file.

	unformatted := fmt.Sprintf("MANIFEST=%v", star.String())

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

func toGolang(data starlark.Value) (interface{}, error) {
	switch x := data.(type) {
	case nil:
		// indicates a bug
		return nil, nil

	case starlark.NoneType:
		return nil, nil

	case starlark.Int:
		i, err := strconv.Atoi(x.String())
		if err != nil {
			return i, err
		}
		return i, nil

	case starlark.Bool:
		if x {
			return true, nil
		} else {
			return false, nil
		}

	case starlark.String:
		return string(x), nil

	case *starlark.List:
		r := []interface{}{}

		for i := 0; i < x.Len(); i++ {
			v, err := toGolang(x.Index(i))
			if err != nil {
				return nil, err
			}

			r = append(r, v)
		}

		return r, nil

	case *starlark.Dict:
		r := map[string]interface{}{}

		for _, k := range x.Keys() {
			kk, err := toGolang(k)
			if err != nil {
				return nil, err
			}

			kkk, ok := kk.(string)
			if !ok {
				panic(fmt.Sprintf("expected string, got %T", kk))
			}

			v, _, err := x.Get(k)
			if err != nil {
				return nil, err
			}

			vv, err := toGolang(v)
			if err != nil {
				return nil, err
			}

			r[kkk] = vv
		}

		return r, nil

	default:
		return nil, fmt.Errorf("not supported: %T", x)
	}
}

func toStarlark(data interface{}) (starlark.Value, error) {
	switch x := data.(type) {
	case nil:
		return starlark.None, nil

	case int:
		return starlark.MakeInt(x), nil

	case bool:
		if x {
			return starlark.True, nil
		} else {
			return starlark.False, nil
		}

	case string:
		return starlark.String(x), nil

	case []interface{}:
		elems := []starlark.Value{}

		for _, v := range x {
			v, err := toStarlark(v)
			if err != nil {
				return nil, err
			}

			elems = append(elems, v)
		}

		return starlark.NewList(elems), nil

	case map[string]interface{}:
		r := starlark.NewDict(len(x))

		for k, v := range x {
			vv, err := toStarlark(v)
			if err != nil {
				return nil, err
			}

			// No need to recuse for key; only string is supported.
			r.SetKey(starlark.String(k), vv)
		}

		return r, nil

	default:
		return nil, fmt.Errorf("not supported: %T", x)
	}
}
