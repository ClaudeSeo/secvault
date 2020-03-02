package cli

import (
	"fmt"
    "context"

	"github.com/claudeseo/secvault/internal/aws"
)

func List() {
    c := aws.New()
	secrets := c.ListSecrets(context.Background())
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
