package main

import (
	"log"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

type sampleProcessor struct {
	logger               *log.Logger
	latestSequenceNumber *string
}

func (s *sampleProcessor) Initialize(input *kcl.InitializationInput) {
	s.printInput("Initialize", input)
}

func (s *sampleProcessor) ProcessRecords(input *kcl.ProcessRecordsInput) {
	s.printInput("ProcessRecords", input)
	for _, record := range input.Records {
		s.latestSequenceNumber = &record.SequenceNumber
	}

	err := input.Checkpoint(s.latestSequenceNumber)
	if err != nil {
		s.logger.Printf("Got error %s", err.Error())
	}
}

func (s *sampleProcessor) LeaseLost(input *kcl.LeaseLostInput) {
	s.printInput("LeaseLost", input)
}

func (s *sampleProcessor) ShardEnded(input *kcl.ShardEndedInput) {
	s.printInput("ShardEnded", input)

	// ShardEnded means the shard is closed and the processor should checkpoint
	// with a nil sequence number to indicate that it has successfully processed
	// all records in the shard.
	err := input.Checkpoint(nil)
	if err != nil {
		s.logger.Printf("Got error %s", err.Error())
	}
}

func (s *sampleProcessor) ShutdownRequested(input *kcl.ShutdownRequestedInput) {
	s.printInput("ShutdownRequested", input)

	// Upon ShutdownRequested, the record processor has the option to checkpoint; however,
	// it is not required.
	if s.latestSequenceNumber != nil {
		err := input.Checkpoint(s.latestSequenceNumber)
		if err != nil {
			s.logger.Printf("Got error %s", err.Error())
		}
	}
}

func (s *sampleProcessor) printInput(inputType string, input interface{}) {
	s.logger.Printf("Sample processor received \"%s\", %v", inputType, input)
}
