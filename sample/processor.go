package main

import (
	"log"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

type sampleProcessor struct {
	logger *log.Logger
}

func (s *sampleProcessor) Initialize(input *kcl.InitializationInput) {
	s.printInput("Initialize", input)
}

func (s *sampleProcessor) ProcessRecords(input *kcl.ProcessRecordsInput) {
	s.printInput("ProcessRecords", input)
}

func (s *sampleProcessor) LeaseLost(input *kcl.LeaseLostInput) {
	s.printInput("LeastLost", input)
}

func (s *sampleProcessor) ShardEnded(input *kcl.ShardEndedInput) {
	s.printInput("ShardEnded", input)
}

func (s *sampleProcessor) ShutdownRequested(input *kcl.ShutdownRequestedInput) {
	s.printInput("ShutdownRequested", input)
}

func (s *sampleProcessor) printInput(inputType string, input interface{}) {
	s.logger.Printf("Sample processor received \"%s\", %v", inputType, input)
}