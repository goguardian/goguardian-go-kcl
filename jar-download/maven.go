package main

import (
	"fmt"
	"strings"
)

type mavenPackage struct {
	Group    string
	Artifact string
	Version  string
}

func (pkg *mavenPackage) URL(mavenBaseURL string) string {
	paths := strings.Split(pkg.Group, ".")
	paths = append(paths, pkg.Artifact, pkg.Version, pkg.Name())
	return mavenBaseURL + strings.Join(paths, "/")
}

func (pkg *mavenPackage) Name() string {
	return fmt.Sprintf("%s-%s.jar", pkg.Artifact, pkg.Version)
}

const mavenBaseHTTPURL = "https://search.maven.org/remotecontent?filepath="

// To update the packages to newer versions follow instructions here:
// https://github.com/awslabs/amazon-kinesis-client-python/blob/master/scripts/build_deps.py
var mavenPackages = []mavenPackage{
	{"software.amazon.kinesis", "amazon-kinesis-client-multilang", "2.3.4"},
	{"software.amazon.kinesis", "amazon-kinesis-client", "2.3.4"},
	{"software.amazon.awssdk", "kinesis", "2.15.31"},
	{"software.amazon.awssdk", "aws-cbor-protocol", "2.15.31"},
	{"com.fasterxml.jackson.dataformat", "jackson-dataformat-cbor", "2.10.4"},
	{"software.amazon.awssdk", "aws-json-protocol", "2.15.31"},
	{"software.amazon.awssdk", "dynamodb", "2.15.31"},
	{"software.amazon.awssdk", "cloudwatch", "2.15.31"},
	{"software.amazon.awssdk", "netty-nio-client", "2.15.31"},
	{"io.netty", "netty-codec-http", "4.1.65.Final"},
	{"io.netty", "netty-codec-http2", "4.1.65.Final"},
	{"io.netty", "netty-codec", "4.1.65.Final"},
	{"io.netty", "netty-transport", "4.1.65.Final"},
	{"io.netty", "netty-resolver", "4.1.65.Final"},
	{"io.netty", "netty-common", "4.1.65.Final"},
	{"io.netty", "netty-buffer", "4.1.65.Final"},
	{"io.netty", "netty-handler", "4.1.65.Final"},
	{"io.netty", "netty-transport-native-epoll", "4.1.65.Final"},
	{"io.netty", "netty-transport-native-unix-common", "4.1.65.Final"},
	{"com.typesafe.netty", "netty-reactive-streams-http", "2.0.4"},
	{"com.typesafe.netty", "netty-reactive-streams", "2.0.4"},
	{"org.reactivestreams", "reactive-streams", "1.0.2"},
	{"software.amazon.glue", "schema-registry-serde", "1.0.0"},
	{"software.amazon.awssdk", "arns", "2.15.26"},
	{"org.apache.avro", "avro", "1.9.2"},
	{"org.apache.commons", "commons-compress", "1.19"},
	{"org.projectlombok", "lombok-utils", "1.18.12"},
	{"software.amazon.glue", "schema-registry-common", "1.0.0"},
	{"software.amazon.awssdk", "glue", "2.15.32"},
	{"org.apache.kafka", "kafka-clients", "2.5.0"},
	{"com.github.luben", "zstd-jni", "1.4.4-7"},
	{"org.lz4", "lz4-java", "1.7.1"},
	{"org.xerial.snappy", "snappy-java", "1.1.7.3"},
	{"org.apache.kafka", "kafka-streams", "2.5.0"},
	{"org.apache.kafka", "connect-json", "2.5.0"},
	{"org.apache.kafka", "connect-api", "2.5.0"},
	{"com.fasterxml.jackson.datatype", "jackson-datatype-jdk8", "2.10.2"},
	{"org.rocksdb", "rocksdbjni", "5.18.3"},
	{"com.google.guava", "guava", "26.0-jre"},
	{"com.google.code.findbugs", "jsr305", "3.0.2"},
	{"org.checkerframework", "checker-qual", "2.5.2"},
	{"com.google.errorprone", "error_prone_annotations", "2.1.3"},
	{"com.google.j2objc", "j2objc-annotations", "1.1"},
	{"org.codehaus.mojo", "animal-sniffer-annotations", "1.14"},
	{"com.google.protobuf", "protobuf-java", "3.11.4"},
	{"org.apache.commons", "commons-lang3", "3.8.1"},
	{"org.slf4j", "slf4j-api", "1.7.25"},
	{"io.reactivex.rxjava2", "rxjava", "2.1.14"},
	{"software.amazon.awssdk", "sts", "2.15.31"},
	{"software.amazon.awssdk", "aws-query-protocol", "2.15.31"},
	{"software.amazon.awssdk", "protocol-core", "2.15.31"},
	{"software.amazon.awssdk", "profiles", "2.15.31"},
	{"software.amazon.awssdk", "sdk-core", "2.15.31"},
	{"com.fasterxml.jackson.core", "jackson-core", "2.10.4"},
	{"com.fasterxml.jackson.core", "jackson-databind", "2.10.5.1"},
	{"software.amazon.awssdk", "auth", "2.15.31"},
	{"software.amazon.eventstream", "eventstream", "1.0.1"},
	{"software.amazon.awssdk", "http-client-spi", "2.15.31"},
	{"software.amazon.awssdk", "regions", "2.15.31"},
	{"com.fasterxml.jackson.core", "jackson-annotations", "2.10.4"},
	{"software.amazon.awssdk", "annotations", "2.15.31"},
	{"software.amazon.awssdk", "utils", "2.15.31"},
	{"software.amazon.awssdk", "aws-core", "2.15.31"},
	{"software.amazon.awssdk", "metrics-spi", "2.15.31"},
	{"software.amazon.awssdk", "apache-client", "2.15.31"},
	{"org.apache.httpcomponents", "httpclient", "4.5.13"},
	{"commons-codec", "commons-codec", "1.11"},
	{"org.apache.httpcomponents", "httpcore", "4.4.11"},
	{"com.amazonaws", "aws-java-sdk-core", "1.11.477"},
	{"commons-logging", "commons-logging", "1.1.3"},
	{"software.amazon.ion", "ion-java", "1.0.2"},
	{"joda-time", "joda-time", "2.8.1"},
	{"ch.qos.logback", "logback-classic", "1.2.3"},
	{"ch.qos.logback", "logback-core", "1.2.3"},
	{"com.beust", "jcommander", "1.72"},
	{"commons-io", "commons-io", "2.6"},
	{"org.apache.commons", "commons-collections4", "4.2"},
	{"commons-beanutils", "commons-beanutils", "1.9.4"},
	{"commons-collections", "commons-collections", "3.2.2"},
}
