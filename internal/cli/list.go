package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/claudeseo/secvault/internal/aws"
)

func List() {
	c := aws.New()
	secrets, err := c.ListSecrets(context.Background())
	if err != nil {
		log.Fatalf("Unable to retrieving secrets\nError: %s", err.Error())
	}

	if len(secrets) == 0 {
		fmt.Println("Not found secrets in AWS Secrets Manager")
		return
	}

	for idx, secret := range secrets {
		fmt.Printf("[%d] Secret\n", idx)
		fmt.Printf("\tSecretName: %s\n", secret.SecretName)
		fmt.Printf("\tTagName: %s\n", secret.TagName)
		fmt.Printf("\tARN: %s\n\n", secret.ARN)
	}
}
