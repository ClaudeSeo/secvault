package cli

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/claudeseo/secvault/internal/aws"
)

func Put(secretName string, file string) {
	if secretName == "" {
		log.Fatal("secretName is required. Set --secret-name flag.")
	}

	if file == "" {
		log.Fatal("fileName is required. Set --file flag.")
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	var result map[string]string
	if err := json.Unmarshal(f, &result); err != nil {
		log.Fatal(err.Error())
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := aws.New()
	secrets := c.DescribeSecret(context.Background(), secretName)
	if secrets == nil {
		c.CreateSecret(context.Background(), secretName)
	}

	c.PutSecretValue(context.Background(), secretName, string(data))
}
