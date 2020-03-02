package aws

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/secretsmanageriface"
)

type Cloud struct {
	config  *aws.Config
	secrets secretsmanageriface.ClientAPI
}

func New() *Cloud {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal("Unalbe to load SDK config, " + err.Error())
	}

	return &Cloud{
		config:  &cfg,
		secrets: secretsmanager.New(cfg),
	}
}
