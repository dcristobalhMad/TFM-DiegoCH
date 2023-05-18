package main

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type Record struct {
	Name string `json:"name"`
}

func Handler(request json.RawMessage) (json.RawMessage, error) {
	// Unmarshal the input data to a slice of Record structs
	var records []Record
	if err := json.Unmarshal(request, &records); err != nil {
		return nil, err
	}

	// Transform the name field of each record to lowercase
	for i := range records {
		records[i].Name = strings.ToLower(records[i].Name)
	}

	// Marshal the transformed data back to JSON
	response, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
