package kcl

type KCLProcess interface {
	Run()
}

func GetKCLProcess(r RecordProcessor) KCLProcess {
	return &kclProcess{
		recordProcessor: r,
	}
}

type kclProcess struct {
	recordProcessor RecordProcessor
}

func (k *kclProcess) Run() {
	// TODO: implement this
}
