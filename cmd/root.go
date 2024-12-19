package cmd

import (
	"os"

	"github.com/lum8rjack/RSS-Monitor/cmd/discord"
	"github.com/lum8rjack/RSS-Monitor/cmd/slack"
	"github.com/spf13/cobra"
)

var (
	DebugLogger  bool
	RssFile      string
	TemplateFile string
	Timewindow   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "RSS-Monitor",
	Short: "Monitor RSS feeds and send recent posts.",
	Long:  `Monitor RSS feeds based on the timeframe provided and send recent posts to different services`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandPallets() {
	rootCmd.AddCommand(discord.DiscordCmd)
	rootCmd.AddCommand(slack.SlackCmd)
}

func init() {
	// Flags
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&DebugLogger, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&RssFile, "rss", "r", "", "File containing the RSS links to scan")
	rootCmd.PersistentFlags().StringVarP(&TemplateFile, "template", "t", "", "Files used as a template for the message to send")
	rootCmd.PersistentFlags().IntVarP(&Timewindow, "timeframe", "f", 24, "Only get articles that were posted in the past number of hours")

	// Required flags
	rootCmd.MarkPersistentFlagRequired("rss")
	rootCmd.MarkPersistentFlagRequired("template")

	// Add subcommands
	addSubcommandPallets()
}
