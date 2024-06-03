package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/haijima/cobrax"
	"github.com/haijima/spreadit/internal"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(v *viper.Viper, fs afero.Fs) *cobra.Command {
	rootCmd := cobrax.NewRoot(v)
	rootCmd.Use = "spreadit --id <spreadsheet_id> --title <title> [--file <file>] [--range <range>] [--append]"
	rootCmd.DisableFlagsInUseLine = true
	rootCmd.Short = "Add CSV data to Google Sheets"
	rootCmd.Example = `  spreadit -f data.csv -i 1X2Y3Z4W5V6U7T8S9R0Q -t 'New Sheet'
  cat data.csv | spreadit -i 1X2Y3Z4W5V6U7T8S9R0Q -t 'New Sheet'`
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Colorization settings
		color.NoColor = color.NoColor || v.GetBool("no-color")
		// Set Logger
		lv := cobrax.VerbosityLevel(v)
		l := slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{Level: lv, AddSource: lv < slog.LevelDebug}))
		slog.SetDefault(l)
		cobrax.SetLogger(l)

		return cobrax.RootPersistentPreRunE(cmd, v, fs, args)
	}
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error { return run(cmd, v, fs, args) }

	rootCmd.Flags().StringP("file", "f", "", "The file name to read CSV data from. If not specified, read from stdin.")
	rootCmd.Flags().StringP("id", "i", "", "The ID of the Google Sheets spreadsheet to add the new sheet to")
	rootCmd.Flags().StringP("title", "t", "", "The name of the new sheet to create")
	rootCmd.Flags().String("range", "A1", "The range to append the CSV data to.")
	rootCmd.Flags().BoolP("append", "a", false, "Append the CSV data to the end of the existing sheet instead of creating a new sheet")
	rootCmd.Flags().String("format", "csv", "The format of the data to read. Valid values are 'csv' and 'tsv'")
	rootCmd.Flags().Bool("no_color", false, "disable colorized output")

	return rootCmd
}

func run(cmd *cobra.Command, v *viper.Viper, fs afero.Fs, _ []string) error {
	spreadsheetId := v.GetString("id")
	title := v.GetString("title")
	a1Range := v.GetString("range")
	doesAppend := v.GetBool("append")
	file := v.GetString("file")
	format := v.GetString("format")
	format = strings.ToLower(format)
	if spreadsheetId == "" {
		return fmt.Errorf("spreadsheet ID is required. Use --id or -i")
	}
	if title == "" {
		return fmt.Errorf("title is required. Use --title or -t")
	}
	if format != "csv" && format != "tsv" {
		return fmt.Errorf("invalid format: %s", format)
	}

	ctx := cmd.Context()
	service, err := internal.NewSheetsService(ctx)
	if err != nil {
		return err
	}

	r, err := cobrax.OpenOrStdIn(file, fs, cobrax.WithStdin(cmd.InOrStdin()))
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
			fmt.Sprintf("Read %s data from %s", format, in),
			"Retrieve spreadsheet info",
			"Create a new sheet",
			fmt.Sprintf("Add %s data to the sheet", format),
		},
		spinner.WithWriter(cmd.ErrOrStderr()))
	defer tasks.Close()
	tasks.Start()

	// 1. Read CSV data
	sep := ','
	if format == "tsv" {
		sep = '\t'
	}
	valueRange, err := service.ReadCsv(r, sep)
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
		insertedSheetId, err := service.CreateNewSheet(ctx, spreadsheetId, title)
		if err != nil {
			return err
		}
		sheetId = insertedSheetId
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
	cmd.PrintErrf("Add %s data to Google Sheets successfully!\n", format)
	cmd.PrintErrf("Open %s#gid=%d\n", spreadsheet.SpreadsheetUrl, sheetId)
	cmd.PrintErrln()
	return nil
}
