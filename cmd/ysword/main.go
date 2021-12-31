package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zlog-fun/ysword/cmd/ysword/internal/project"
)

const release = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "ysword",
	Short:   "ysword: An simple toolkit for Go restful api.",
	Long:    `ysword: An simple toolkit for Go restful api.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
