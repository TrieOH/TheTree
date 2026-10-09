// Package sqsjobs is the SQS carrier for lib/jobs: Enqueue sends the job's
// Envelope as the message body. The queue is drained by a worker function
// whose event source is the queue (see lib/jobs/lambdaevents). The endpoint,
// region and credentials come from the standard AWS environment.
package sqsjobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"lib/jobs"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Enqueuer struct {
	client   *sqs.Client
	queueURL string
}

var _ jobs.Enqueuer = (*Enqueuer)(nil)

// New builds an Enqueuer for queueURL from the default AWS configuration.
func New(ctx context.Context, queueURL string) (*Enqueuer, error) {
	if queueURL == "" {
		return nil, errors.New("sqsjobs: queue URL is required")
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqsjobs: load aws config: %w", err)
	}
	return &Enqueuer{client: sqs.NewFromConfig(cfg), queueURL: queueURL}, nil
}

func (e *Enqueuer) Enqueue(ctx context.Context, job jobs.Job) error {
	env, err := jobs.Encode(job)
	if err != nil {
		return err
	}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("sqsjobs: encode envelope: %w", err)
	}
	_, err = e.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(e.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("sqsjobs: send %s: %w", env.Kind, err)
	}
	return nil
}
