package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

func main() {
	now := time.Now()
	file, err := os.Create(fmt.Sprintf("sample-log-%d", now.Unix()))
	if err != nil {
		log.Fatal(err)
	}

	processor := &sampleProcessor{
		logger: log.New(file, "", log.LstdFlags),
	}

	kclLogFile, err := os.Create(fmt.Sprintf("kcl-log-%d", now.Unix()))
	if err != nil {
		log.Fatal(err)
	}
	kclLogger := log.New(kclLogFile, "", log.LstdFlags)

	process := kcl.GetKCLProcess(processor, kcl.WithLogger(kclLogger))
	err = process.Run()
	if err != nil {
		processor.logger.Fatal(err)
	}

	processor.logger.Println("KCL Processor exited")
}
