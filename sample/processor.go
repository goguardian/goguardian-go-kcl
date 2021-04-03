package main

import (
	"github.com/goguardian/goguardian-go-kcl/kcl"
)

type sampleProcessor struct {
}

func (s *sampleProcessor) Initialize(input *kcl.InitializationInput) {
}

func (s *sampleProcessor) ProcessRecords(input *kcl.ProcessRecordsInput) {
}

func (s *sampleProcessor) LeaseLost(input *kcl.LeaseLostInput) {
}

func (s *sampleProcessor) ShardEnded(input *kcl.ShardEndedInput) {
}

func (s *sampleProcessor) ShutdownRequested(input *kcl.ShutdownRequestedInput) {
}
