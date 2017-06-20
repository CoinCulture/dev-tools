package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"
)

// Match all functions and methods in the git diff
var regexFuncs = regexp.MustCompile(`[\-\+]func (\((.*?)\) )?([A-Z]\w*)\(.*?\).*\{`)

// The old and new definitions of a function in the code
type function struct {
	def1 string
	def2 string
}

func findMatches(dir, baseBranch string) (output string, err error) {
	funcs := make(map[string]function)    // all matches in the diff
	breaking := make(map[string]function) // all matches which broke in the diff

	buf := new(bytes.Buffer)
	cmd := exec.Command("git", "diff", baseBranch, dir)
	if dir == "" {
		cmd = exec.Command("git", "diff", baseBranch)
	}
	cmd.Stdout = buf
	cmd.Run()

	matches := regexFuncs.FindAllSubmatch(buf.Bytes(), -1)
	for _, m := range matches {
		def := string(m[0])
		receiver := string(m[2])
		name := string(m[3])
		if receiver != "" {
			name = fmt.Sprintf("(%s).%s", receiver, name)
		}

		funcInfo, ok := funcs[name]
		if ok {
			// if nothings changed, ignore
			// NOTE: we exclude the prefix -/+ from the diff
			if def[1:] == funcInfo.def1[1:] {
				continue
			}

			// if its changed, it should only change once
			if funcInfo.def2 != "" && def[1:] != funcInfo.def2[1:] {
				err = errors.Errorf("%v\n %v\n %v\n %v\n",
					name, def, funcInfo.def1, funcInfo.def2)
				return
			}

			funcInfo.def2 = def
			funcs[name] = funcInfo
			continue
		}

		funcs[name] = function{def1: def}
	}

	// suss out just the breaking changes
	for n, f := range funcs {
		if f.def2 != "" {
			breaking[n] = f
		}
	}

	for n, f := range breaking {
		output += fmt.Sprintf("%v\n%v\n%v\n\n", n, f.def1, f.def2)
	}
	return
}
