package cmd

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/haijima/cobrax"
	"github.com/haijima/spreadit/internal"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(v *viper.Viper, fs afero.Fs) *cobrax.Command {
	rootCmd := cobrax.NewCommand(v, fs)
	rootCmd.Use = "spreadit {--id|-i} <spreadsheet_id> [--file|-f <file>] [--title|-t <title>] [--range|-r <range>] [--append|-a]"
	rootCmd.DisableFlagsInUseLine = true
	rootCmd.Short = "Add CSV data to Google Sheets"
	rootCmd.Example = `  spreadit -f data.csv -i 1X2Y3Z4W5V6U7T8S9R0Q -t 'New Sheet'
  cat data.csv | spreadit -i 1X2Y3Z4W5V6U7T8S9R0Q -t 'New Sheet'`
	rootCmd.Args = cobra.NoArgs
	rootCmd.RunE = func(cmd *cobrax.Command, args []string) error {
		spreadsheetId := cmd.Viper().GetString("id")
		title := cmd.Viper().GetString("title")
		a1Range := cmd.Viper().GetString("range")
		doesAppend := cmd.Viper().GetBool("append")
		file := cmd.Viper().GetString("file")

		ctx := cmd.Context()
		service, err := internal.NewSheetsService(ctx)
		if err != nil {
			return err
		}

		r, err := cmd.OpenOrStdIn(file)
		if err != nil {
			return err
		}
		defer r.Close()

		in := file
		if in == "" {
			in = "stdin"
		}
		tasks := NewTasks(
			[]string{
				fmt.Sprintf("Read CSV data from %s", in),
				"Retrieve spreadsheet info",
				"Create a new sheet",
				"Add CSV data to the sheet",
			},
			spinner.WithWriter(cmd.ErrOrStderr()))
		defer tasks.Close()
		tasks.Start()

		// 1. Read CSV data
		valueRange, err := service.ReadCsv(r)
		if err != nil {
			return err
		}
		tasks.Next()

		// 2. Retrieve spreadsheet
		spreadsheet, err := service.RetrieveSpreadsheet(ctx, spreadsheetId)
		if err != nil {
			return err
		}
		var sheetId int64
		for _, sheet := range spreadsheet.Sheets {
			if sheet.Properties.Title == title {
				sheetId = sheet.Properties.SheetId
				break
			}
		}
		if sheetId != 0 && !doesAppend {
			return fmt.Errorf("a sheet with the title \"%s\" already exists. Please enter another name or use --apend flag", title)
		}
		tasks.Next()

		// 3. Create a new sheet
		if !doesAppend {
			if err := service.CreateNewSheet(ctx, spreadsheetId, title); err != nil {
				return err
			}
			tasks.Next()
		} else {
			tasks.Skip()
		}

		// 4. Add CSV data to sheet
		if err := service.AddCsvToSheet(ctx, spreadsheetId, title, a1Range, valueRange); err != nil {
			return err
		}
		tasks.Next()

		cmd.PrintErrln()
		cmd.PrintErrln("Add CSV data to Google Sheets successfully!")
		cmd.PrintErrf("Open %s#gid=%d\n", spreadsheet.SpreadsheetUrl, sheetId)
		cmd.PrintErrln()
		return nil
	}

	rootCmd.Flags().StringP("file", "f", "", "The file name to read CSV data from. If not specified, read from stdin.")
	rootCmd.Flags().StringP("id", "i", "", "The ID of the Google Sheets spreadsheet to add the new sheet to")
	rootCmd.Flags().StringP("title", "t", "", "The name of the new sheet to create")
	rootCmd.Flags().String("range", "A1", "The range to append the CSV data to.")
	rootCmd.Flags().BoolP("append", "a", false, "Append the CSV data to the end of the existing sheet instead of creating a new sheet")

	return rootCmd
}
