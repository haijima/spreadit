# spreadit

[![CI Status](https://github.com/haijima/spreadit/workflows/CI/badge.svg?branch=main)](https://github.com/haijima/spreadit/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/haijima/spreadit.svg)](https://pkg.go.dev/github.com/haijima/spreadit)
[![Go report](https://goreportcard.com/badge/github.com/haijima/spreadit)](https://goreportcard.com/report/github.com/haijima/spreadit)

CLI to add csv data to Google Sheet

## Usage

``` sh
# specify csv file
spreadit -f data.csv --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet'

# or pipe data
cat data.csv | spreadit --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet'
```

## Install

``` sh
go install github.com/haijima/spreadit@latest
```

## License

This tool is licensed under the MIT License. See the [LICENSE](https://github.com/haijima/spreadit/blob/main/LICENSE) file for details.
