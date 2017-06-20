package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	//nolint - exported so can be added to other packages
	RootCmd = &cobra.Command{
		Use:   "breaking [packages]",
		Short: "Find all breaking changes to functions and methods in a Go project, using git",
		RunE:  rootCmd,
	}

	baseBranch    string
	pathChangelog string
)

func init() {
	RootCmd.Flags().StringVar(&baseBranch, "branch", "master", "Base branch to compare code too")
	RootCmd.Flags().StringVar(&pathChangelog, "changelog", "CHANGELOG.md", "Changelog file to add information too")
}

func main() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func rootCmd(cmd *cobra.Command, args []string) (err error) {

	var out string
	if len(args) > 0 {
		for _, dir := range args {
			dir = strings.Trim(dir, ".")
			dir = strings.Trim(dir, "/")
			fmt.Println("-----------------------------------------------")
			fmt.Println(dir)
			fmt.Println("-----------------------------------------------")
			out, err = findMatches(dir, baseBranch)
			if err != nil {
				return
			}
		}
	} else {
		out, err = findMatches("", baseBranch)
		if err != nil {
			return
		}
	}

	//Load changelog, write output
	if _, err := os.Stat(pathChangelog); os.IsNotExist(err) {
		//if the changelog file does not exist simply write to the stdout
		fmt.Println("Could not load changelog file, here are the results:")
		fmt.Println(out)
	} else {
		//if the changelog does exist, insert output into the second line
		fmt.Println("Writing to changelog file")
		InsertStringToFile(pathChangelog, out, 2)
	}

	return
}
