package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLS      bool
	Insecure bool
}

type Client struct {
	cfg Config
}

func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

type Message struct {
	To      []string
	Subject string
	Body    string
	HTML    bool
}

func (c *Client) Send(msg Message) error {
	addr := net.JoinHostPort(c.cfg.Host, fmt.Sprintf("%d", c.cfg.Port))

	contentType := "text/plain"
	if msg.HTML {
		contentType = "text/html"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "From: %s\r\n", c.cfg.From)
	fmt.Fprintf(&sb, "To: %s\r\n", strings.Join(msg.To, ", "))
	fmt.Fprintf(&sb, "Subject: %s\r\n", msg.Subject)
	sb.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(&sb, "Content-Type: %s; charset=\"UTF-8\"\r\n", contentType)
	sb.WriteString("\r\n")
	sb.WriteString(msg.Body)

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("email: failed to dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		return fmt.Errorf("email: failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	if c.cfg.TLS {
		tlsConfig := &tls.Config{
			ServerName:         c.cfg.Host,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: c.cfg.Insecure,
		}
		err = client.StartTLS(tlsConfig)
		if err != nil {
			return fmt.Errorf("email: STARTTLS failed: %w", err)
		}
	}

	if c.cfg.Username != "" || c.cfg.Password != "" {
		auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("email: auth failed: %w", err)
		}
	}

	err = client.Mail(c.cfg.From)
	if err != nil {
		return fmt.Errorf("email: MAIL FROM failed: %w", err)
	}

	for _, to := range msg.To {
		err = client.Rcpt(to)
		if err != nil {
			return fmt.Errorf("email: RCPT TO failed: %w", err)
		}
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: DATA failed: %w", err)
	}

	_, err = wc.Write([]byte(sb.String()))
	if err != nil {
		return fmt.Errorf("email: write failed: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("email: close failed: %w", err)
	}

	return nil
}
