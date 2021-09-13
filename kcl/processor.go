package kcl

type CheckpointFunc = func(*string) error

type (
	InitializationInput struct {
		ShardID string
	}
	ProcessRecordsInput struct {
		Records    []Record
		Checkpoint CheckpointFunc `json:"-"`
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
