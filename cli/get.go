package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/claudeseo/secvalut/internal/aws"
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

func Get(secretId string, outputType string) {
	if secretId == "" {
		log.Fatal("secretId is required. Set --secret-id flag.")
	}

	if outputType != "" &&
		outputType != "kubernetes" &&
		outputType != "json" &&
		outputType != "dotenv" {
		log.Fatal("outputType is not allowd. [" + outputType + "]")
	}

	describe := aws.DescribeSecret(secretId)
	if describe == nil {
		log.Fatal("SecretsManager not found")
	}

	var secretName string
	if describe.TagName != "" {
		secretName = describe.TagName
	} else {
		secretName = describe.SecretName
	}

	secret := aws.GetSecret(secretId)
	switch outputType {
	case "kubernetes":
		makeKubernetesSecret(secretName, secret)
		break
	case "json":
		makeJson(secretName, secret)
		break
	case "dotenv":
		makeDotEnv(secretName, secret)
		break
	default:
		printSecret(secret)
	}
}
