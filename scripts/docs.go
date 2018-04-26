package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/template"
)

/**
 * Documentation generation and testing.
 * This script looks through all .late files in docs/, generates the final
 * markdown files and runs the documentation-driven test suite to ensure what's
 * documented is exactly how the system works.
 */

type TestDoc struct {
	FilePath string
	Segments []*Segment
	Errors   []string
	Failed   bool
}

type Segment struct {
	IsLiquid bool
	Input    string
	Expected string
	Output   string
}

func (s *Segment) Matches() bool {
	return s.Output == s.Expected
}

func main() {
	docsDir := os.Args[1]

	var lateFiles []string

	entries, _ := ioutil.ReadDir(docsDir)
	for _, file := range entries {
		if path.Ext(file.Name()) == ".late" {
			lateFiles = append(lateFiles, path.Join(docsDir, file.Name()))
		}
	}

	success := true

	for _, file := range lateFiles {
		fmt.Printf("\nRendering %s...", file)
		testDoc := splitDocFile(file)

		// TODO Load up any data files and build up context

		for _, segment := range testDoc.Segments {
			if !segment.IsLiquid {
				continue
			}

			t := template.New(segment.Input)
			ctx := context.New()
			segment.Output = t.Render(ctx)

			if len(t.Errors) == 0 && segment.Matches() {
				printSuccess()
			} else {
				testDoc.Failed = true
				success = false

				printFailure()
				testDoc.Errors = t.Errors
			}
		}

		if !testDoc.Failed {
			continue
		}

		fmt.Println()
		fmt.Printf("%s had the following errors:\n", testDoc.FilePath)

		for _, err := range testDoc.Errors {
			fmt.Println(err)
		}

		for _, segment := range testDoc.Segments {
			if segment.IsLiquid && !segment.Matches() {
				fmt.Println()
				fmt.Printf("  E: %#v\n", segment.Expected)
				fmt.Printf("  G: %#v\n", segment.Output)
			}
		}
	}

	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func printSuccess() {
	fmt.Printf("\x1b[32m✓\x1b[0m")
}

func printFailure() {
	fmt.Printf("\x1b[31m𝙓\x1b[0m")
}

func splitDocFile(filePath string) *TestDoc {
	content, _ := ioutil.ReadFile(filePath)
	testDoc := &TestDoc{FilePath: filePath}

	splitRegex := regexp.MustCompile("(?m)^$")
	removeLeader := regexp.MustCompile(`(?m)^([<>]\s)`)
	parts := splitRegex.Split(string(content), -1)
	segment := &Segment{}

	for _, part := range parts {
		partStr := string(part)
		clean := strings.TrimSpace(partStr)

		if clean[0] == '>' {
			segment.IsLiquid = true
			segment.Input = removeLeader.ReplaceAllString(clean, "")
		} else if clean[0] == '<' {
			segment.Expected = removeLeader.ReplaceAllString(clean, "")
			testDoc.Segments = append(testDoc.Segments, segment)
			segment = &Segment{}
		} else {
			segment.Input = part
			testDoc.Segments = append(testDoc.Segments, segment)
			segment = &Segment{}
		}
	}

	return testDoc
}
