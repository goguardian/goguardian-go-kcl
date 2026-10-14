package integration_tests

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/pkg/errors"
)

const localstackEndpoint = "http://localhost:4566"

type localKinesis struct {
	kClient *kinesis.Client
}

func GetLocalKinesisClient() (*localKinesis, error) {
	awsConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("anyAccessKeyID", "anySecretAccessKeyID", ""),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load aws config")
	}

	err = waitForLocalstack(30 * time.Second)
	if err != nil {
		errors.New("localstack is not running yet")
	}

	return &localKinesis{
		kClient: kinesis.NewFromConfig(awsConfig, func(o *kinesis.Options) {
			o.BaseEndpoint = aws.String(localstackEndpoint)
		}),
	}, nil
}

func (t *localKinesis) CreateStream(streamName string, shardCount int32, timeout time.Duration) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		resp, err := t.kClient.DescribeStream(
			context.Background(),
			&kinesis.DescribeStreamInput{
				StreamName: &streamName,
			},
		)
		// Create it if it doesn't exist
		if err != nil {
			var notFound *types.ResourceNotFoundException
			if errors.As(err, &notFound) {
				t.kClient.CreateStream(
					context.Background(),
					&kinesis.CreateStreamInput{
						ShardCount: &shardCount,
						StreamName: &streamName,
					},
				)
				continue
			}
		}

		if err == nil && resp.StreamDescription.StreamStatus == types.StreamStatusActive {
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

func (t *localKinesis) DeleteStream(streamName string, timeout time.Duration) error {
	_, err := t.kClient.DeleteStream(
		context.Background(),
		&kinesis.DeleteStreamInput{
			StreamName: &streamName,
		},
	)
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil // stream already deleted
		}
		return errors.Wrap(err, "failed to delete stream")
	}

	ticker := time.Tick(timeout)
	for {
		resp, err := t.kClient.ListStreams(
			context.Background(),
			&kinesis.ListStreamsInput{
				ExclusiveStartStreamName: &streamName,
				Limit:                    aws.Int32(1),
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to list streams")
		}

		found := false
		for _, name := range resp.StreamNames {
			if name == streamName {
				found = true
			}
		}
		if !found {
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

func (t *localKinesis) PutRecords(streamName string, records []string) error {
	entries := []types.PutRecordsRequestEntry{}
	for _, record := range records {
		record := record
		entries = append(entries, types.PutRecordsRequestEntry{
			Data:         []byte(record),
			PartitionKey: &record,
		})
	}

	resp, err := t.kClient.PutRecords(
		context.Background(),
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
