package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	client    *sendgrid.Client
	fromEmail string
	fromName  string
}

func NewEmailService(apiKey, fromEmail, fromName string) *EmailService {
	if apiKey == "" {
		return nil // Email service disabled
	}

	client := sendgrid.NewSendClient(apiKey)
	return &EmailService{
		client:    client,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *EmailService) SendEmail(ctx context.Context, toEmail, toName, subject, htmlContent string) error {
	if s == nil || s.client == nil {
		return nil // Email service not configured
	}

	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("email send failed with status %d: %s", response.StatusCode, response.Body)
	}

	return nil
}

func (s *EmailService) SendApplicationUpdateEmail(ctx context.Context, toEmail, toName string, applicationID uint, status string) error {
	if s == nil {
		return nil
	}

	subject := "Visa Application Status Update"

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #0ea5e9; color: white; padding: 20px; text-align: center; }
			.content { padding: 20px; background-color: #f9fafb; }
			.status { display: inline-block; padding: 5px 15px; border-radius: 5px; font-weight: bold; }
			.status-{{.Status}} { background-color: #dbeafe; color: #1e40af; }
			.footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Visa Help Center</h1>
			</div>
			<div class="content">
				<h2>Application Status Update</h2>
				<p>Dear {{.ToName}},</p>
				<p>Your visa application (ID: {{.ApplicationID}}) status has been updated.</p>
				<p><strong>New Status:</strong> <span class="status status-{{.Status}}">{{.Status}}</span></p>
				<p>Please log in to your account to view more details.</p>
				<p>Best regards,<br>Visa Help Center Team</p>
			</div>
			<div class="footer">
				<p>This is an automated message. Please do not reply.</p>
			</div>
		</div>
	</body>
	</html>
	`

	t := template.Must(template.New("email").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{
		"ToName":        toName,
		"ApplicationID": applicationID,
		"Status":        status,
	}); err != nil {
		return err
	}

	return s.SendEmail(ctx, toEmail, toName, subject, buf.String())
}

func (s *EmailService) SendDocumentRequestEmail(ctx context.Context, toEmail, toName string, applicationID uint, docType string) error {
	if s == nil {
		return nil
	}

	subject := "Document Request - Visa Application"

	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2>Document Request</h2>
			<p>Dear %s,</p>
			<p>We need you to upload the following document for your visa application (ID: %d):</p>
			<p><strong>%s</strong></p>
			<p>Please log in to your account to upload the document.</p>
			<p>Best regards,<br>Visa Help Center Team</p>
		</div>
	</body>
	</html>
	`, toName, applicationID, docType)

	return s.SendEmail(ctx, toEmail, toName, subject, htmlContent)
}
