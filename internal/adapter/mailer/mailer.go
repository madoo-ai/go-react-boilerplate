// Package mailer sends the verification and password-reset emails over SMTP.
// When SMTP_HOST is empty the link is logged, so local development works
// without a mail server.
package mailer

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"html"
	"log/slog"
	"mime"
	"net"
	"net/smtp"
	"strings"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	AppName  string
	Secure   bool
}

type Mailer struct {
	cfg Config
	log *slog.Logger
}

func New(cfg Config, log *slog.Logger) *Mailer {
	if cfg.Port == 0 {
		cfg.Port = 1025
	}
	if cfg.From == "" {
		cfg.From = "hello@example.com"
	}
	return &Mailer{cfg: cfg, log: log}
}

func (m *Mailer) SendVerification(ctx context.Context, to, name, url string) error {
	subject := "Verify your email address"
	heading := fmt.Sprintf("Welcome, %s", name)
	body := "Confirm your email address to finish setting up your account."
	action := "Verify email"
	footer := "If you did not create this account, you can ignore this email."
	text := fmt.Sprintf("Welcome, %s.\n\nVerify your email address: %s\n", name, url)
	return m.send(ctx, to, subject, htmlEmail(m.cfg.AppName, heading, body, action, url, footer), text, url)
}

func (m *Mailer) SendPasswordReset(ctx context.Context, to, name, url string) error {
	subject := "Reset your password"
	heading := fmt.Sprintf("Reset your password, %s", name)
	body := "We received a request to reset your password. This link expires in one hour."
	action := "Choose a new password"
	footer := "If you did not request a reset, no action is needed."
	text := fmt.Sprintf("Reset your password, %s.\n\nChoose a new password: %s\n", name, url)
	return m.send(ctx, to, subject, htmlEmail(m.cfg.AppName, heading, body, action, url, footer), text, url)
}

func (m *Mailer) send(ctx context.Context, to, subject, htmlBody, textBody, fallbackURL string) error {
	if strings.TrimSpace(m.cfg.Host) == "" {
		m.log.Info("mail disabled; use this link", "to", to, "subject", subject, "url", fallbackURL)
		return nil
	}

	boundary, err := randomBoundary()
	if err != nil {
		return err
	}
	message := buildMessage(m.cfg.From, m.cfg.AppName, to, subject, boundary, htmlBody, textBody)

	addr := net.JoinHostPort(m.cfg.Host, fmt.Sprintf("%d", m.cfg.Port))
	var auth smtp.Auth
	if m.cfg.User != "" {
		auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)
	}

	if err := m.deliver(addr, auth, m.envelopeFrom(), to, message); err != nil {
		m.log.Error("mail send failed", "to", to, "subject", subject, "err", err)
		m.log.Info("mail fallback link", "to", to, "url", fallbackURL)
		return err
	}
	_ = ctx
	m.log.Info("mail sent", "to", to, "subject", subject)
	return nil
}

func (m *Mailer) deliver(addr string, auth smtp.Auth, from, to string, message []byte) error {
	if m.cfg.Secure {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.cfg.Host, MinVersion: tls.VersionTLS12})
		if err != nil {
			return err
		}
		client, err := smtp.NewClient(conn, m.cfg.Host)
		if err != nil {
			return err
		}
		defer client.Close()
		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
		if err := client.Mail(from); err != nil {
			return err
		}
		if err := client.Rcpt(to); err != nil {
			return err
		}
		writer, err := client.Data()
		if err != nil {
			return err
		}
		if _, err := writer.Write(message); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}
		return client.Quit()
	}
	return smtp.SendMail(addr, auth, from, []string{to}, message)
}

func (m *Mailer) envelopeFrom() string {
	value := strings.TrimSpace(m.cfg.From)
	if start := strings.Index(value, "<"); start >= 0 {
		if end := strings.Index(value[start:], ">"); end >= 0 {
			return strings.TrimSpace(value[start+1 : start+end])
		}
	}
	return value
}

func buildMessage(from, appName, to, subject, boundary, htmlBody, textBody string) []byte {
	var builder strings.Builder
	builder.WriteString("From: " + from + "\r\n")
	builder.WriteString("To: " + to + "\r\n")
	builder.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n")
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n\r\n")
	builder.WriteString("--" + boundary + "\r\n")
	builder.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	builder.WriteString(textBody + "\r\n")
	builder.WriteString("--" + boundary + "\r\n")
	builder.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
	builder.WriteString(htmlBody + "\r\n")
	builder.WriteString("--" + boundary + "--\r\n")
	_ = appName
	return []byte(builder.String())
}

func htmlEmail(appName, heading, body, action, url, footer string) string {
	escapedName := html.EscapeString(appName)
	escapedHeading := html.EscapeString(heading)
	escapedBody := html.EscapeString(body)
	escapedAction := html.EscapeString(action)
	escapedURL := html.EscapeString(url)
	escapedFooter := html.EscapeString(footer)
	return `<!doctype html>
<html><body style="background:#f6f9fc;font-family:Inter,-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;margin:0;padding:40px 12px">
<div style="background:#ffffff;border:1px solid #e3e8ee;border-radius:12px;margin:0 auto;max-width:520px;padding:32px">
<p style="color:#533afd;font-size:13px;font-weight:600;letter-spacing:.08em;margin:0 0 24px;text-transform:uppercase">` + escapedName + `</p>
<h1 style="color:#0d253d;font-size:26px;font-weight:300;letter-spacing:-.26px;line-height:1.2;margin:0 0 12px">` + escapedHeading + `</h1>
<p style="color:#273951;font-size:15px;font-weight:300;line-height:1.6;margin:0 0 24px">` + escapedBody + `</p>
<p style="margin:0 0 24px"><a href="` + escapedURL + `" style="background:#533afd;border-radius:9999px;color:#ffffff;display:inline-block;font-size:15px;font-weight:500;padding:10px 20px;text-decoration:none">` + escapedAction + `</a></p>
<p style="color:#64748d;font-size:12px;line-height:1.6;margin:0 0 24px">If the button does not work, copy this link into your browser:<br><span style="color:#533afd;word-break:break-all">` + escapedURL + `</span></p>
<hr style="border-color:#e3e8ee;margin:0 0 16px">
<p style="color:#64748d;font-size:12px;line-height:1.6;margin:0">` + escapedFooter + `</p>
</div></body></html>`
}

func randomBoundary() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
