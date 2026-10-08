// Package awssm is the AWS Secrets Manager adapter for lib/secrets. One
// secret holds one Backend's values as a flat JSON object
// ({"HMAC_SECRET": "...", ...}). The endpoint, region and credentials come
// from the standard AWS environment (AWS_REGION, AWS_ENDPOINT_URL for Floci,
// the Lambda execution role in AWS), so nothing here is environment-specific.
package awssm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"lib/secrets"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// Store reads and writes one secret by id (name or ARN).
type Store struct {
	client *secretsmanager.Client
	id     string
}

var _ secrets.Store = (*Store)(nil)

// New builds a Store from the default AWS configuration chain.
func New(ctx context.Context, id string) (*Store, error) {
	if id == "" {
		return nil, errors.New("awssm: secret id is required")
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("awssm: load aws config: %w", err)
	}
	return &Store{client: secretsmanager.NewFromConfig(cfg), id: id}, nil
}

func (s *Store) Load(ctx context.Context) (map[string]string, error) {
	out, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: aws.String(s.id)})
	if err != nil {
		return nil, fmt.Errorf("awssm: get %s: %w", s.id, err)
	}
	values := make(map[string]string)
	err = json.Unmarshal([]byte(aws.ToString(out.SecretString)), &values)
	if err != nil {
		return nil, fmt.Errorf("awssm: secret %s is not a flat JSON object of strings: %w", s.id, err)
	}
	return values, nil
}

// Put replaces the secret's value, creating the secret on first use.
func (s *Store) Put(ctx context.Context, values map[string]string) error {
	body, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("awssm: encode: %w", err)
	}
	_, err = s.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(s.id),
		SecretString: aws.String(string(body)),
	})
	if _, notFound := errors.AsType[*types.ResourceNotFoundException](err); notFound {
		_, err = s.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(s.id),
			SecretString: aws.String(string(body)),
		})
	}
	if err != nil {
		return fmt.Errorf("awssm: put %s: %w", s.id, err)
	}
	return nil
}
