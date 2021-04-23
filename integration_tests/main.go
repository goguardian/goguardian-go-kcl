package main

import (
	"log"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

func main() {
	processor := &testProcessor{}

	process := kcl.GetKCLProcess(processor)
	err := process.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("KCL Processor exited")
}
