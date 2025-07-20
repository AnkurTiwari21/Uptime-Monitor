package sendgrid

import (
	"bytes"
	"text/template"

	config "github.com/ankur12345678/uptime-monitor/Config"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const emailTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Status Notification</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f9fc;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        margin: 40px auto;
        background-color: #ffffff;
        padding: 30px;
        border-radius: 8px;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
      }
      .header {
        font-size: 22px;
        font-weight: bold;
        color: #333333;
        margin-bottom: 20px;
      }
      .content {
        font-size: 16px;
        color: #555555;
        line-height: 1.6;
      }
      .highlight {
        font-weight: bold;
        color: #007bff;
      }
      .footer {
        margin-top: 30px;
        font-size: 13px;
        color: #999999;
        text-align: center;
      }
      a {
        color: #007bff;
        text-decoration: none;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">Website Status Update</div>
      <div class="content">
        Hello,<br /><br />
        This is an automated message to inform you that the website at
        <a href="{{.WebsiteURL}}" class="highlight">{{.WebsiteURL}}</a> is currently in
        <span class="highlight">{{.Status}}</span> status.<br /><br />
        Please take appropriate action if needed.
      </div>
      <div class="footer">
        &copy; {{.Year}} Uptime Mon8or. All rights reserved.
      </div>
    </div>
  </body>
</html>

`

const plainTextTemplate = `
Website Status Update

Hello,

This is an automated message to inform you that the website at {{.WebsiteURL}} is currently in {{.Status}} status.

Please take appropriate action if needed.

© {{.Year}} Uptime Mon8or. All rights reserved.
`

type EmailData struct {
	WebsiteURL string
	Status     string
	Year       int
}

func RenderHTMLBody(data EmailData) (string, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderPlainTextBody(data EmailData) (string, error) {
	tmpl, err := template.New("plain").Parse(plainTextTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendEmail(toEmail, toName string, cfg *config.Creds, data EmailData) error {

	htmlBody, err := RenderHTMLBody(data)
	if err != nil {
		logger.Error("error in preparing email from html template | err: ", err)
		return err
	}

	plainText, err := RenderPlainTextBody(data)
	if err != nil {
		logger.Error("error in preparing email from plain template | err: ", err)
		return err
	}

	from := mail.NewEmail(cfg.ServiceName, cfg.SendgridFromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, "Webiste Status Update", to, plainText, htmlBody)

	client := sendgrid.NewSendClient(cfg.SendgridApiKey)
	resp, err := client.Send(message)
	if resp != nil {
		logger.Info("RESp: ", resp)
	}

	return err
}
