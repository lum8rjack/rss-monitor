package slack

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
	slackWebhook       string
	slackServiceRegex  string = "https://hooks.slack.com/services/T[A-Z0-9]+/B[A-Z0-9]+/[A-Za-z0-9]{23,25}"
	slackWorkflowRegex string = "https://hooks.slack.com/workflows/T[A-Z0-9]+/A[A-Z0-9]+/[0-9]{17,19}/[A-Za-z0-9]{23,25}"
)

type SlackWebhook struct {
	Webhook string
}

type SlackMessage struct {
	Text string `json:"text"`
}

// SlackCmd represents the slack command
var SlackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Slack webhook notification",
	Long:  `Check the RSS feeds and send the recent posts to a Slack webhook.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get main arguments
		ba, err := utils.GetArgsData(cmd)
		utils.CheckError(err)

		// Create the Slack struct
		swh, err := NewSlackWebhook(slackWebhook)
		if err != nil {
			ba.Logger.Error(err.Error())
			os.Exit(0)
		}

		// Get the rss feeds
		posts, err := utils.GetRssUpdates(ba.RssFile, ba.Timewindow)
		utils.CheckError(err)

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
		err = swh.SendWebhook(message)
		if err != nil {
			ba.Logger.Error(err.Error())
			os.Exit(0)
		}

		ba.Logger.Info("successfully sent Slack webhook")

	},
}

func NewSlackWebhook(webhook string) (SlackWebhook, error) {
	dw := SlackWebhook{
		Webhook: webhook,
	}

	if webhook == "" {
		return dw, errors.New("Slack webhook is empty")
	}

	// Check slack service webhook
	match, err := regexp.MatchString(slackServiceRegex, webhook)
	if err != nil {
		return dw, errors.New("error parsing regex")
	}

	if match {
		return dw, nil
	}

	// Check slack workflow webhook
	match, err = regexp.MatchString(slackWorkflowRegex, webhook)
	if err != nil {
		return dw, errors.New("error parsing regex")
	}

	if !match {
		return dw, errors.New("invalid Slack webhook")
	}

	return dw, nil
}

func (d *SlackWebhook) SendWebhook(message string) error {
	// Create json message
	dm := SlackMessage{
		Text: message,
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
	SlackCmd.Flags().StringVarP(&slackWebhook, "webhook", "w", "", "Slack webhook to use")

	// Required flags
	SlackCmd.MarkFlagRequired("webhook")
}
