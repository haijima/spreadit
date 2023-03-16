package internal

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type AddOption struct {
	SpreadsheetID string
	NewSheetTitle string
	Range         string
	Append        bool
	Debug         *log.Logger
	Verbose       *log.Logger
}

var DefaultOptions []option.ClientOption

type SheetsService struct {
	srv *sheets.Service
}

func NewSheetsService(ctx context.Context) (*SheetsService, error) {
	srv, err := sheets.NewService(ctx, clientOptions()...)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}
	return &SheetsService{srv}, nil
}

func (s *SheetsService) AddFileDataToNewSheet(ctx context.Context, r io.Reader, opt AddOption) error {
	// Read CSV data
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read CSV data from file or stdin: %v", err)
	}

	// Convert CSV data to Google Sheets API value range format
	var values [][]interface{}
	for _, record := range records {
		var row []interface{}
		for _, value := range record {
			row = append(row, value)
		}
		values = append(values, row)
	}
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Create a new sheet
	if !opt.Append {
		req := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{AddSheet: &sheets.AddSheetRequest{Properties: &sheets.SheetProperties{Title: opt.NewSheetTitle}}},
			},
		}
		if _, err := s.srv.Spreadsheets.BatchUpdate(opt.SpreadsheetID, req).Context(ctx).Do(); err != nil {
			return fmt.Errorf("unable to add sheet to spreadsheet: %v", err)
		}
	}

	// Add CSV data to sheet
	if _, err := s.srv.Spreadsheets.Values.Update(opt.SpreadsheetID, fmt.Sprintf("'%s'!%s", opt.NewSheetTitle, opt.Range), valueRange).ValueInputOption("RAW").Context(ctx).Do(); err != nil {
		return fmt.Errorf("unable to add csv data to sheet: %v", err)
	}
	return nil
}

func clientOptions() []option.ClientOption {
	credential := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	return append(DefaultOptions, option.WithCredentialsFile(credential))
}
