package kcl

type InitializationInput struct{}
type ProcessRecordsInput struct{}
type LeaseLostInput struct{}
type ShardEndedInput struct{}
type ShutdownRequestedInput struct{}

type RecordProcessor interface {
	Initialize(InitializationInput)
	ProcessRecords(ProcessRecordsInput)
	LeaseLost(LeaseLostInput)
	ShardEnded(ShardEndedInput)
	ShutdownRequested(ShutdownRequestedInput)
}
