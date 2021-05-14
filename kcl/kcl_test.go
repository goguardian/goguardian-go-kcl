package kcl

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"log"
	"os"
	"strings"
	"testing"
)

type mockProcessor struct {
	initializeCall        *InitializationInput
	processRecordsCall    *ProcessRecordsInput
	leaseLostCall         *LeaseLostInput
	shardEndedCall        *ShardEndedInput
	shutdownRequestedCall *ShutdownRequestedInput
}

func (m *mockProcessor) Initialize(input *InitializationInput) {
	m.initializeCall = input
}
func (m *mockProcessor) ProcessRecords(input *ProcessRecordsInput) {
	m.processRecordsCall = input
}
func (m *mockProcessor) LeaseLost(input *LeaseLostInput) {
	m.leaseLostCall = input
}
func (m *mockProcessor) ShardEnded(input *ShardEndedInput) {
	m.shardEndedCall = input
}
func (m *mockProcessor) ShutdownRequested(input *ShutdownRequestedInput) {
	m.shutdownRequestedCall = input
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

	if mProcessor.shutdownRequestedCall == nil {
		t.Errorf("expected shutdownRequested to have been called, but it was not")
	}
}

func TestRun_ProcessRecords(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	data := "testData"
	base64Data := base64.StdEncoding.EncodeToString([]byte(data))
	inputLines := `{"action": "processRecords", "records": ` +
		`[{"data": "` + base64Data + `", "partitionKey": "somePartitionKey", "sequenceNumber": "someSequenceNumber"}]}` +
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
	if string(processRecordsCall.Data) != data {
		t.Errorf("expected '%s', but got '%s'", data, string(processRecordsCall.Data))
	}

	if processRecordsCall.PartitionKey != "somePartitionKey" {
		t.Errorf("expected 'somePartitionKey', but got '%s'", processRecordsCall.PartitionKey)
	}

	if processRecordsCall.SequenceNumber != "someSequenceNumber" {
		t.Errorf("expected 'someSequenceNumber', but got '%s'", processRecordsCall.SequenceNumber)
	}

	if mProcessor.shutdownRequestedCall == nil {
		t.Errorf("expected shutdownRequested to have been called, but it was not")
	}
}

func TestRun_LeaseLost(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	inputLines := `{"action": "leaseLost"}` + "\n" + `{"action": "shutdownRequested"}` + "\n"

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
{"action":"status","responseFor":"leaseLost"}

{"action":"status","responseFor":"shutdownRequested"}
`

	output := outputBuffer.String()
	if expectedOutput != output {
		t.Errorf("expected the kclProcess to write '%s', but instead it wrote '%s'", expectedOutput, output)
	}

	if mProcessor.leaseLostCall == nil {
		t.Errorf("expected leaseLost to have been called, but it was not")
	}

	if mProcessor.shutdownRequestedCall == nil {
		t.Errorf("expected shutdownRequested to have been called, but it was not")
	}
}

func TestRun_ShardEnded(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	inputLines := `{"action": "shardEnded"}` + "\n"

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
{"action":"status","responseFor":"shardEnded"}
`

	output := outputBuffer.String()
	if expectedOutput != output {
		t.Errorf("expected the kclProcess to write '%s', but instead it wrote '%s'", expectedOutput, output)
	}

	if mProcessor.shardEndedCall == nil {
		t.Errorf("expected shardEnded to have been called, but it was not")
	}
}

func TestRun_ShutdownRequested(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	inputLines := `{"action": "shutdownRequested"}` + "\n"

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
{"action":"status","responseFor":"shutdownRequested"}
`

	output := outputBuffer.String()
	if expectedOutput != output {
		t.Errorf("expected the kclProcess to write '%s', but instead it wrote '%s'", expectedOutput, output)
	}

	if mProcessor.shutdownRequestedCall == nil {
		t.Errorf("expected shutdownRequested to have been called, but it was not")
	}
}

func TestRun_UnknownAction(t *testing.T) {
	mProcessor := &mockProcessor{}
	outputBuffer := &bytes.Buffer{}
	inputLines := `{"action": "unknownAction"}` + "\n"

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
	if err == nil {
		t.Error("expected a non nil error when action is unknown")
	}
}

func TestCheckpoint(t *testing.T) {
	someSequenceNumber := "123"
	testCases := []struct {
		sequenceNumber *string
		shouldErr      bool
		inputLines     string
		expectedOutput string
	}{
		{
			sequenceNumber: nil,
			inputLines:     `{"action": "checkpoint"}` + "\n",
			expectedOutput: "\n" + `{"action":"checkpoint","sequenceNumber":null}` + "\n",
		},
		{
			sequenceNumber: &someSequenceNumber,
			inputLines:     `{"action": "checkpoint"}` + "\n",
			expectedOutput: "\n" + `{"action":"checkpoint","sequenceNumber":"123"}` + "\n",
		},
		{
			sequenceNumber: &someSequenceNumber,
			inputLines:     `{"action": "checkpoint", "error": "someError"}` + "\n",
			expectedOutput: "\n" + `{"action":"checkpoint","sequenceNumber":"123"}` + "\n",
			shouldErr:      true,
		},
	}

	for _, testCase := range testCases {
		mProcessor := &mockProcessor{}
		outputBuffer := &bytes.Buffer{}

		k := &kclProcess{
			recordProcessor: mProcessor,
			logger:          defaultLogger,
			reader:          bufio.NewReader(strings.NewReader(testCase.inputLines)),
			writer:          bufio.NewWriter(outputBuffer),
		}

		err := k.checkpoint(testCase.sequenceNumber)

		if testCase.shouldErr && err == nil {
			t.Error("expected an error but did not get one")
		}

		if !testCase.shouldErr && err != nil {
			t.Errorf("unexpected error: %+v", err)
		}

		actualOutput := outputBuffer.String()
		if testCase.expectedOutput != actualOutput {
			t.Errorf("expected output '%s' but was '%s'", testCase.expectedOutput, actualOutput)
		}
	}
}

func TestGetKCLProcess_WithNoOptions(t *testing.T) {
	mProcessor := &mockProcessor{}
	processInterface := GetKCLProcess(mProcessor)
	process, ok := processInterface.(*kclProcess)

	if !ok {
		t.Fatal("GetKCLProcess did not return a *kclProcess")
	}
	if process.recordProcessor != mProcessor {
		t.Error("GetKCLProcess did not set recordProcessor correctly")
	}
	if process.logger != defaultLogger {
		t.Error("GetKCLProcess did not set default logger correctly")
	}
}

func TestGetKCLProcess_WithOptions(t *testing.T) {
	mProcessor := &mockProcessor{}
	anyCustomLogger := log.New(os.Stderr, "", log.LUTC)
	processInterface := GetKCLProcess(mProcessor, WithLogger(anyCustomLogger))
	process, ok := processInterface.(*kclProcess)

	if !ok {
		t.Fatal("GetKCLProcess did not return a *kclProcess")
	}
	if process.logger != anyCustomLogger {
		t.Error("GetKCLProcess did not set default logger correctly")
	}
}
