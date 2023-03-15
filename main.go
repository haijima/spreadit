package main

import (
	"context"
	"os"

	"github.com/haijima/spreadit/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := cmd.NewRootCmd(viper.New(), afero.NewOsFs())
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
