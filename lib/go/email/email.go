package email

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
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

// InlineImage is one inline MIME part referenced from the HTML body via
// `<img src="cid:<ContentID>">`. Inline (CID) images are the standard way
// to embed images in email — Gmail and Outlook render them, while `data:`
// URIs in <img> tags are blocked by most clients.
type InlineImage struct {
	ContentID string
	MIMEType  string
	Data      []byte
}

// Send delivers a plain or HTML message as a single MIME part. Callers that
// need embedded images use SendWithInlineImage instead; this method is
// deliberately untouched so existing callers keep the exact same wire format.
func (c *Client) Send(msg Message) error {
	return c.transmit(msg, c.buildPayload(msg, nil))
}

// SendWithInlineImage sends an HTML message with inline images embedded as
// multipart/related parts. The body must reference each image with
// `<img src="cid:<ContentID>">`. With no images it falls back to the same
// single-part payload as Send, so callers can always pass a non-empty
// message through this method.
func (c *Client) SendWithInlineImage(msg Message, images ...InlineImage) error {
	return c.transmit(msg, c.buildPayload(msg, images))
}

// buildPayload assembles the MIME message payload: a single part when no
// inline images are given (identical to Send's wire format), otherwise
// multipart/related with the HTML part and one base64 part per inline
// image, each carrying a Content-ID for `cid:` references.
func (c *Client) buildPayload(msg Message, images []InlineImage) []byte {
	if len(images) == 0 || !msg.HTML {
		contentType := "text/plain"
		if msg.HTML {
			contentType = "text/html"
		}
		var sb strings.Builder
		c.writeHeaders(&sb, msg, contentType)
		sb.WriteString("\r\n")
		sb.WriteString(msg.Body)
		return []byte(sb.String())
	}

	boundary := "_payssage_multipart_" + strconv.FormatInt(time.Now().UnixNano(), 36)

	var sb strings.Builder
	c.writeHeaders(&sb, msg, "multipart/related; boundary=\""+boundary+"\"")
	sb.WriteString("\r\n")

	sb.WriteString("--" + boundary + "\r\n")
	sb.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(msg.Body)
	sb.WriteString("\r\n")

	for _, img := range images {
		sb.WriteString("--" + boundary + "\r\n")
		fmt.Fprintf(&sb, "Content-Type: %s\r\n", img.MIMEType)
		sb.WriteString("Content-Transfer-Encoding: base64\r\n")
		fmt.Fprintf(&sb, "Content-ID: <%s>\r\n", img.ContentID)
		sb.WriteString("Content-Disposition: inline\r\n")
		sb.WriteString("\r\n")
		sb.WriteString(base64.StdEncoding.EncodeToString(img.Data))
		sb.WriteString("\r\n")
	}
	sb.WriteString("--" + boundary + "--\r\n")

	return []byte(sb.String())
}

// writeHeaders writes the shared RFC 5322 headers (From/To/Subject/MIME).
func (c *Client) writeHeaders(sb *strings.Builder, msg Message, contentType string) {
	fmt.Fprintf(sb, "From: %s\r\n", c.cfg.From)
	fmt.Fprintf(sb, "To: %s\r\n", strings.Join(msg.To, ", "))
	fmt.Fprintf(sb, "Subject: %s\r\n", msg.Subject)
	sb.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(sb, "Content-Type: %s; charset=\"UTF-8\"\r\n", contentType)
}

// transmit performs the SMTP transaction for an already-built message payload.
func (c *Client) transmit(msg Message, payload []byte) error {
	addr := net.JoinHostPort(c.cfg.Host, strconv.Itoa(c.cfg.Port))

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		return fmt.Errorf("email: failed to dial: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		return fmt.Errorf("email: failed to create SMTP client: %w", err)
	}
	defer func() {
		_ = client.Quit()
	}()

	if c.cfg.TLS {
		tlsConfig := &tls.Config{
			ServerName:         c.cfg.Host,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: c.cfg.Insecure, //nolint:gosec
		}
		err = client.StartTLS(tlsConfig)
		if err != nil {
			return fmt.Errorf("email: STARTTLS failed: %w", err)
		}
	}

	if c.cfg.Username != "" || c.cfg.Password != "" {
		auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
		err := client.Auth(auth)
		if err != nil {
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

	_, err = wc.Write(payload)
	if err != nil {
		return fmt.Errorf("email: write failed: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("email: close failed: %w", err)
	}

	return nil
}
