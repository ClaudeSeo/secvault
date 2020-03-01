package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	SecretName string
	TagName    string
	ARN        string
}

func DescribeSecret(secretId string) *Secret {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal("Unable to load SDK Config, " + err.Error())
	}

	svc := secretsmanager.New(cfg)
	payload := secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	}

	req := svc.DescribeSecretRequest(&payload)
	result, err := req.Send(context.Background())
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

	output := &Secret{
		SecretName: *result.Name,
		ARN:        *result.ARN,
	}

	for _, tag := range result.Tags {
		if *tag.Key == "SecretName" {
			output.TagName = *tag.Value
		}
	}

	return output
}

func GetSecret(secretId string) map[string]string {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal("Unable to load SDK Config, " + err.Error())
	}

	svc := secretsmanager.New(cfg)
	payload := secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	req := svc.GetSecretValueRequest(&payload)
	result, err := req.Send(context.Background())
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

func ListSecret() []Secret {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal("Unable to load SDK Config, " + err.Error())
	}

	svc := secretsmanager.New(cfg)
	payload := secretsmanager.ListSecretsInput{}

	req := svc.ListSecretsRequest(&payload)
	result, err := req.Send(context.Background())
	if err != nil {
		panic(err)
	}

	var secrets []Secret
	for _, s := range result.SecretList {
		sec := &Secret{
			SecretName: *s.Name,
			ARN:        *s.ARN,
		}

		for _, tag := range s.Tags {
			if *tag.Key == "SecretName" {
				sec.TagName = *tag.Value
			}
		}

		secrets = append(secrets, *sec)
	}

	return secrets
}
