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
	LeaseLostInput  struct{}
	ShardEndedInput struct {
		Checkpoint CheckpointFunc `json:"-"`
	}
	ShutdownRequestedInput struct {
		Checkpoint CheckpointFunc `json:"-"`
	}
)

type RecordProcessor interface {
	Initialize(*InitializationInput)
	ProcessRecords(*ProcessRecordsInput)
	LeaseLost(*LeaseLostInput)
	ShardEnded(*ShardEndedInput)
	ShutdownRequested(*ShutdownRequestedInput)
}
