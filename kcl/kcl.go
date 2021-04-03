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

// Writes a line to the output file. The line is preceeded and followed by a new
// line because other libraries could be writing to the output file as well
// (e.g. some libs might write debugging info to STDOUT) so we would like to
// prevent our lines from being interlaced with other messages so the
// MultiLangDaemon can understand them.  (e.g. '{"action" : "status",
// "responseFor" : "<someAction>"}')
func writeStatus(action string) {
	fmt.Printf("\n{\"action\": \"status\", \"responseFor\": \"%s\"}\n", action)
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
			// TODO: Implement this
		case "leaseLost":
			k.recordProcessor.LeaseLost(&LeaseLostInput{})
		case "shardEnded":
			// TODO: Implement this
		case "shutdownRequested":
			// TODO: Implement this
		default:
			// TODO: Implement this
		}
		writeStatus(msg.Action)
	}
}
