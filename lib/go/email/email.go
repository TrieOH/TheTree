package email

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"lib/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// sendTimeout bounds one whole SMTP transaction, from dial to the final
// dot, when the caller's context has no earlier deadline. A server that
// stops answering mid-transaction would otherwise hold the job forever.
const sendTimeout = 30 * time.Second

// ErrHeaderInjection is a message whose recipients or subject hold a line
// break, which would let the text write headers of its own (a Bcc, say).
// Subjects carry names people choose, like an event's, so this is checked
// on every send.
var ErrHeaderInjection = errors.New("email: line break in a header value")

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLS      bool
	Insecure bool
}

// Sender is the transport seam: what a Backend depends on to deliver a
// message. Client (SMTP) and ses.Sender (AWS SES) satisfy it, so moving
// between providers is a wiring change.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

var _ Sender = (*Client)(nil)

// Client is the SMTP Sender.
type Client struct {
	cfg   Config
	quota Quota
}

func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

// ErrMonthlyLimit is a send refused because this month's email budget is
// spent. Cloudflare bills past its included quota, so the cap is ours.
var ErrMonthlyLimit = errors.New("email: monthly send limit reached")

// Quota reserves n sends against this month's budget. It returns
// ErrMonthlyLimit when they don't fit, and reserves nothing then.
type Quota interface {
	Reserve(ctx context.Context, n int) error
}

// WithQuota makes every send reserve its recipients against q first.
func (c *Client) WithQuota(q Quota) *Client {
	c.quota = q
	return c
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
func (c *Client) Send(ctx context.Context, msg Message) error {
	return c.deliver(ctx, msg, nil)
}

// SendWithInlineImage sends an HTML message with inline images embedded as
// multipart/related parts. The body must reference each image with
// `<img src="cid:<ContentID>">`. With no images it falls back to the same
// single-part payload as Send, so callers can always pass a non-empty
// message through this method.
func (c *Client) SendWithInlineImage(ctx context.Context, msg Message, images ...InlineImage) error {
	return c.deliver(ctx, msg, images)
}

// deliver checks the headers, builds the payload and runs the SMTP
// transaction under one span. The span holds no recipient or subject:
// both are personal data.
func (c *Client) deliver(ctx context.Context, msg Message, images []InlineImage) (err error) {
	ctx, span := telemetry.StartSpan(ctx, "email.send")
	span.SetAttributes(attribute.Int("email.recipients", len(msg.To)), attribute.Int("email.inline_images", len(images)))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()
	defer func() { recordSend(ctx, len(msg.To), err) }()
	err = checkHeaders(msg)
	if err != nil {
		return err
	}
	if c.quota != nil {
		err = c.quota.Reserve(ctx, len(msg.To))
		if err != nil {
			return err
		}
	}
	return c.transmit(ctx, msg, c.buildPayload(msg, images))
}

func checkHeaders(msg Message) error {
	values := append([]string{msg.Subject}, msg.To...)
	for _, v := range values {
		if strings.ContainsAny(v, "\r\n") {
			return ErrHeaderInjection
		}
	}
	return nil
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
	// RFC 2047: a subject with accents or a dash (most of ours) is sent
	// encoded, so every client shows it as written. Plain ASCII is left as is.
	fmt.Fprintf(sb, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", msg.Subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(sb, "Content-Type: %s; charset=\"UTF-8\"\r\n", contentType)
}

// transmit performs the SMTP transaction for an already-built message payload.
func (c *Client) transmit(ctx context.Context, msg Message, payload []byte) error {
	addr := net.JoinHostPort(c.cfg.Host, strconv.Itoa(c.cfg.Port))

	deadline := time.Now().Add(sendTimeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	dialer := &net.Dialer{Timeout: 10 * time.Second, Deadline: deadline}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("email: failed to dial: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	// The deadline covers every read and write that follows, TLS and the
	// data included.
	err = conn.SetDeadline(deadline)
	if err != nil {
		return fmt.Errorf("email: set deadline: %w", err)
	}

	tlsConfig := &tls.Config{
		ServerName:         c.cfg.Host,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: c.cfg.Insecure, //nolint:gosec
	}
	// Port 465 is SMTPS: TLS from the first byte, no STARTTLS. Cloudflare
	// Email Service only offers this.
	implicitTLS := c.cfg.TLS && c.cfg.Port == 465
	if implicitTLS {
		conn = tls.Client(conn, tlsConfig)
	}

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		return fmt.Errorf("email: failed to create SMTP client: %w", err)
	}
	defer func() {
		_ = client.Quit()
	}()

	if c.cfg.TLS && !implicitTLS {
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
