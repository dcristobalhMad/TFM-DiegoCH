package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(event events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {

	fmt.Printf("InvocationID: %s\n", event.InvocationID)
	fmt.Printf("DeliveryStreamArn: %s\n", event.DeliveryStreamArn)
	fmt.Printf("Region: %s\n", event.Region)

	var response events.KinesisFirehoseResponse

	for _, record := range event.Records {
		fmt.Printf("RecordID: %s\n", record.RecordID)
		fmt.Printf("ApproximateArrivalTimestamp: %s\n", record.ApproximateArrivalTimestamp)

		// Transform data: ToUpper the data
		var transformedRecord events.KinesisFirehoseResponseRecord
		transformedRecord.RecordID = record.RecordID
		transformedRecord.Result = events.KinesisFirehoseTransformedStateOk
		transformedRecord.Data = []byte(strings.ToUpper(string(record.Data)))

		response.Records = append(response.Records, transformedRecord)
	}

	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}