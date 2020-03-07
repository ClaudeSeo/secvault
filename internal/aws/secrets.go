package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	SecretName string
	TagName    string
	ARN        string
}

func (c *Cloud) DescribeSecret(ctx context.Context, secretId string) (*Secret, error) {
	payload := secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	}

	req := c.secrets.DescribeSecretRequest(&payload)
	result, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	secret := &Secret{
		SecretName: *result.Name,
		ARN:        *result.ARN,
	}

	for _, tag := range result.Tags {
		if *tag.Key == "SecretName" {
			secret.TagName = *tag.Value
		}

	}

	return secret, nil
}

func (c *Cloud) GetSecret(ctx context.Context, secretId string) (map[string]string, error) {
	payload := secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	req := c.secrets.GetSecretValueRequest(&payload)
	result, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	secretMap := map[string]string{}
	json.Unmarshal([]byte(*result.SecretString), &secretMap)
	return secretMap, nil
}

func (c *Cloud) ListSecrets(ctx context.Context) ([]Secret, error) {
	payload := secretsmanager.ListSecretsInput{}
	req := c.secrets.ListSecretsRequest(&payload)
	result, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}

	var secrets []Secret
	for _, v := range result.SecretList {
		s := &Secret{
			SecretName: *v.Name,
			ARN:        *v.ARN,
		}

		for _, tag := range v.Tags {
			if *tag.Key == "SecretName" {
				s.TagName = *tag.Value
			}
		}

		secrets = append(secrets, *s)
	}

	return secrets, nil
}

func (c *Cloud) CreateSecret(ctx context.Context, secretName string) error {
	payload := secretsmanager.CreateSecretInput{
		SecretString: aws.String("{\"created\": \"success\"}"),
		Name:         aws.String(secretName),
		Description:  aws.String("Created by secvault"),
		Tags: []secretsmanager.Tag{
			{
				Key:   aws.String("SecretName"),
				Value: aws.String(secretName),
			},
		},
	}
	req := c.secrets.CreateSecretRequest(&payload)
	_, err := req.Send(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cloud) PutSecretValue(ctx context.Context, secretName string, data string) error {
	payload := secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretName),
		SecretString: aws.String(data),
	}

	req := c.secrets.PutSecretValueRequest(&payload)
	_, err := req.Send(ctx)
	if err != nil {
		return err
	}

	return nil
}
