package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/claudeseo/secvault/internal/aws"
	"gopkg.in/yaml.v2"
)

type KubernetesResource struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Type       string            `yaml:"type"`
	StringData map[string]string `yaml:"stringData"`
}

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
	res := KubernetesResource{
		ApiVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		StringData: secret,
	}
	res.Metadata.Name = secretName

	d, err := yaml.Marshal(&res)
	if err != nil {
		log.Fatal("Marshal error: " + err.Error())
	}

	makeFile(secretName+".yaml", string(d))
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

func validateGetParameter(secretName string, fileType string) {
	if secretName == "" {
		log.Fatal("secretName is required. Set --secret-name flag.")
	}

	if fileType != "" &&
		fileType != "kubernetes" &&
		fileType != "json" &&
		fileType != "dotenv" {
		log.Fatal("fileType is not allowd. [" + fileType + "]")
	}
}

func Get(secretName string, fileType string) {
	validateGetParameter(secretName, fileType)

	c := aws.New()
	describe, err := c.DescribeSecret(context.Background(), secretName)
	if err != nil {
		log.Fatalf("Unable to find the %s secrets\nError: %s", secretName, err.Error())
	}

	var secName string
	if describe.TagName != "" {
		secName = describe.TagName
	} else {
		secName = describe.SecretName
	}

	secret, err := c.GetSecret(context.Background(), secretName)
	if err != nil {
		log.Fatalf("Unable to find the %s secrets\nError: %s", secretName, err.Error())
	}
	switch fileType {
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
