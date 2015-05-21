package main

import (
	"fmt"
	"log"
	"os"

	"github.com/darnould/dynamic-dynamodb-lambda/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/darnould/dynamic-dynamodb-lambda/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/service/dynamodb"
)

func main() {
	if len(os.Args[1:]) < 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "<region>", "<table>", "<read/write>", "<up/down>")
		os.Exit(1)
	}

	region := os.Args[1]
	table := os.Args[2]
	capacity := os.Args[3]
	direction := os.Args[4]
	scaleFactor := int64(2)

	svc := dynamodb.New(&aws.Config{Region: region})

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
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
	if capacity == "read" {
		if direction == "up" {
			updatedReadCapacity = currentReadCapacity * scaleFactor
		} else {
			updatedReadCapacity = currentReadCapacity / scaleFactor
		}
	}

	if capacity == "write" {
		if direction == "up" {
			updatedWriteCapacity = currentWriteCapacity * scaleFactor
		} else {
			updatedWriteCapacity = currentWriteCapacity / scaleFactor
		}
	}

	updateTableInput := &dynamodb.UpdateTableInput{
		TableName: aws.String(table),
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

	log.Printf("Read Capacity: %d -> %d for %s\n", currentReadCapacity, updatedReadCapacity, table)
	log.Printf("Write Capacity: %d -> %d for %s\n", currentWriteCapacity, updatedWriteCapacity, table)
}
