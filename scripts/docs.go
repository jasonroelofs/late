package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jasonroelofs/late"
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

	var errors []string
	var testCase string
	var expected string

	for _, file := range lateFiles {
		testCase, expected = parseAndTestDocFile(file)
		t := late.NewTemplate(testCase)
		results := t.Render()

		if expected != results {
			errors = append(errors, fmt.Sprintf("%s did not render as expected", file))
			continue
		}

		// TODO
		// If match, write out the resulting docs in docs/generated
		// If not match, track as an error to output at the end
	}

	if len(errors) > 0 {
		fmt.Println("There were errors rendering documentation")
		for _, err := range errors {
			fmt.Println(err)
		}
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
