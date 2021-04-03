package main

import (
	"log"
	"os"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

func main() {
	processor := &sampleProcessor{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	process := kcl.GetKCLProcess(processor)
	err := process.Run()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("KCL Processor exited")
}
