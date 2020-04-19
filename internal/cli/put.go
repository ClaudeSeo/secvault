package cli

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/claudeseo/secvault/internal/aws"
)

func loadDataFromJson(data []byte) []byte {
	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err.Error())
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err.Error())
	}

	return data
}

func loadDataFromDotenv(data []byte) []byte {
	var result map[string]string

	result = make(map[string]string)
	for _, s := range strings.Split(string(data), "\n") {
		d := strings.Split(s, "=")
		if len(d) != 2 {
			continue
		}

		result[d[0]] = d[1]
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err.Error())
	}

	return data
}

func validatePutParameter(secretName string, file string, fileType string) {
	if secretName == "" {
		log.Fatal("secretName is required. Set --secret-name flag.")
	}

	if file == "" {
		log.Fatal("fileName is required. Set --file flag.")
	}

	if fileType != "" && fileType != "json" && fileType != "dotenv" {
		log.Fatal("fileType is not allowd. [" + fileType + "]")
	}
}

func Put(secretName string, file string, fileType string) {
	validatePutParameter(secretName, file, fileType)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	var data []byte
	if fileType == "json" {
		data = loadDataFromJson(f)
	} else if fileType == "dotenv" {
		data = loadDataFromDotenv(f)
	}

	c := aws.New()
	secrets, _ := c.DescribeSecret(context.Background(), secretName)
	if secrets == nil {
		c.CreateSecret(context.Background(), secretName)
	}

	c.PutSecretValue(context.Background(), secretName, string(data))
}
