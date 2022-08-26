[![Test][test-img]][test]
[![Go Report Card][go-report-img]][go-report]
[![Codecov][codecov-img]][codecov]
[![License: MIT][license-img]][license]

# vulcan-check-catalog

Vulcan Check Catalog Generator

## ⚠️ Alpha status

This tool is under active development and the expected input and output may change.

## Installing

From source code

```sh
# Last release version
go install github.com/adevinta/vulcan-check-catalog/cmd/vulcan-check-catalog@latest

# The main version
go install github.com/adevinta/vulcan-check-catalog/cmd/vulcan-check-catalog@main
```

## Running

Usage:
```sh
usage: vulcan-check-catalog [optional flags] [required flags] <path to checks>

Flags:
  -checktypes-url-list string
        Checktypes URL list. Optional.
  -output string
        Output file path. Optional.
  -registry-url string
        Docker image registry base URL. Required.
  -tag string
        Docker image tag. Optional.
```

---

[test]: https://github.com/adevinta/vulcan-check-catalog/actions/workflows/test.yaml
[test-img]: https://github.com/adevinta/vulcan-check-catalog/actions/workflows/test.yaml/badge.svg
[go-report]: https://goreportcard.com/report/github.com/adevinta/vulcan-check-catalog
[go-report-img]: https://goreportcard.com/badge/github.com/adevinta/vulcan-check-catalog
[codecov]: https://codecov.io/gh/adevinta/vulcan-check-catalog
[codecov-img]: https://codecov.io/gh/adevinta/vulcan-check-catalog/branch/main/graph/badge.svg
[license]: https://github.com/adevinta/vulcan-check-catalog/actions/blob/main/LICENSE
[license-img]: https://img.shields.io/badge/License-MIT-blue.svg
