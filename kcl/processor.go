package kcl

type CheckpointFunc = func(string) error

type (
	InitializationInput struct {
		shardID string
	}
	ProcessRecordsInput struct {
		Records    []Record
		Checkpoint CheckpointFunc
	}
	ShouldCheckpointInput struct {
		SourceCallType string
	}
	LeaseLostInput  struct{}
	ShardEndedInput struct {
		Checkpoint CheckpointFunc
	}
	ShutdownRequestedInput struct {
		SequenceNumber string
		Checkpoint     CheckpointFunc
	}
)

type RecordProcessor interface {
	Initialize(*InitializationInput)
	ProcessRecords(*ProcessRecordsInput)
	LeaseLost(*LeaseLostInput)
	ShardEnded(*ShardEndedInput)
	ShutdownRequested(*ShutdownRequestedInput)
}
