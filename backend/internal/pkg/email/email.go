// Package email provides a small SMTP-backed mailer.
//
// When SMTP is not configured (no SMTP_HOST), Send logs the message instead of
// sending it, so the app works in local/dev without a mail server. This mirrors
// the AI service's mock-mode behaviour and keeps flows testable end-to-end.
package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"mime"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Sender sends emails through SMTP.
type Sender struct {
	host string
	port string
	user string
	pass string
	from string
}

const displayName = "ViệcLàm AI"

// NewSender builds a Sender from SMTP settings. host may be empty (dev mode).
func NewSender(host, port, user, pass, from string) *Sender {
	return &Sender{host: host, port: port, user: user, pass: pass, from: from}
}

// Enabled reports whether a real SMTP host is configured.
func (s *Sender) Enabled() bool {
	return strings.TrimSpace(s.host) != ""
}

// Send delivers a plain-text email to a single recipient.
// In dev mode (no host) it logs the message and returns nil.
func (s *Sender) Send(to, subject, body string) error {
	return s.send(to, subject, buildMessage(s.from, to, subject, body, "text/plain; charset=\"UTF-8\""))
}

// SendHTML delivers a multipart email with a plain-text fallback and HTML body.
func (s *Sender) SendHTML(to, subject, textBody, htmlBody string) error {
	boundary := "=_viieclam-ai-otp"
	var b strings.Builder
	b.WriteString("From: " + formatAddress(s.from) + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n\r\n")
	b.WriteString("--" + boundary + "\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
	b.WriteString(textBody + "\r\n")
	b.WriteString("--" + boundary + "\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	b.WriteString(htmlBody + "\r\n")
	b.WriteString("--" + boundary + "--\r\n")
	return s.send(to, subject, b.String())
}

func (s *Sender) send(to, subject, msg string) error {
	if !s.Enabled() {
		log.Printf("email(dev): to=%s subject=%q status=accepted", maskEmail(to), subject)
		if os.Getenv("DEVELOPMENT_VERBOSE_LOG") == "true" {
			log.Printf("email(dev): body=%s", msg)
		}
		return nil
	}
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	conn, err := (&net.Dialer{Timeout: 10 * time.Second}).Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect smtp: %w", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("create smtp client: %w", err)
	}
	defer client.Close()
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}
	if s.user != "" {
		if err := client.Auth(smtp.PlainAuth("", s.user, s.pass, s.host)); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := client.Mail(s.from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(msg)); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func maskEmail(address string) string {
	parts := strings.SplitN(address, "@", 2)
	if len(parts) != 2 || len(parts[0]) < 3 {
		return "***"
	}
	return parts[0][:2] + "***" + parts[0][len(parts[0])-2:] + "@" + parts[1]
}

func formatAddress(address string) string {
	return mime.QEncoding.Encode("UTF-8", displayName) + " <" + address + ">"
}

func buildMessage(from, to, subject, body, contentType string) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: " + contentType + "\r\n\r\n")
	b.WriteString(body)
	return b.String()
}
