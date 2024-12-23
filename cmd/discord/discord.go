package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/lum8rjack/RSS-Monitor/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	discordWebhook string
	discordRegex   string = `https:\/\/discord.com\/api\/webhooks\/[0-9]+/[0-9a-zA-Z-_]+`
)

type DiscordWebhook struct {
	Webhook string
}

type DiscordMessage struct {
	Content string `json:"content"`
}

// DiscordCmd represents the discord command
var DiscordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Discord webhook notification",
	Long:  `Check the RSS feeds and send the recent posts to a Discord webhook.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get main arguments
		ba, err := utils.GetArgsData(cmd)
		utils.CheckError(err)

		// Discord specific arguments
		webhook, err := cmd.Flags().GetString("webhook")
		utils.CheckError(err)

		// Create the Discord struct
		dwh, err := NewDiscordWebhook(webhook)
		if err != nil {
			ba.Logger.Error(err.Error())
			os.Exit(0)
		}

		// Get the rss feeds
		posts, err := utils.GetRssUpdates(ba.RssFile, ba.Timewindow)
		utils.CheckError(err)

		// Send message
		if len(posts) > 0 {
			ba.Logger.Info("received posts", "number", len(posts))

			// Create message from template
			message, err := utils.GenerateMessage(ba.TemplateFile, posts)
			if err != nil {
				ba.Logger.Error(err.Error(), "template_file", ba.TemplateFile)
				os.Exit(0)
			}
			ba.Logger.Debug("generated message from template file", "template_file", ba.TemplateFile)

			if message == "" {
				ba.Logger.Info("empty message, not sending to webhook")
				os.Exit(0)
			}

			// Create message
			err = dwh.SendWebhook(message)
			if err != nil {
				ba.Logger.Error(err.Error())
				os.Exit(0)
			}

			ba.Logger.Info("successfully sent Discord webhook")
		}

	},
}

func NewDiscordWebhook(webhook string) (DiscordWebhook, error) {
	dw := DiscordWebhook{
		Webhook: webhook,
	}

	if webhook == "" {
		return dw, errors.New("Discord webhook is empty")
	}

	match, err := regexp.MatchString(discordRegex, webhook)
	if err != nil {
		return dw, errors.New("error parsing regex")
	}
	if !match {
		return dw, errors.New("invalid Discord webhook")
	}

	return dw, nil
}

func (d *DiscordWebhook) SendWebhook(message string) error {
	// Create json message
	dm := DiscordMessage{
		Content: message,
	}

	jsonData, err := json.Marshal(dm)
	if err != nil {
		return err
	}

	// Send message
	req, err := http.NewRequest("POST", d.Webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Confirm valid response
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(responseBody))
	}

	return nil
}

func init() {
	// Flags
	DiscordCmd.Flags().StringVarP(&discordWebhook, "webhook", "w", "", "Discord webhook to use")

	// Required flags
	DiscordCmd.MarkFlagRequired("webhook")
}
