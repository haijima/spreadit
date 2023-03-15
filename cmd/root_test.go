package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmdExecute(t *testing.T) {
	v := viper.New()
	fs := afero.NewOsFs()
	cmd := NewRootCmd(v, fs)

	cmd.SetArgs([]string{"--id=1nV3vSo0tWNoCzNZnWEkdKnC1Cmg7IxIsYzubR3Fr6ys", "--title=test"})
	err := cmd.Execute()

	assert.NoError(t, err)
}
