package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/template"
)

/**
 * Documentation generation and testing.
 * This script looks through all .late files in docs/, generates the final
 * markdown files and runs the documentation-driven test suite to ensure what's
 * documented is exactly how the system works.
 */

func main() {
	docsDir := os.Args[1]

	var lateFiles []string

	entries, _ := ioutil.ReadDir(docsDir)
	for _, file := range entries {
		if path.Ext(file.Name()) == ".late" {
			lateFiles = append(lateFiles, path.Join(docsDir, file.Name()))
		}
	}

	var testCase string
	var expected string
	success := true

	for _, file := range lateFiles {
		fmt.Printf("Rendering %s...", file)
		testCase, expected = parseAndTestDocFile(file)
		t := template.New(testCase)
		ctx := context.New()
		results := t.Render(ctx)

		if len(t.Errors) > 0 {
			success = false
			fmt.Printf("\x1b[31m𝙓\x1b[0m\n")
			fmt.Println("Rendering had the following errors:")
			for _, err := range t.Errors {
				fmt.Printf("\t %s\n", err)
			}

			continue
		}

		if expected == results {
			fmt.Printf("\x1b[32m✓\x1b[0m\n")
		} else {
			success = false
			fmt.Printf("\x1b[31m𝙓\x1b[0m\n")
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(expected, results, false)
			fmt.Printf("Unexpected template result\n\n%s\n", dmp.DiffPrettyText(diffs))
		}

		// TODO
		// If match, write out the resulting docs in docs/generated
		// If not match, track as an error to output at the end
	}

	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func parseAndTestDocFile(filePath string) (string, string) {
	content, _ := ioutil.ReadFile(filePath)

	var testCase []string
	var expected []string
	var inTestCase bool

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if len(line) == 0 {
			// Drop the new-line between the test case and
			// the expected results to ensure the content lines
			// up as expected.
			if inTestCase {
				inTestCase = false
			} else {
				testCase = append(testCase, line)
				expected = append(expected, line)
			}
			continue
		}

		switch line[0] {
		case '>':
			testCase = append(testCase, line[1:])
			inTestCase = true
		case '<':
			expected = append(expected, line[1:])
		default:
			testCase = append(testCase, line)
			expected = append(expected, line)
		}
	}

	return strings.Join(testCase, "\n"), strings.Join(expected, "\n")
}
