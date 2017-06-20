package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	BASE_BRANCH = "master"

	// match all functions and methods in the git diff
	regexFuncs = regexp.MustCompile(`[\-\+]func (\((.*?)\) )?([A-Z]\w*)\(.*?\).*\{`)
)

type Function struct {
	def1 string
	def2 string
}

func findMatches(dir string) {
	funcs := make(map[string]Function)    // all matches in the diff
	breaking := make(map[string]Function) // all matches which broke in the diff

	buf := new(bytes.Buffer)
	cmd := exec.Command("git", "diff", BASE_BRANCH, dir)
	if dir == "" {
		cmd = exec.Command("git", "diff", BASE_BRANCH)
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
				fmt.Println(name)
				fmt.Println(def)
				fmt.Println(funcInfo.def1)
				fmt.Println(funcInfo.def2)
				panic("")
			}

			funcInfo.def2 = def
			funcs[name] = funcInfo
			continue
		}

		funcs[name] = Function{def1: def}
	}

	// suss out just the breaking changes
	for n, f := range funcs {
		if f.def2 != "" {
			breaking[n] = f
		}
	}

	for n, f := range breaking {
		fmt.Println(n)
		fmt.Println(f.def1)
		fmt.Println(f.def2)
		fmt.Println("")
	}

}

func main() {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	if len(args) > 0 {
		for _, dir := range args {
			dir = strings.Trim(dir, ".")
			dir = strings.Trim(dir, "/")
			fmt.Println("-----------------------------------------------")
			fmt.Println(dir)
			fmt.Println("-----------------------------------------------")
			findMatches(dir)
		}
	} else {
		findMatches("")
	}

}
