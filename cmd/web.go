package cmd

import (
	"github.com/Scorpio69t/gcloc/web"
	"github.com/spf13/cobra"
)

var (
	Port = "8080" // Default port for the web server
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server for gcloc",
	Run: func(cmd *cobra.Command, args []string) {
		if err := web.Start(":8080"); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	webCmd.PersistentFlags().StringVarP(&Port, "port", "p", "8080", "Port to run the web server on")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
