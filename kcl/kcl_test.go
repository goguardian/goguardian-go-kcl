package kcl

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type mockProcessor struct {
	initializeCall     *InitializationInput
	processRecordsCall *ProcessRecordsInput
}

func (m *mockProcessor) Initialize(input *InitializationInput) {
	m.initializeCall = input
}
func (m *mockProcessor) ProcessRecords(input *ProcessRecordsInput) {
	m.processRecordsCall = input
}
func (m *mockProcessor) LeaseLost(input *LeaseLostInput) {
}
func (m *mockProcessor) ShardEnded(input *ShardEndedInput) {
}
func (m *mockProcessor) ShutdownRequested(input *ShutdownRequestedInput) {
}

func TestWriteStatus(t *testing.T) {
	bytesBuffer := &bytes.Buffer{}
	k := &kclProcess{
		logger: defaultLogger,
		writer: bufio.NewWriter(bytesBuffer),
	}

	err := k.writeStatus("someAction")
	if err != nil {
		t.Error(err)
	}

	expected := "\n" + `{"action":"status","responseFor":"someAction"}` + "\n"
	actual := bytesBuffer.String()

	if expected != actual {
		t.Errorf("expected '%s' but got '%s'", expected, actual)
	}
}

func TestRun_Initialize(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	inputLines := `{"action": "initialize", "shardId": "someShardID"}` +
		"\n" +
		`{"action": "shutdownRequested"}` +
		"\n"

	k := &kclProcess{
		recordProcessor: mProcessor,
		logger:          defaultLogger,
		reader:          bufio.NewReader(strings.NewReader(inputLines)),
		writer:          bufio.NewWriter(outputBuffer),
	}

	finishedRun := make(chan bool)
	var err error
	go func() {
		err = k.Run()
		finishedRun <- true
	}()

	<-finishedRun
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	expectedOutput := `
{"action":"status","responseFor":"initialize"}

{"action":"status","responseFor":"shutdownRequested"}
`

	output := outputBuffer.String()
	if expectedOutput != output {
		t.Errorf("expected the kclProcess to write '%s', but instead it wrote '%s'", expectedOutput, output)
	}

	if mProcessor.initializeCall.shardID != "someShardID" {
		t.Errorf("unexpected shardID from initialize call %s", mProcessor.initializeCall.shardID)
	}
}

func TestRun_ProcessRecords(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	inputLines := `{"action": "processRecords", "records": ` +
		`[{"data": "dGVzdERhdGE=", "partitionKey": "somePartitionKey", "sequenceNumber": "someSequenceNumber"}]}` +
		"\n" +
		`{"action": "shutdownRequested"}` +
		"\n"

	k := &kclProcess{
		recordProcessor: mProcessor,
		logger:          defaultLogger,
		reader:          bufio.NewReader(strings.NewReader(inputLines)),
		writer:          bufio.NewWriter(outputBuffer),
	}

	finishedRun := make(chan bool)
	var err error
	go func() {
		err = k.Run()
		finishedRun <- true
	}()

	<-finishedRun
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	expectedOutput := `
{"action":"status","responseFor":"processRecords"}

{"action":"status","responseFor":"shutdownRequested"}
`

	output := outputBuffer.String()
	if expectedOutput != output {
		t.Errorf("expected the kclProcess to write '%s', but instead it wrote '%s'", expectedOutput, output)
	}

	processRecordsCall := mProcessor.processRecordsCall.Records[0]
	if string(processRecordsCall.Data) != "testData" {
		t.Errorf("expected 'testData', but got '%s'", string(processRecordsCall.Data))
	}

	if processRecordsCall.PartitionKey != "somePartitionKey" {
		t.Errorf("expected 'somePartitionKey', but got '%s'", processRecordsCall.PartitionKey)
	}

	if processRecordsCall.SequenceNumber != "someSequenceNumber" {
		t.Errorf("expected 'someSequenceNumber', but got '%s'", processRecordsCall.SequenceNumber)
	}
}

func TestReadMessage(t *testing.T) {
	testCases := []struct {
		inputLine       string
		expectedMessage *message
		expectErr       bool
	}{
		{
			inputLine:       `{"action": "initialize", "shardId": "someShardID"}` + "\n",
			expectedMessage: &message{Action: "initialize", ShardID: "someShardID"},
			expectErr:       false,
		},
		{
			inputLine: `{"action": "processRecords", "records": ` +
				`[{"data": "dGVzdERhdGE=", "partitionKey": "somePartitionKey", "sequenceNumber": "someSequenceNumber"}]}` + "\n",
			expectedMessage: &message{Action: "processRecords", Records: []Record{
				{
					Data:           []byte("testData"),
					PartitionKey:   "somePartitionKey",
					SequenceNumber: "someSequenceNumber",
				},
			}},
			expectErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("line %s", testCase.inputLine), func(t *testing.T) {
			k := &kclProcess{
				logger: defaultLogger,
				reader: bufio.NewReader(strings.NewReader(testCase.inputLine)),
			}

			actualMessage, err := k.readMessage()
			if err != nil && !testCase.expectErr {
				t.Errorf("readMessage returned an error '%+v'", err)
			}
			if err == nil && testCase.expectErr {
				t.Error("expected readMessage to return an error, but it didn't")
			}

			if !reflect.DeepEqual(testCase.expectedMessage, actualMessage) {
				t.Errorf("expected '%+v', but got '%+v'", testCase.expectedMessage, actualMessage)
			}
		})
	}

}
