This folder contains the documentation and acceptance test suite for Late.

In each `.late` file, you'll see segments of late code prefixed with `>` followed by
expected output prefixed with `<`. As part of the test suite, Late pulls these sections
out of the documentation and builds two documents, running one through the engine and comparing
the result with the second.

When documentation requires multiple files to function properly, the extra files should exist
in a similarly named directory. The test suite will explicitly look for a `[test-name]/data.json`
file for data to make available to the templates.
