package email

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

// TestSendPayloadSinglePart pins Send's wire format: a single MIME part
// with no multipart/boundary/Content-ID. Existing consumers (identityx,
// certs, signatures) rely on this shape.
func TestSendPayloadSinglePart(t *testing.T) {
	c := NewClient(Config{From: "badges@test.local"})
	payload := string(c.buildPayload(Message{To: []string{"a@b.c"}, Subject: "s", Body: "<p>hi</p>", HTML: true}, nil))

	if !strings.Contains(payload, "Content-Type: text/html; charset=\"UTF-8\"") {
		t.Fatalf("want single-part text/html, got:\n%s", payload)
	}
	if strings.Contains(payload, "multipart/related") || strings.Contains(payload, "Content-ID") {
		t.Fatalf("single-part send must not be multipart, got:\n%s", payload)
	}
	if !strings.Contains(payload, "<p>hi</p>") {
		t.Fatalf("body must be intact, got:\n%s", payload)
	}
}

// TestSendWithInlineImagePayload builds the multipart/related payload and
// verifies the structure email clients need to render the inline image: a
// boundary split between the HTML part and the base64 PNG part, with the
// right Content-ID so `cid:badge-qr` in the HTML resolves.
func TestSendWithInlineImagePayload(t *testing.T) {
	c := NewClient(Config{From: "badges@test.local"})
	qr := []byte("fake-png-bytes")
	payload := string(c.buildPayload(Message{
		To:      []string{"a@b.c"},
		Subject: "s",
		Body:    `<img src="cid:badge-qr">`,
		HTML:    true,
	}, []InlineImage{{ContentID: "badge-qr", MIMEType: "image/png", Data: qr}}))

	for _, want := range []string{
		"Content-Type: multipart/related",
		"Content-ID: <badge-qr>",
		"Content-Type: image/png",
		base64.StdEncoding.EncodeToString(qr),
		`src="cid:badge-qr"`,
		"--", // boundary markers present (opening + closing)
	} {
		if !strings.Contains(payload, want) {
			t.Fatalf("payload missing %q, got:\n%s", want, payload)
		}
	}
	// Closing boundary: a line of "--<boundary>--".
	if !strings.Contains(payload, "--\r\n") || !strings.HasSuffix(strings.TrimSpace(payload), "--") {
		t.Fatalf("want closing boundary, got:\n%s", payload)
	}
}

// TestSendWithInlineImageNoImagesFallsBack — no images ⇒ identical single
// part as Send.
func TestSendWithInlineImageNoImagesFallsBack(t *testing.T) {
	c := NewClient(Config{From: "badges@test.local"})
	withImages := string(c.buildPayload(Message{To: []string{"a@b.c"}, Subject: "s", Body: "x", HTML: true}, nil))
	plain := string(c.buildPayload(Message{To: []string{"a@b.c"}, Subject: "s", Body: "x", HTML: true}, nil))
	if withImages != plain {
		t.Fatalf("no-image path must equal Send payload:\n--- inline ---\n%s\n--- send ---\n%s", withImages, plain)
	}
}

func TestSendRefusesLineBreaksInHeaders(t *testing.T) {
	c := NewClient(Config{Host: "unused.invalid", Port: 25, From: "a@b.c"})
	for _, msg := range []Message{
		{To: []string{"a@b.c"}, Subject: "Event\r\nBcc: victim@x.y"},
		{To: []string{"a@b.c\nBcc: victim@x.y"}, Subject: "s"},
	} {
		err := c.Send(context.Background(), msg)
		if !errors.Is(err, ErrHeaderInjection) {
			t.Errorf("Send(%q) = %v, want ErrHeaderInjection", msg.Subject, err)
		}
	}
}

func TestSubjectIsEncodedWhenNotASCII(t *testing.T) {
	c := NewClient(Config{From: "a@b.c"})
	payload := string(c.buildPayload(Message{To: []string{"a@b.c"}, Subject: "Certificado — Semana Acadêmica"}, nil))
	if !strings.Contains(payload, "Subject: =?UTF-8?q?") {
		t.Fatalf("non-ASCII subject not RFC 2047 encoded:\n%s", payload)
	}
	plain := string(c.buildPayload(Message{To: []string{"a@b.c"}, Subject: "Reset your password"}, nil))
	if !strings.Contains(plain, "Subject: Reset your password\r\n") {
		t.Fatalf("ASCII subject should be left as is:\n%s", plain)
	}
}
