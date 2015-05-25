package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/darnould/dynamic-dynamodb-lambda/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/darnould/dynamic-dynamodb-lambda/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/service/dynamodb"
)

type Config struct {
	Region    string
	Table     string
	Unit      string
	Direction string
	Percent   int
}

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "<json-config>")
		os.Exit(1)
	}

	var config Config
	err := json.Unmarshal([]byte(os.Args[1]), &config)
	if err != nil {
		log.Fatal("Can't parse JSON parameters: ", err)
	}

	svc := dynamodb.New(&aws.Config{Region: config.Region})

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(config.Table),
	}

	resp, err := svc.DescribeTable(params)
	if awserr := aws.Error(err); awserr != nil {
		log.Fatal("Error: ", awserr.Code, awserr.Message)
	} else if err != nil {
		log.Fatal("Error: ", err)
	}

	currentReadCapacity := *resp.Table.ProvisionedThroughput.ReadCapacityUnits
	currentWriteCapacity := *resp.Table.ProvisionedThroughput.WriteCapacityUnits

	var updatedReadCapacity int64
	var updatedWriteCapacity int64
	if config.Unit == "read" {
		if config.Direction == "up" {
			updatedReadCapacity = int64(float64(currentReadCapacity) * (float64(config.Percent) / float64(100)))
		} else {
			updatedReadCapacity = int64(float64(currentReadCapacity) / (float64(config.Percent) / float64(100)))
		}
	}

	if config.Unit == "write" {
		if config.Direction == "up" {
			updatedWriteCapacity = int64(float64(currentWriteCapacity) * (float64(config.Percent) / float64(100)))
		} else {
			updatedWriteCapacity = int64(float64(currentWriteCapacity) / (float64(config.Percent) / float64(100)))
		}
	}

	updateTableInput := &dynamodb.UpdateTableInput{
		TableName: aws.String(config.Table),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  &updatedReadCapacity,
			WriteCapacityUnits: &updatedWriteCapacity,
		},
	}

	_, err = svc.UpdateTable(updateTableInput)
	if awserr := aws.Error(err); awserr != nil {
		log.Fatal("Error: ", awserr.Code, awserr.Message)
	} else if err != nil {
		log.Fatal("Error: ", err)
	}

	if config.Unit == "read" {
		log.Printf("Read Capacity: %d -> %d for %s\n", currentReadCapacity, updatedReadCapacity, config.Table)
	} else if config.Unit == "write" {
		log.Printf("Write Capacity: %d -> %d for %s\n", currentWriteCapacity, updatedWriteCapacity, config.Table)
	}
}
