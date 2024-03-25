# spreadit

[![CI Status](https://github.com/haijima/spreadit/workflows/CI/badge.svg?branch=main)](https://github.com/haijima/spreadit/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/haijima/spreadit.svg)](https://pkg.go.dev/github.com/haijima/spreadit)
[![Go report](https://goreportcard.com/badge/github.com/haijima/spreadit)](https://goreportcard.com/report/github.com/haijima/spreadit)

`spreadit` is a simple CLI tool to write CSV data to Google Sheets.

## Usage

``` sh
spreadit --id <spreadsheet_id> --title <title> [--file <file>] [--range <range>] [--append]
```

### Examples

``` sh
# specify csv file
spreadit --file data.csv --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet'

# or pipe data
cat data.csv | spreadit --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet'

# or redirect into stdin
spreadit --id 1X2Y3Z4W5V6U7T8S9R0Q --title 'New Sheet' < data.csv
```

### Options

| Option          | Type   | Description                                                | Default |
|-----------------|--------|------------------------------------------------------------|---------|
| `--id` `-i`     | string | Spreadsheet ID                                             |         |
| `--title` `-t`  | string | Sheet title                                                |         |
| `--file` `-f`   | string | CSV file path. If not specified read from stdin (optional) |         |
| `--range` `-r`  | string | Range to write (optional)                                  | "A1"    |
| `--append` `-a` | bool   | Append data to the end of the sheet (optional)             | false   |
| `--format`      | string | Input format. `csv` or `tsv` (optional)                    | "csv"   |
| `--config`      | string | Config file path (optional)                                |         |

### Config file

You can specify the default options in the config file.

``` yaml
# ~/.config/spreadit/.spreadit.yaml
id: 1X2Y3Z4W5V6U7T8S9R0Q
title: New Sheet
file: data.csv
range: A1
append: true
format: csv
```

Config file is searched in the following order:

1. `--config` option
2. `$CURRENT_DIR/.spreadit.yaml`
3. `$XDG_CONFIG_HOME/spreadit/.spreadit.yaml`
4. `$HOME/.config/spreadit/.spreadit.yaml` when `$XDG_CONFIG_HOME` is not set
5. `$HOME/.spreadit.yaml`

[YAML](https://yaml.org/), [JSON](https://www.json.org/json-en.html) or [TOML](https://toml.io/en/) format is supported.

## Requirements

`spreadit` requires the following environment variables to be set:

- `GOOGLE_APPLICATION_CREDENTIALS`: path to the service account key file

See [here](https://cloud.google.com/docs/authentication/getting-started) for more details.

## Install

You can install `spreadit` using the following command:

``` sh
go install github.com/haijima/spreadit@latest
```

MacOS users can install stool using [Homebrew](https://brew.sh/) (See also [haijima/homebrew-tap](http://github.com/haijima/homebrew-tap)):

``` sh
brew install haijima/tap/spreadit
```

or you can download binaries from [Releases](https://github.com/haijima/spreadit/releases).

## License

This tool is licensed under the MIT License. See the [LICENSE](https://github.com/haijima/spreadit/blob/main/LICENSE)
file for details.
