package main

import (
	"fmt"
	"os"
	"strconv"
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

	mavenBaseURL := "https://repo1.maven.org/maven2/"
	if os.Getenv("MAVEN_BASE_URL") != "" {
		mavenBaseURL = os.Getenv("MAVEN_BASE_URL")
	}
	fmt.Printf("Using Maven base URL: %s\n", mavenBaseURL)

	maxRetries := 3
	if os.Getenv("MAX_MAVEN_HTTP_RETRIES") != "" {
		var err error
		maxRetries, err = strconv.Atoi(os.Getenv("MAX_MAVEN_HTTP_RETRIES"))
		if err != nil {
			fmt.Printf("failed to parse MAX_MAVEN_HTTP_RETRIES: %+v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Using max Maven HTTP retries: %d\n", maxRetries)

	d := getDownloader(maxRetries, mavenBaseURL)
	err := d.download(dstFolder)
	if err != nil {
		fmt.Printf("failed to download due to error: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("Completed download")
}
