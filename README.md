# Amazon Kinesis Client Library for Go

This package provides an interface to the Amazon Kinesis Client Library (KCL) MultiLangDaemon,
which is part of the [Amazon KCL for Java][kinesis-github].
Developers can use the [Amazon KCL][amazon-kcl] to build distributed applications that
process streaming data reliably at scale. The [Amazon KCL][amazon-kcl] takes care of
many of the complex tasks associated with distributed computing, such as load-balancing
across multiple instances, responding to instance failures, checkpointing processed records,
and reacting to changes in stream volume.
This interface manages the interaction with the MultiLangDaemon so that developers can focus on
implementing their record processor executable. A record processor executable
typically looks something like:

```go
package main

import (
	"github.com/goguardian/goguardian-go-kcl/kcl"
)

type myProcessor struct{}

func (m *myProcessor) Initialize(*kcl.InitializationInput)           { /* handle init */ }
func (m *myProcessor) ProcessRecords(*kcl.ProcessRecordsInput)       { /* handle process */ }
func (m *myProcessor) LeaseLost(*kcl.LeaseLostInput)                 { /* handle lease lost */ }
func (m *myProcessor) ShardEnded(*kcl.ShardEndedInput)               { /* handle shard end */ }
func (m *myProcessor) ShutdownRequested(*kcl.ShutdownRequestedInput) { /* handle shutdown */ }

func main() {
	processor := &myProcessor{}
	process := kcl.GetKCLProcess(processor)
	err := process.Run()
	if err != nil {
		panic(err)
	}
}
```

## Before You Get Started

Install [Go][go-install] and make sure your go version matches the go version
in the `go.mod` file.

Before running the sample, you'll want to make sure that your environment is
configured to allow the sample to use your [AWS Security
Credentials](http://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html).

By default, the sample uses the [DefaultCredentialsProvider][DefaultCredentialsProvider]
so you'll want to make your credentials available to one of the credentials providers in that
provider chain. There are several ways to do this such as providing a ~/.aws/credentials file,
or if you're running on EC2, you can associate an IAM role with your instance with appropriate
access.

## Running the Sample

There is a single Makefile target to run the sample. Simply run:

```bash
make run_sample
```

This command will do the following:
1) Download jars necessary to run the MultiLangDaemon into a `./jar` folder.
2) Build the runner app located at [./runner](runner) which is used to run the Java MultiLangDaemon.
3) Build the sample processor executable located at [./sample](sample).
4) Run the Java MultiLangDaemon which will spawn the sample processor.

### Running integration tests
Ensure you have [docker-compose][docker-compose-install]. We
leverage [LocalStack][localstack] to emulate Kinesis locally. Also ensure you have the
`$JAVA_HOME` environment variable set. For MacOS see this [guide](https://mkyong.com/java/how-to-set-java_home-environment-variable-on-mac-os-x/).

```bash
make run_integ_test
```

## Under the Hood - What You Should Know about Amazon KCL's [MultiLangDaemon][multi-lang-daemon]
Amazon KCL for Go uses [Amazon KCL for Java][kinesis-github] internally. AWS
implemented a Java-based daemon, called the *MultiLangDaemon* that does all the
heavy lifting. This approach has the daemon spawn the user-defined record
processor script/program as a sub-process. The *MultiLangDaemon* communicates
with this sub-process over standard input/output using a simple protocol, and
therefore the record processor script/program can be written in any language.

At runtime, there will always be a one-to-one correspondence between a record processor, a child process,
and an [Amazon Kinesis Shard][amazon-kinesis-shard]. The *MultiLangDaemon* will make sure of
that, without any need for the developer to intervene.

In this release, we have abstracted these implementation details away and exposed an interface that enables
you to focus on writing record processing logic in Go. This approach enables [Amazon KCL][amazon-kcl] to
be language agnostic, while providing identical features and similar parallel processing model across
all languages.

## See Also
* [Developing Consumer Applications for Amazon Kinesis Using the Amazon Kinesis Client Library][amazon-kcl]
* The [Amazon KCL for Java][kinesis-github]
* The [Amazon KCL for Ruby][amazon-kinesis-ruby-github]
* The [Amazon Kinesis Documentation][amazon-kinesis-docs]
* The [Amazon Kinesis Forum][kinesis-forum]

## Release Notes
TBD

[amazon-kinesis-shard]: http://docs.aws.amazon.com/kinesis/latest/dev/key-concepts.html
[amazon-kinesis-docs]: http://aws.amazon.com/documentation/kinesis/
[amazon-kcl]: http://docs.aws.amazon.com/kinesis/latest/dev/kinesis-record-processor-app.html
[multi-lang-daemon]: https://github.com/awslabs/amazon-kinesis-client/blob/master/amazon-kinesis-client-multilang/src/main/java/software/amazon/kinesis/multilang/package-info.java
[kinesis]: http://aws.amazon.com/kinesis
[amazon-kinesis-ruby-github]: https://github.com/awslabs/amazon-kinesis-client-ruby
[kinesis-github]: https://github.com/awslabs/amazon-kinesis-client
[boto]: http://boto.readthedocs.org/en/latest/
[DefaultCredentialsProvider]: https://sdk.amazonaws.com/java/api/latest/software/amazon/awssdk/auth/credentials/DefaultCredentialsProvider.html
[kinesis-forum]: http://developer.amazonwebservices.com/connect/forum.jspa?forumID=169
[go-install]: https://golang.org/doc/install
[docker-compose-install]: https://docs.docker.com/compose/install/
[localstack]: https://github.com/localstack/localstack
[license]: https://github.com/goguardian/goguardian-go-kcl/blob/main/LICENSE

## License
[License][license]
