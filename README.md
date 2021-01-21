# GoGuardian-Go-KCL
Kinesis Client Library for Go

## Overview

## How to run

### Prerequisites

- Install [Go](https://golang.org/)

Make sure your go version matches the go version in the `go.mod` file.

## Documentation

### Using Amazon KCL's MultiLangDaemon
Amazon has implemented a [MultiLangDaemon](https://github.com/awslabs/amazon-kinesis-client/tree/master/amazon-kinesis-client-multilang) which is a java-daemon that implements most of the important kinesis based functionalities under the hood. Our approach was to run this daemon in the parent thread which itself spawns child threads for the consumer/record processor and communicates with it via STDIN/STDOUT.
