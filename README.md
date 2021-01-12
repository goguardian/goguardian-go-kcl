# GoGuardian-Go-KCL
Kinesis Client Library for Go

## Overview

## How to run

## Documentation

### Using Amazon KCL's MultiLangDaemon
Amazon has implemented a MultiLangDaemon which is a java-daemon that implements most of the important kinesis based functionalities under the hood. Our approach was to run this daemon in the parent thread which itself spawns child threads for the consumer/record processor and communicates with it via STDIN/STDOUT.
