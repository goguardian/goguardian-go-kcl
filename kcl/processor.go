package kcl

type (
	InitializationInput struct {
		shardID string
	}
	ProcessRecordsInput struct {
		Records []Record
	}
	ShouldCheckpointInput struct {
		SourceCallType string
	}
	LeaseLostInput         struct{}
	ShardEndedInput        struct{}
	ShutdownRequestedInput struct {
		SequenceNumber string
	}
)

type RecordProcessor interface {
	Initialize(*InitializationInput)
	ProcessRecords(*ProcessRecordsInput)
	ShouldCheckpoint(*ShouldCheckpointInput) (ShouldCheckpoint bool, Checkpoint string)
	LeaseLost(*LeaseLostInput)
	ShardEnded(*ShardEndedInput)
	ShutdownRequested(*ShutdownRequestedInput)
}
