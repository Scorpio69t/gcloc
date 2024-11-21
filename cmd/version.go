/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

// Variables to store version and Git info
var (
	Version   = "dev"  // Default value
	GitCommit = "none" // Default value
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gcloc",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Version: ", Version)
		cmd.Println("Git Commit: ", GitCommit)
		cmd.Println("Build Date: ", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// getGitHash gets the git commit hash
func getGitHash() string {
	// get git commit hash
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	hash := strings.TrimSuffix(string(out), "\n")
	return hash
}
