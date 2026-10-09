// Package ses is the AWS SES (v2 API) adapter for lib/email.Sender. The
// endpoint, region and credentials come from the standard AWS environment
// (AWS_REGION, AWS_ENDPOINT_URL for Floci, the Lambda execution role in
// AWS); the From address must be a verified SES identity.
package ses

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lib/email"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// sendTimeout bounds one SendEmail call; Sender.Send carries no context.
const sendTimeout = 10 * time.Second

type Sender struct {
	client *sesv2.Client
	from   string
}

var _ email.Sender = (*Sender)(nil)

// New builds a Sender from the default AWS configuration chain.
func New(ctx context.Context, from string) (*Sender, error) {
	if from == "" {
		return nil, errors.New("ses: from address is required")
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("ses: load aws config: %w", err)
	}
	return &Sender{client: sesv2.NewFromConfig(cfg), from: from}, nil
}

func (s *Sender) Send(msg email.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	body := &types.Body{}
	content := &types.Content{Data: aws.String(msg.Body), Charset: aws.String("UTF-8")}
	if msg.HTML {
		body.Html = content
	} else {
		body.Text = content
	}
	_, err := s.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.from),
		Destination:      &types.Destination{ToAddresses: msg.To},
		Content: &types.EmailContent{Simple: &types.Message{
			Subject: &types.Content{Data: aws.String(msg.Subject), Charset: aws.String("UTF-8")},
			Body:    body,
		}},
	})
	if err != nil {
		return fmt.Errorf("ses: send: %w", err)
	}
	return nil
}
