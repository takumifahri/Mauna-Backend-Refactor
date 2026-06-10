package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type SMTPMailer struct {
	host      string
	port      string
	username  string
	password  string
	fromEmail string
	fromName  string
}

func NewSMTPMailerFromEnv() *SMTPMailer {
	return &SMTPMailer{
		host:      os.Getenv("SMTP_HOST"),
		port:      envOrDefault("SMTP_PORT", "587"),
		username:  os.Getenv("SMTP_USERNAME"),
		password:  os.Getenv("SMTP_PASSWORD"),
		fromEmail: envOrDefault("SMTP_FROM_EMAIL", os.Getenv("SMTP_USERNAME")),
		fromName:  envOrDefault("SMTP_FROM_NAME", "Mauna"),
	}
}

func (m *SMTPMailer) SendPasswordReset(ctx context.Context, to string, name string, resetToken string, resetURL string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.host == "" || m.port == "" || m.username == "" || m.password == "" || m.fromEmail == "" {
		return fmt.Errorf("smtp configuration is incomplete")
	}

	subject := "Reset Password Mauna"
	body := passwordResetBody(name, resetToken, resetURL)
	message := m.message(to, subject, body)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	return smtp.SendMail(m.host+":"+m.port, auth, m.fromEmail, []string{to}, []byte(message))
}

func (m *SMTPMailer) SendRegistrationOTP(ctx context.Context, to string, name string, otp string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.host == "" || m.port == "" || m.username == "" || m.password == "" || m.fromEmail == "" {
		return fmt.Errorf("smtp configuration is incomplete")
	}

	subject := "Kode Verifikasi Registrasi Mauna"
	body := registrationOTPBody(name, otp)
	message := m.message(to, subject, body)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	return smtp.SendMail(m.host+":"+m.port, auth, m.fromEmail, []string{to}, []byte(message))
}

func (m *SMTPMailer) message(to string, subject string, body string) string {
	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", m.fromName, m.fromEmail),
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=UTF-8",
	}

	var builder strings.Builder
	for key, value := range headers {
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(value)
		builder.WriteString("\r\n")
	}
	builder.WriteString("\r\n")
	builder.WriteString(body)
	return builder.String()
}

func passwordResetBody(name string, resetToken string, resetURL string) string {
	if strings.TrimSpace(name) == "" {
		name = "Mauna user"
	}

	var builder strings.Builder
	builder.WriteString("Halo ")
	builder.WriteString(name)
	builder.WriteString(",\n\nKami menerima permintaan reset password untuk akun Mauna kamu.\n\n")
	if resetURL != "" {
		builder.WriteString("Buka link berikut untuk reset password:\n")
		builder.WriteString(resetURL)
		builder.WriteString("\n\n")
	}
	builder.WriteString("Token reset password kamu:\n")
	builder.WriteString(resetToken)
	builder.WriteString("\n\nToken ini berlaku 15 menit. Abaikan email ini kalau kamu tidak meminta reset password.\n")
	return builder.String()
}

func registrationOTPBody(name string, otp string) string {
	if strings.TrimSpace(name) == "" {
		name = "Mauna user"
	}

	var builder strings.Builder
	builder.WriteString("Halo ")
	builder.WriteString(name)
	builder.WriteString(",\n\nMasukkan kode berikut untuk menyelesaikan registrasi akun Mauna kamu:\n")
	builder.WriteString(otp)
	builder.WriteString("\n\nKode ini berlaku 10 menit. Abaikan email ini kalau kamu tidak membuat akun Mauna.\n")
	return builder.String()
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func EnvIntOrDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
