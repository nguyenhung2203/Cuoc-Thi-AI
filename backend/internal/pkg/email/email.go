// Package email provides a small SMTP-backed mailer.
//
// When SMTP is not configured (no SMTP_HOST), Send logs the message instead of
// sending it, so the app works in local/dev without a mail server. This mirrors
// the AI service's mock-mode behaviour and keeps flows testable end-to-end.
package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Sender sends plain-text emails.
type Sender struct {
	host string
	port string
	user string
	pass string
	from string
}

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
	if !s.Enabled() {
		log.Printf("email(dev): to=%s subject=%q\n%s", to, subject, body)
		return nil
	}

	msg := buildMessage(s.from, to, subject, body)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}
	return nil
}

func buildMessage(from, to, subject, body string) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	return b.String()
}
