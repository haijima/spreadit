package cmd

import (
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
	rootCmd.Short = "Add CSV data to an existing Google Sheets"
	rootCmd.Example = `  spreadit -f data.csv -i 1X2Y3Z4W5V6U7T8S9R0Q -t 'New Sheet'
  cat data.csv | spreadit -i 1X2Y3Z4W5V6U7T8S9R0Q -t 'New Sheet'`
	rootCmd.Args = cobra.NoArgs
	rootCmd.RunE = func(cmd *cobrax.Command, args []string) error {
		r, err := cmd.OpenOrStdIn(cmd.Viper().GetString("file"))
		if err != nil {
			return err
		}
		defer r.Close()
		opt := internal.AddOption{
			SpreadsheetID: cmd.Viper().GetString("id"),
			NewSheetTitle: cmd.Viper().GetString("title"),
			Range:         cmd.Viper().GetString("range"),
			Append:        cmd.Viper().GetBool("append"),
			Debug:         cmd.D,
			Verbose:       cmd.V,
		}
		service, err := internal.NewSheetsService(cmd.Context())
		if err != nil {
			return err
		}
		return service.AddFileDataToNewSheet(cmd.Context(), r, opt)
	}

	rootCmd.Flags().StringP("file", "f", "", "The file name to read CSV data from. If not specified, read from stdin.")
	rootCmd.Flags().StringP("id", "i", "", "The ID of the Google Sheets spreadsheet to add the new sheet to")
	rootCmd.Flags().StringP("title", "t", "", "The name of the new sheet to create")
	rootCmd.Flags().String("range", "A1", "The range to append the CSV data to.")
	rootCmd.Flags().BoolP("append", "a", false, "Append the CSV data to the end of the existing sheet instead of creating a new sheet")

	return rootCmd
}
