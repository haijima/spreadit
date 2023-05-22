package main

import (
	"context"
	"os"

	"github.com/haijima/spreadit/cmd"
	"github.com/mattn/go-colorable"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := cmd.NewRootCmd(viper.New(), afero.NewOsFs())
	rootCmd.SetOut(colorable.NewColorableStdout())
	rootCmd.SetErr(colorable.NewColorableStderr())
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
