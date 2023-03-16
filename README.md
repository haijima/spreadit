# spreadit
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