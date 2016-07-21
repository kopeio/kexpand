package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	cobraCommand *cobra.Command
}

var rootCommand = RootCmd{
	cobraCommand: &cobra.Command{
		Use: "kexpand",
	},
}

func Execute() {
	goflag.CommandLine.Parse([]string{})
	if err := rootCommand.cobraCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cmd := rootCommand.cobraCommand

	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
}
