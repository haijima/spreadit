package cmd

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/haijima/spreadit/internal"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func TestNewRootCmd(t *testing.T) {
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "testdata/credentials.json")
	cmd := NewRootCmd(viper.New(), afero.NewMemMapFs())

	assert.Equal(t, "spreadit", cmd.Name())
}

func TestRootCmdExecute(t *testing.T) {
	defer gock.Off()
	gock.New("https://sheets.googleapis.com").Reply(200).JSON(sheets.Spreadsheet{Sheets: []*sheets.Sheet{}})
	gock.New("https://sheets.googleapis.com").Reply(200).JSON(sheets.BatchUpdateSpreadsheetResponse{UpdatedSpreadsheet: &sheets.Spreadsheet{Sheets: []*sheets.Sheet{{Properties: &sheets.SheetProperties{SheetId: 1, Title: "test"}}}}})
	gock.New("https://sheets.googleapis.com").Reply(200).JSON(nil)

	internal.DefaultOptions = []option.ClientOption{option.WithHTTPClient(&http.Client{})}
	cmd := NewRootCmd(viper.New(), afero.NewMemMapFs())

	cmd.SetArgs([]string{"-i", "1X2Y3Z4W5V6U7T8S9R0Q", "-t", "test"})
	err := cmd.ExecuteContext(context.Background())

	assert.NoError(t, err)
}
