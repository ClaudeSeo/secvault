package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	SecretName string
	TagName    string
	ARN        string
}

type SecretsAPI interface {
	DescribeSecret(ctx context.Context, secretId string) *Secret
	GetSecret(ctx context.Context, secretId string) map[string]string
	ListSecrets(ctx context.Context) []Secret
}

func (c *Cloud) DescribeSecret(ctx context.Context, secretId string) *Secret {
	payload := secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	}

	req := c.secrets.DescribeSecretRequest(&payload)
	result, err := req.Send(ctx)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
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
	return secret
}

func (c *Cloud) GetSecret(ctx context.Context, secretId string) map[string]string {
	payload := secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}
	req := c.secrets.GetSecretValueRequest(&payload)
	result, err := req.Send(ctx)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}
	secretMap := map[string]string{}
	json.Unmarshal([]byte(*result.SecretString), &secretMap)
	return secretMap
}

func (c *Cloud) ListSecrets(ctx context.Context) []Secret {
	payload := secretsmanager.ListSecretsInput{}
	req := c.secrets.ListSecretsRequest(&payload)
	result, err := req.Send(ctx)
	if err != nil {
		panic(err)
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

	return secrets
}
