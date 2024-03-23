package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/haijima/cobrax"
	"github.com/haijima/spreadit/cmd"
	"github.com/mattn/go-colorable"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string
	commit  string
	date    string

	v       *viper.Viper
	rootCmd *cobra.Command
)

func init() {
	cobra.OnInitialize(func() {
		// Colorization settings
		color.NoColor = color.NoColor || v.GetBool("no-color")
		// Set Logger
		lv := cobrax.VerbosityLevel(v)
		l := slog.New(slog.NewTextHandler(rootCmd.ErrOrStderr(), &slog.HandlerOptions{Level: lv, AddSource: lv < slog.LevelDebug}))
		slog.SetDefault(l)
		cobrax.SetLogger(l)
	})
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelError)
	v = viper.NewWithOptions(viper.WithLogger(slog.Default()))
	fs := afero.NewOsFs()
	v.SetFs(fs)
	rootCmd = cmd.NewRootCmd(v, fs)
	rootCmd.Version = cobrax.VersionFunc(version, commit, date)
	rootCmd.SetOut(colorable.NewColorableStdout())
	rootCmd.SetErr(colorable.NewColorableStderr())
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		if slog.Default().Enabled(rootCmd.Context(), slog.LevelDebug) {
			slog.Error(fmt.Sprintf("%+v", err))
		} else {
			slog.Error(err.Error())
		}
		os.Exit(1)
	}
}
