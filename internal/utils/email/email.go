package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Config struct {
	Host       string `validate:"required"`
	Port       int    `validate:"required"`
	Username   string `validate:"required"`
	Password   string `validate:"required"`
	DomainName string `validate:"required"`
}

type EmailSender struct {
	cfg Config
}

func NewEmailSender(cfg Config) EmailSender {
	return EmailSender{cfg: cfg}
}

func (es EmailSender) SendPasswordResetEmail(email, token string) error {
	subject := "BOOKTOK: Reset Password"

	body := fmt.Sprintf("Reset code: %s", token)

	err := es.sendEmail(email, subject, body)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (es EmailSender) SendVerificationEmail(email, token string) error {
	subject := "BOOKTOK: Verification Password"

	body := fmt.Sprintf("http://%s/verify-email?token=%s", es.cfg.DomainName, token)

	err := es.sendEmail(email, subject, body)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (es EmailSender) sendEmail(destinationEmail, subject, body string) error {
	auth := smtp.PlainAuth("", es.cfg.Username, es.cfg.Password, es.cfg.Host)

	message := buildEmailMessage(
		es.cfg.Username,
		destinationEmail,
		subject,
		body,
	)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", es.cfg.Host, es.cfg.Port),
		auth,
		es.cfg.Username,
		[]string{destinationEmail},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("send email to %s: %w", destinationEmail, err)
	}

	return nil
}

func buildEmailMessage(sourceEmail, destinationEmail, subject, body string) string {
	var msg strings.Builder

	fmt.Fprintf(&msg, "From: %s\r\n", sourceEmail)
	fmt.Fprintf(&msg, "To: %s\r\n", destinationEmail)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprintf(&msg, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&msg, "Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	fmt.Fprintf(&msg, "\r\n")
	fmt.Fprintf(&msg, "%s", body)

	return msg.String()
}
