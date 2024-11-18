/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gcloc/pkg/language"

	"github.com/spf13/cobra"
)

// showLangCmd represents the showLang command
var showLangCmd = &cobra.Command{
	Use:   "showLang",
	Short: "List all supported languages and their extensions",
	Run: func(cmd *cobra.Command, args []string) {
		languages := language.NewDefinedLanguages()
		cmd.Print(languages.GetFormattedString())
	},
}

func init() {
	rootCmd.AddCommand(showLangCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showLangCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showLangCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
