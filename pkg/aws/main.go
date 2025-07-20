package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ankur12345678/uptime-monitor/jobs"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func LoadAWSConfig(profile, region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func NewClient(cfg aws.Config) *sqs.Client {
	return sqs.NewFromConfig(cfg)
}

func SendMessage(client *sqs.Client, queueURL string, msg *jobs.SQSIncidentEventType) error {
	// Marshal the struct to JSON
	bodyBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal message: %w", err)
		return err
	}

	body := string(bodyBytes)

	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &body,
	})
	if err != nil {
		return err
	}
	fmt.Println("✅ Message sent:", body)
	return nil
}

func ReceiveMessage(client *sqs.Client, queueURL string) (*types.Message, error) {
	resp, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     10,
		VisibilityTimeout:   30,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Messages) == 0 {
		fmt.Println("No messages.")
		return nil, nil
	}
	return &resp.Messages[0], nil
}

func DeleteMessage(client *sqs.Client, queueURL string, receiptHandle *string) error {
	_, err := client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return err
	}
	fmt.Println("🗑️ Message deleted.")
	return nil
}
