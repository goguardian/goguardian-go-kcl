package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

type testProcessor struct{}

func (s *testProcessor) Initialize(input *kcl.InitializationInput) {
}

func (s *testProcessor) ProcessRecords(input *kcl.ProcessRecordsInput) {
	body, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://localhost"+receiverPort+"/process_records", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (s *testProcessor) LeaseLost(input *kcl.LeaseLostInput) {
}

func (s *testProcessor) ShardEnded(input *kcl.ShardEndedInput) {
}

func (s *testProcessor) ShutdownRequested(input *kcl.ShutdownRequestedInput) {
}
