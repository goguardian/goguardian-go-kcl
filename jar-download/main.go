package main

import (
	"fmt"
	"os"
)

// Credit to https://github.com/arthurbailao/aws-kcl/blob/master/cmd/aws-kcl/download.go
// since this code was adapted from that. It was adapted to work with the KCL 2.0 jar
// dependencies which was taken from here
// https://github.com/awslabs/amazon-kinesis-client-python/blob/master/setup.py#L49
func main() {
	fmt.Println("Starting download")
	dstFolder := "./jar"
	if len(os.Args) > 1 {
		dstFolder = os.Args[1]
	}

	d := getDownloader()
	err := d.download(dstFolder)
	if err != nil {
		fmt.Printf("failed to download due to error: %+v\n", err)
	}

	fmt.Println("Completed download")
}
