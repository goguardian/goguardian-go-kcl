package kcl_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
		record := record
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

func waitForLocalstack(timeout time.Duration) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return errors.New("timed out waiting for localstack to start")
		default:
		}

		resp, err := http.Get(fmt.Sprintf("%s/health", localstackEndpoint))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read localstack health")
		}

		health := struct {
			Services map[string]string `json:"services"`
		}{}
		err = json.Unmarshal(body, &health)
		if err != nil {
			return errors.Wrap(err, "failed to parse localstack health")
		}

		if status, found := health.Services["kinesis"]; found && status == "running" {
			return nil
		}
	}
}

func TestIntegration(t *testing.T) {
	log.Println("waiting for localstack to start")
	err := waitForLocalstack(30 * time.Second)
	if err != nil {
		t.Fatal("localstack is not running yet")
	}
	tClient := getTestClient()

	log.Println("deleting kinesis stream if present")
	err = tClient.deleteStream(testStreamName, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("creating new kinesis stream")
	err = tClient.createStream(testStreamName, 4, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("putting records in kinesis stream")
	err = tClient.putRecords(testStreamName, []string{"alice", "bob", "charlie"})
	if err != nil {
		t.Fatal(err)
	}
}
