package main

import (
	"log"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

type sampleProcessor struct {
	logger               *log.Logger
	latestSequenceNumber string
}

func (s *sampleProcessor) Initialize(input *kcl.InitializationInput) {
	s.printInput("Initialize", input)
}

func (s *sampleProcessor) ProcessRecords(input *kcl.ProcessRecordsInput) error {
	s.printInput("ProcessRecords", input)
	for _, record := range input.Records {
		s.latestSequenceNumber = record.SequenceNumber
	}
	return input.Checkpoint(s.latestSequenceNumber)
}

func (s *sampleProcessor) LeaseLost(input *kcl.LeaseLostInput) {
	s.printInput("LeastLost", input)
}

func (s *sampleProcessor) ShardEnded(input *kcl.ShardEndedInput) error {
	s.printInput("ShardEnded", input)
	return input.Checkpoint(s.latestSequenceNumber)
}

func (s *sampleProcessor) ShutdownRequested(input *kcl.ShutdownRequestedInput) error {
	s.printInput("ShutdownRequested", input)
	return input.Checkpoint(input.SequenceNumber)
}

func (s *sampleProcessor) printInput(inputType string, input interface{}) {
	s.logger.Printf("Sample processor received \"%s\", %v", inputType, input)
}
