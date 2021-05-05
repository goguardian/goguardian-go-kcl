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

### Running integration tests
* Ensure you have docker-compose https://docs.docker.com/compose/install/. 
* Also ensure you have the $JAVA_HOME environment variable set. For MacOS see the guide here https://mkyong.com/java/how-to-set-java_home-environment-variable-on-mac-os-x/.

```bash
make run_integ_test
```

### TODOs
- [ ] Add unit tests to the kcl package
- [ ] Add a license and contributing.md
- [ ] Add a docs folder with more information
- [ ] Update README
  - [ ] Information on how to run the sample
  - [ ] give attribution to the Python KCL
