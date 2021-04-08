package kcl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

type KCLProcess interface {
	Run() error
}

func GetKCLProcess(r RecordProcessor) KCLProcess {
	return &kclProcess{
		recordProcessor: r,
	}
}

type kclProcess struct {
	recordProcessor RecordProcessor
	shardID         string
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

// message represents a status
type statusMessage struct {
	Action      string `json:"action"`
	ResponseFor string `json:"responseFor"`
}

// message represents a checkpoint
type checkpointMessage struct {
	Action         string `json:"action"`
	SequenceNumber string `json:"sequenceNumber"`
}

// Writes a line to the output file. The line is preceeded and followed by a new
// line because other libraries could be writing to the output file as well
// (e.g. some libs might write debugging info to STDOUT) so we would like to
// prevent our lines from being interlaced with other messages so the
// MultiLangDaemon can understand them.  (e.g. '{"action" : "status",
// "responseFor" : "<someAction>"}')
func writeStatus(action string) {
	rawStatus, err := json.Marshal(statusMessage{
		Action:      "status",
		ResponseFor: action,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(rawStatus))
}

func readMessage() (*message, error) {
	reader := bufio.NewReader(os.Stdin)
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, errors.Wrap(err, "failed to read bytes")
	}

	var msg message
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall message")
	}

	return &msg, nil
}

func (k *kclProcess) Run() error {
	shouldExit := false
	for {
		msg, err := readMessage()
		if err != nil {
			if err == io.EOF {
				panic("Unexpected end of file. Exiting!")
			} else {
				return errors.New("failed to read message")
			}
		}

		switch msg.Action {
		case "initialize":
			k.shardID = msg.ShardID
			k.recordProcessor.Initialize(&InitializationInput{
				shardID: k.shardID,
			})

		case "processRecords":
			k.recordProcessor.ProcessRecords(&ProcessRecordsInput{
				Records: msg.Records,
			})

		case "leaseLost":
			k.recordProcessor.LeaseLost(&LeaseLostInput{})

		case "shardEnded":
			k.recordProcessor.ShardEnded(&ShardEndedInput{})
			shouldExit = true

		case "shutdownRequested":
			k.recordProcessor.ShutdownRequested(&ShutdownRequestedInput{
				SequenceNumber: msg.Checkpoint,
			})
			shouldExit = true

		default:
			// TODO: Implement this

		}
		writeStatus(msg.Action)

		if shouldExit {
			return nil
		}
	}
}

func (k *kclProcess) checkpointFunc(sequenceNumber string) error {
	// Write checkpoint and immediately check for acknowledgement

	rawCheckpoint, err := json.Marshal(checkpointMessage{
		Action:         "checkpoint",
		SequenceNumber: sequenceNumber,
	})
	if err != nil {
		return err
	}

	fmt.Printf(string(rawCheckpoint))

	checkpointMsgOutput, err := readMessage()
	if err != nil {
		return errors.Wrap(err, "failed to read message for checkpoint")
	}

	if checkpointMsgOutput.Error != "" {
		return errors.New(fmt.Sprintf("Error %s when checkpointing", checkpointMsgOutput.Error))
	}

	switch checkpointMsgOutput.Action {
	case "checkpoint":
		// successful checkpoint
	default:
		// unsuccessful checkpoint
		return errors.New("Unknown message. Expecting checkpoint message")
	}
	return nil
}
