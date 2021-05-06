package kcl

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
)

var defaultLogger = log.New(os.Stderr, "", log.LstdFlags)

type KCLProcess interface {
	Run() error
}

// Record format comes from https://github.com/awslabs/amazon-kinesis-client/blob/master/amazon-kinesis-client-multilang/src/main/java/software/amazon/kinesis/multilang/package-info.java
type Record struct {
	Data           []byte `json:"data"`
	PartitionKey   string `json:"partitionKey"`
	SequenceNumber string `json:"sequenceNumber"`
}

// message format comes from https://github.com/awslabs/amazon-kinesis-client/blob/master/amazon-kinesis-client-multilang/src/main/java/software/amazon/kinesis/multilang/package-info.java
type message struct {
	Action     string   `json:"action"`
	ShardID    string   `json:"shardId"`
	Checkpoint string   `json:"checkpoint"`
	Records    []Record `json:"records"`
	Error      string   `json:"error"`
}

// statusMessage represents a status
type statusMessage struct {
	Action      string `json:"action"`
	ResponseFor string `json:"responseFor"`
}

// checkpointMessage represents a checkpoint
type checkpointMessage struct {
	Action         string `json:"action"`
	SequenceNumber string `json:"sequenceNumber"`
}

type kclProcess struct {
	logger          *log.Logger
	recordProcessor RecordProcessor
	shardID         string

	reader *bufio.Reader
	writer *bufio.Writer
}

// Option signifies the type of options that can be passed to the kclProcess
type Option func(*kclProcess)

// WithLogger can be used to specify a logger
func WithLogger(logger *log.Logger) Option {
	return func(k *kclProcess) {
		k.logger = logger
	}
}

func GetKCLProcess(r RecordProcessor, opts ...Option) KCLProcess {
	kclProcess := &kclProcess{
		recordProcessor: r,
		logger:          defaultLogger,

		writer: bufio.NewWriter(os.Stdout),
		reader: bufio.NewReader(os.Stdin),
	}

	for _, opt := range opts {
		opt(kclProcess)
	}

	return kclProcess
}

// Writes a line to the output file. The line is preceeded and followed by a new
// line because other libraries could be writing to the output file as well
// (e.g. some libs might write debugging info to STDOUT) so we would like to
// prevent our lines from being interlaced with other messages so the
// MultiLangDaemon can understand them.  (e.g. '{"action" : "status",
// "responseFor" : "<someAction>"}')
// Similar to https://github.com/awslabs/amazon-kinesis-client-python/blob/master/amazon_kclpy/kcl.py#L31
func (k *kclProcess) writeLine(bytes []byte) error {
	err := k.writer.WriteByte('\n')
	if err != nil {
		return errors.Wrap(err, "failed to write line")
	}

	_, err = k.writer.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to write line")
	}

	err = k.writer.WriteByte('\n')
	if err != nil {
		return errors.Wrap(err, "failed to write line")
	}

	if err := k.writer.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush line")
	}
	return nil
}

func (k *kclProcess) writeStatus(action string) error {
	status, err := json.Marshal(statusMessage{
		Action:      "status",
		ResponseFor: action,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal status message")
	}

	k.logger.Printf("Writing status %s", status)
	err = k.writeLine([]byte(status))
	if err != nil {
		return err
	}

	return nil
}

func (k *kclProcess) readMessage() (*message, error) {
	bytes, err := k.readLine()
	if err != nil {
		return nil, err
	}

	var msg message
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall message")
	}

	return &msg, nil
}

func (k *kclProcess) readLine() ([]byte, error) {
	bytes, err := k.reader.ReadBytes('\n')
	if err != nil {
		return nil, errors.Wrap(err, "failed to read bytes")
	}
	return bytes, nil
}

func (k *kclProcess) Run() error {
	shouldExit := false
	for {
		msg, err := k.readMessage()
		if err != nil {
			return err
		}

		switch msg.Action {
		case "initialize":
			k.shardID = msg.ShardID
			k.recordProcessor.Initialize(&InitializationInput{
				shardID: k.shardID,
			})

		case "processRecords":
			k.recordProcessor.ProcessRecords(&ProcessRecordsInput{
				Records:    msg.Records,
				Checkpoint: k.checkpoint,
			})

		case "leaseLost":
			k.recordProcessor.LeaseLost(&LeaseLostInput{})

		case "shardEnded":
			k.recordProcessor.ShardEnded(&ShardEndedInput{
				Checkpoint: k.checkpoint,
			})
			shouldExit = true

		case "shutdownRequested":
			k.recordProcessor.ShutdownRequested(&ShutdownRequestedInput{
				SequenceNumber: msg.Checkpoint,
				Checkpoint:     k.checkpoint,
			})
			shouldExit = true

		default:
			return errors.New("Unknown message")

		}
		k.writeStatus(msg.Action)

		if shouldExit {
			return nil
		}
	}
}

func (k *kclProcess) checkpoint(sequenceNumber *string) error {
	// Write checkpoint and immediately check for acknowledgement
	var checkpoint []byte
	var err error

	checkpoint = []byte(`{"action":"checkpoint","sequenceNumber":null}`)
	if sequenceNumber != nil {
		checkpoint, err = json.Marshal(
			checkpointMessage{
				Action:         "checkpoint",
				SequenceNumber: *sequenceNumber,
			})
		if err != nil {
			return errors.Wrap(err, "failed to marshal checkpoint")
		}
	}

	k.logger.Printf("Writing checkpoint %s", checkpoint)
	k.writeLine(checkpoint)

	checkpointMsgOutput, err := k.readMessage()
	if err != nil {
		return errors.Wrap(err, "failed to read message for checkpoint")
	}

	if checkpointMsgOutput.Error != "" {
		return errors.Errorf("Error when checkpointing: %s", checkpointMsgOutput.Error)
	}

	switch checkpointMsgOutput.Action {
	case "checkpoint":
		// successful checkpoint
	default:
		// unsuccessful checkpoint
		return errors.Errorf("Unknown message '%s'. Expecting checkpoint message", checkpointMsgOutput.Action)
	}
	return nil
}
