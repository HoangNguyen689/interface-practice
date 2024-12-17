package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HoangNguyen689/interface-practice/app/genmigration"
	"github.com/HoangNguyen689/interface-practice/app/queuesample"
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "server",
		Short:         "Server components.",
		SilenceErrors: true,
	}

	cmds := []*cobra.Command{
		queuesample.NewCommand(),
		genmigration.NewCommand(),
	}

	for _, cmd := range cmds {
		rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("failed to execute command: %w", err))
	}
}
