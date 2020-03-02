package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/claudeseo/secvault/internal/aws"
)

func makeFile(fileName string, data string) {
	f, err := os.Create("./" + fileName)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err := f.Write([]byte(data)); err != nil {
		panic(err)
	}

	f.Sync()
}

func makeKubernetesSecret(secretName string, secret map[string]string) {
	template := `apiVersion: v1
kind: Secret
metadata:
  name: %s
type: Opaque
stringData:
`
	s := fmt.Sprintf(template, secretName)

	for k, v := range secret {
		s += fmt.Sprintf("  %s: %s", k, v)
	}

	makeFile(secretName+".yaml", s)
}

func makeDotEnv(secretName string, secret map[string]string) {
	s := ""
	for k, v := range secret {
		s += fmt.Sprintf("%s=%s\n", k, v)
	}
	makeFile(secretName+".env", s)
}

func makeJson(secretName string, secret map[string]string) {
	result, err := json.Marshal(secret)
	if err != nil {
		panic(err)
	}

	makeFile(secretName+".json", string(result))
}

func printSecret(secret map[string]string) {
	for k, v := range secret {
		fmt.Printf("%s = %s\n", k, v)
	}
}

func Get(secretName string, outputType string) {
	if secretName == "" {
		log.Fatal("secretName is required. Set --secret-name flag.")
	}

	if outputType != "" &&
		outputType != "kubernetes" &&
		outputType != "json" &&
		outputType != "dotenv" {
		log.Fatal("outputType is not allowd. [" + outputType + "]")
	}

	c := aws.New()
	describe := c.DescribeSecret(context.Background(), secretName)
	if describe == nil {
		log.Fatal("SecretsManager not found")
	}

	var secName string
	if describe.TagName != "" {
		secName = describe.TagName
	} else {
		secName = describe.SecretName
	}

	secret := c.GetSecret(context.Background(), secretName)
	switch outputType {
	case "kubernetes":
		makeKubernetesSecret(secName, secret)
		break
	case "json":
		makeJson(secName, secret)
		break
	case "dotenv":
		makeDotEnv(secName, secret)
		break
	default:
		printSecret(secret)
	}
}
