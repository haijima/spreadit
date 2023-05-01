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

func (s *SheetsService) RetrieveSpreadsheet(ctx context.Context, spreadsheetId string) (*sheets.Spreadsheet, error) {
	spreadsheet, err := s.srv.Spreadsheets.Get(spreadsheetId).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}
	return spreadsheet, nil
}

func (s *SheetsService) ReadCsv(r io.Reader, sep rune) (*sheets.ValueRange, error) {
	// Read CSV data
	reader := csv.NewReader(r)
	reader.Comma = sep
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read CSV data from file or stdin: %v", err)
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

	return valueRange, nil
}

func (s *SheetsService) CreateNewSheet(ctx context.Context, spreadSheetId, sheetTitle string) (int64, error) {
	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{AddSheet: &sheets.AddSheetRequest{Properties: &sheets.SheetProperties{Title: sheetTitle}}},
		},
		IncludeSpreadsheetInResponse: true,
	}
	res, err := s.srv.Spreadsheets.BatchUpdate(spreadSheetId, req).Context(ctx).Do()
	if err != nil {
		return 0, fmt.Errorf("unable to add sheet to spreadsheet: %v", err)
	}

	for _, sheet := range res.UpdatedSpreadsheet.Sheets {
		if sheet.Properties.Title == sheetTitle {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("created sheet was not found")
}

func (s *SheetsService) AddCsvToSheet(ctx context.Context, spreadSheetId, sheetTitle, a1Range string, valueRange *sheets.ValueRange) error {
	if _, err := s.srv.Spreadsheets.Values.Update(spreadSheetId, fmt.Sprintf("'%s'!%s", sheetTitle, a1Range), valueRange).ValueInputOption("USER_ENTERED").Context(ctx).Do(); err != nil {
		return fmt.Errorf("unable to add csv data to sheet: %v", err)
	}
	return nil
}

func clientOptions() []option.ClientOption {
	credential := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	return append(DefaultOptions, option.WithCredentialsFile(credential))
}
