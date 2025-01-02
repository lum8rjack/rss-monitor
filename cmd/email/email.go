package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/lum8rjack/RSS-Monitor/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	fromAddress  string
	fromPassword string
	toAddress    string
	subject      string
	smtpHost     string
	smtpPort     int
	htmlEmail    bool
)

// EmailCmd represents the email command
var EmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Email notification",
	Long: `Check the RSS feeds and send the recent posts as an email.
	
	Gmail info: smtp.gmail.com:587`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get main arguments
		ba, err := utils.GetArgsData(cmd)
		utils.CheckError(err)

		// Check arguments
		if smtpPort < 1 || smtpPort > 65535 {
			ba.Logger.Error("invalid smtp port", "smtp_port", smtpPort)
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
			ba.Logger.Info("empty message, not sending email")
			os.Exit(0)
		}

		var body bytes.Buffer

		// Write the headers
		body.Write([]byte("MIME-version: 1.0\n"))
		body.Write([]byte("From: " + fromAddress + "\n"))
		body.Write([]byte("To: " + toAddress + "\n"))
		body.Write([]byte("Subject: " + subject + "\n"))

		if htmlEmail {
			body.Write([]byte("Content-Type: text/html; charset=\"UTF-8\"\n"))
		} else {
			body.Write([]byte("Content-Type: text/plain; charset=\"UTF-8\"\n"))
		}

		// Write the body
		body.Write([]byte(fmt.Sprintf("\n%s", message)))

		// Send email
		to := strings.Split(toAddress, ",")
		auth := smtp.PlainAuth("", fromAddress, fromPassword, smtpHost)
		if fromPassword == "" {
			auth = nil
		}
		err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, fromAddress, to, body.Bytes())
		if err != nil {
			ba.Logger.Error(err.Error())
			os.Exit(0)
		}

		ba.Logger.Info("successfully sent email")

	},
}

func init() {
	// Flags
	EmailCmd.Flags().StringVarP(&fromAddress, "from", "", "", "Sending FROM address")
	EmailCmd.Flags().StringVarP(&fromPassword, "password", "", "", "Sending FROM password (leave empty for no auth)")
	EmailCmd.Flags().StringVarP(&toAddress, "to", "", "", "Sending TO addresses")
	EmailCmd.Flags().StringVarP(&subject, "subject", "", "RSS Monitor", "Email subject")
	EmailCmd.Flags().StringVarP(&smtpHost, "host", "", "", "SMTP host address")
	EmailCmd.Flags().IntVarP(&smtpPort, "port", "", 25, "SMTP host address")
	EmailCmd.Flags().BoolVarP(&htmlEmail, "html", "", false, "Send as an HTML email")

	// Required flags
	EmailCmd.MarkFlagRequired("from")
	EmailCmd.MarkFlagRequired("to")
	EmailCmd.MarkFlagRequired("host")
	EmailCmd.MarkFlagRequired("port")
}
