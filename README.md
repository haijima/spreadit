# spreadit

[![CI Status](https://github.com/haijima/spreadit/workflows/CI/badge.svg?branch=main)](https://github.com/haijima/spreadit/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/haijima/spreadit.svg)](https://pkg.go.dev/github.com/haijima/spreadit)
[![Go report](https://goreportcard.com/badge/github.com/haijima/spreadit)](https://goreportcard.com/report/github.com/haijima/spreadit)

CLI to add csv data to Google Sheet

## Synopsis

``` sh
spreadit {--id|-i} <spreadsheet_id> [--file|-f <file>] [--title|-t <title>] [--range|-r <range>] [--append|-a]
```

### Options

| Option     | Shorthand | Type   | Description                                     | Default  |
|------------|-----------|--------|-------------------------------------------------|----------|
| `--id`     | `-i`      | string | Spreadsheet ID                                  |          |
| `--file`   | `-f`      | string | CSV file path. If not specified read from stdin (optional) |          |
| `--title`  | `-t`      | string | Sheet title (optional)                                     | "Sheet1" |
| `--range`  | `-r`      | string | Range to write (optional)                                 | "A1"     |
| `--append` | `-a`      | bool   | Append data to the end of the sheet (optional)             | false    |

### Examples

``` sh
# specify csv file
spreadit --file data.csv --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet'

# or pipe data
cat data.csv | spreadit --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet'

# or redirect into stdin
spreadit --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet' < data.csv
```

## Requirements

`spreadit` requires the following environment variables to be set:

- `GOOGLE_APPLICATION_CREDENTIALS`: path to the service account key file

## Install

You can install `spreadit` using the following command:

``` sh
go install github.com/haijima/spreadit@latest
```

or you can download binaries from [Releases](https://github.com/haijima/spreadit/releases).


## License

This tool is licensed under the MIT License. See the [LICENSE](https://github.com/haijima/spreadit/blob/main/LICENSE)
file for details.
