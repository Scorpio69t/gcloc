/*
Copyright © 2024 Yang Ruitao yangruitao6@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gcloc [flags] PATH...",
	Short: "A tool for counting source code files and lines in various programming languages",
	Long: `gcloc is a tool for counting source code files and lines in various programming languages.
It supports a variety of programming languages, and can be customized to support more languages.
It is a simple and easy-to-use tool that can help you quickly count the number of source code files and lines in a project.`,

	Args: cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runGCloc,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gcloc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runGCloc(cmd *cobra.Command, args []string) {
	// Do Stuff Here
}
