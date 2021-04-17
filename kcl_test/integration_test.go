package kcl_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/pkg/errors"
)

const (
	localstackEndpoint = "http://localhost:4566"
	testStreamName     = "sample_kinesis_stream"
)

type testClient struct {
	kClient *kinesis.Kinesis
}

func getTestClient() *testClient {
	awsConfig := &aws.Config{
		Endpoint: aws.String(localstackEndpoint),
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     "anyAccessKeyID",
				SecretAccessKey: "anySecretAccessKeyID",
			},
		}),
	}

	return &testClient{
		kClient: kinesis.New(session.Must(session.NewSession(awsConfig))),
	}
}

func (t *testClient) createStream(streamName string, shardCount int64, timeout time.Duration) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		resp, err := t.kClient.DescribeStream(
			&kinesis.DescribeStreamInput{
				StreamName: &streamName,
			},
		)

		// Create it if it doesn't exist
		if err != nil {
			if e, ok := err.(awserr.Error); ok && e.Code() == "ResourceNotFoundException" {
				t.kClient.CreateStream(
					&kinesis.CreateStreamInput{
						ShardCount: &shardCount,
						StreamName: &streamName,
					},
				)
				continue
			}
		}

		if err == nil && *resp.StreamDescription.StreamStatus == "ACTIVE" {
			return nil
		}

		select {
		case <-ticker.C:
			return errors.New("Timed out waiting to create stream")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *testClient) deleteStream(streamName string, timeout time.Duration) error {
	_, err := t.kClient.DeleteStream(
		&kinesis.DeleteStreamInput{
			StreamName: &streamName,
		},
	)

	if err != nil {
		switch e := err.(type) {
		case awserr.Error:
			if e.Code() == "ResourceNotFoundException" {
				return nil // stream already deleted
			}
		default:
			return errors.Wrap(err, "failed to delete stream")
		}
	}

	ticker := time.Tick(timeout)
	for {
		resp, err := t.kClient.ListStreams(
			&kinesis.ListStreamsInput{
				ExclusiveStartStreamName: &streamName,
				Limit:                    aws.Int64(1),
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to list streams")
		}
		if len(resp.StreamNames) == 0 {
			return nil
		}

		select {
		case <-ticker:
			return errors.New("Timed out waiting to delete stream")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *testClient) putRecords(streamName string, records []string) error {
	entries := []*kinesis.PutRecordsRequestEntry{}
	for _, record := range records {
		entries = append(entries, &kinesis.PutRecordsRequestEntry{
			Data:         []byte(record),
			PartitionKey: &record,
		})
	}
	resp, err := t.kClient.PutRecords(
		&kinesis.PutRecordsInput{
			StreamName: &streamName,
			Records:    entries,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to put records")
	}

	if *resp.FailedRecordCount > 0 {
		return errors.New("failed to put some records")
	}

	return nil
}

func TestIntegration(t *testing.T) {
	tClient := getTestClient()
	err := tClient.deleteStream(testStreamName, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = tClient.createStream(testStreamName, 4, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = tClient.putRecords(testStreamName, []string{"alice", "bob", "charlie"})
	if err != nil {
		t.Fatal(err)
	}
}
