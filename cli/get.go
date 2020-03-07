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

type kubernetesMetaData struct {
	Name string `yaml:"name"`
}

type KubernetesResource struct {
	ApiVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   kubernetesMetaData `yaml:"metadata"`
	Type       string             `yaml:"type"`
	StringData map[string]string  `yaml:"stringData"`
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
		Metadata: kubernetesMetaData{
			Name: secretName,
		},
		Type:       "Opaque",
		StringData: secret,
	}

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
