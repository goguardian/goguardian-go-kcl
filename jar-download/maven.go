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
	{"software.amazon.kinesis", "amazon-kinesis-client-multilang", "2.5.8"},
	{"software.amazon.kinesis", "amazon-kinesis-client", "2.5.8"},
	{"software.amazon.glue", "schema-registry-common", "1.1.19"},
	{"software.amazon.awssdk", "glue", "2.22.12"},
	{"software.amazon.glue", "schema-registry-build-tools", "1.1.19"},
	{"software.amazon.awssdk", "url-connection-client", "2.22.12"},
	{"org.apache.avro", "avro", "1.11.3"},
	{"org.apache.commons", "commons-compress", "1.21"},
	{"software.amazon.awssdk", "kinesis", "2.25.11"},
	{"software.amazon.awssdk", "dynamodb", "2.25.11"},
	{"software.amazon.awssdk", "cloudwatch", "2.25.11"},
	{"software.amazon.awssdk", "netty-nio-client", "2.25.11"},
	{"io.netty", "netty-transport-classes-epoll", "4.1.107.Final"},
	{"software.amazon.awssdk", "metrics-spi", "2.25.11"},
	{"software.amazon.awssdk", "sts", "2.25.11"},
	{"software.amazon.awssdk", "protocol-core", "2.25.11"},
	{"software.amazon.awssdk", "aws-query-protocol", "2.25.11"},
	{"software.amazon.awssdk", "aws-cbor-protocol", "2.25.11"},
	{"software.amazon.awssdk", "aws-json-protocol", "2.25.11"},
	{"software.amazon.awssdk", "json-utils", "2.25.11"},
	{"software.amazon.awssdk", "third-party-jackson-core", "2.25.11"},
	{"software.amazon.awssdk", "third-party-jackson-dataformat-cbor", "2.25.11"},
	{"software.amazon.awssdk", "profiles", "2.25.11"},
	{"software.amazon.awssdk", "sdk-core", "2.25.11"},
	{"software.amazon.awssdk", "aws-core", "2.25.11"},
	{"software.amazon.eventstream", "eventstream", "1.0.1"},
	{"software.amazon.awssdk", "endpoints-spi", "2.25.11"},
	{"software.amazon.awssdk", "auth", "2.25.11"},
	{"software.amazon.awssdk", "http-client-spi", "2.25.11"},
	{"software.amazon.awssdk", "regions", "2.25.11"},
	{"software.amazon.awssdk", "annotations", "2.25.11"},
	{"software.amazon.awssdk", "utils", "2.25.11"},
	{"software.amazon.awssdk", "apache-client", "2.25.11"},
	{"software.amazon.awssdk", "arns", "2.25.11"},
	{"software.amazon.awssdk", "http-auth-spi", "2.25.11"},
	{"software.amazon.awssdk", "http-auth", "2.25.11"},
	{"software.amazon.awssdk", "http-auth-aws", "2.25.11"},
	{"software.amazon.awssdk", "checksums-spi", "2.25.11"},
	{"software.amazon.awssdk", "checksums", "2.25.11"},
	{"software.amazon.awssdk", "identity-spi", "2.25.11"},
	{"io.netty", "netty-codec-http", "4.1.108.Final"},
	{"io.netty", "netty-codec-http2", "4.1.108.Final"},
	{"io.netty", "netty-codec", "4.1.108.Final"},
	{"io.netty", "netty-transport", "4.1.108.Final"},
	{"io.netty", "netty-resolver", "4.1.108.Final"},
	{"io.netty", "netty-common", "4.1.108.Final"},
	{"io.netty", "netty-buffer", "4.1.108.Final"},
	{"io.netty", "netty-handler", "4.1.108.Final"},
	{"io.netty", "netty-transport-native-epoll", "4.1.108.Final"},
	{"io.netty", "netty-transport-native-unix-common", "4.1.108.Final"},
	{"com.typesafe.netty", "netty-reactive-streams-http", "2.0.6"},
	{"com.typesafe.netty", "netty-reactive-streams", "2.0.6"},
	{"org.reactivestreams", "reactive-streams", "1.0.3"},
	{"com.google.guava", "guava", "32.1.1-jre"},
	{"com.google.guava", "failureaccess", "1.0.1"},
	{"com.google.guava", "listenablefuture", "9999.0-empty-to-avoid-conflict-with-guava"},
	{"com.google.code.findbugs", "jsr305", "3.0.2"},
	{"org.checkerframework", "checker-qual", "2.5.2"},
	{"com.google.errorprone", "error_prone_annotations", "2.7.1"},
	{"com.google.j2objc", "j2objc-annotations", "1.3"},
	{"org.codehaus.mojo", "animal-sniffer-annotations", "1.20"},
	{"com.google.protobuf", "protobuf-java", "3.21.7"},
	{"org.apache.commons", "commons-lang3", "3.12.0"},
	{"org.slf4j", "slf4j-api", "2.0.5"},
	{"io.reactivex.rxjava3", "rxjava", "3.1.5"},
	{"com.fasterxml.jackson.dataformat", "jackson-dataformat-cbor", "2.13.5"},
	{"com.fasterxml.jackson.core", "jackson-core", "2.13.5"},
	{"com.fasterxml.jackson.core", "jackson-databind", "2.13.5"},
	{"com.fasterxml.jackson.core", "jackson-annotations", "2.13.5"},
	{"software.amazon", "flow", "1.7"},
	{"org.apache.httpcomponents", "httpclient", "4.5.13"},
	{"commons-codec", "commons-codec", "1.15"},
	{"org.apache.httpcomponents", "httpcore", "4.4.15"},
	{"com.amazonaws", "aws-java-sdk-core", "1.12.668"},
	{"com.amazonaws", "aws-java-sdk-sts", "1.12.668"},
	{"com.amazonaws", "jmespath-java", "1.12.668"},
	{"com.amazon.ion", "ion-java", "1.11.4"},
	{"software.amazon.glue", "schema-registry-serde", "1.1.13"},
	{"org.apache.kafka", "kafka-clients", "2.8.1"},
	{"com.github.luben", "zstd-jni", "1.4.9-1"},
	{"org.lz4", "lz4-java", "1.7.1"},
	{"org.xerial.snappy", "snappy-java", "1.1.8.1"},
	{"com.kjetland", "mbknor-jackson-jsonschema_2.12", "1.0.39"},
	{"org.scala-lang", "scala-library", "2.12.10"},
	{"javax.validation", "validation-api", "2.0.1.Final"},
	{"io.github.classgraph", "classgraph", "4.8.120"},
	{"com.github.erosb", "everit-json-schema", "1.12.2"},
	{"org.json", "json", "20201115"},
	{"commons-validator", "commons-validator", "1.6"},
	{"commons-digester", "commons-digester", "1.8.1"},
	{"com.damnhandy", "handy-uri-templates", "2.1.8"},
	{"com.google.re2j", "re2j", "1.3"},
	{"org.jetbrains.kotlin", "kotlin-stdlib", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-stdlib-common", "1.7.10"},
	{"org.jetbrains", "annotations", "13.0"},
	{"org.jetbrains.kotlin", "kotlin-stdlib-jdk8", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-stdlib-jdk7", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-reflect", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-scripting-compiler-impl-embeddable", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-scripting-common", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-scripting-jvm", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-script-runtime", "1.7.10"},
	{"org.jetbrains.kotlin", "kotlin-scripting-compiler-embeddable", "1.7.10"},
	{"org.jetbrains.kotlinx", "kotlinx-serialization-core-jvm", "1.4.0"},
	{"com.squareup.wire", "wire-schema", "3.7.1"},
	{"com.squareup.wire", "wire-runtime", "3.7.1"},
	{"com.squareup.okio", "okio", "2.8.0"},
	{"com.squareup.wire", "wire-compiler", "3.7.1"},
	{"com.squareup.wire", "wire-kotlin-generator", "3.7.1"},
	{"com.squareup", "kotlinpoet", "1.7.2"},
	{"com.squareup.wire", "wire-grpc-server-generator", "3.7.1"},
	{"com.squareup.wire", "wire-java-generator", "3.7.1"},
	{"com.squareup", "javapoet", "1.13.0"},
	{"com.squareup.wire", "wire-swift-generator", "3.7.1"},
	{"io.outfoxx", "swiftpoet", "1.0.0"},
	{"com.squareup.wire", "wire-profiles", "3.7.1"},
	{"com.charleskorn.kaml", "kaml", "0.20.0"},
	{"org.snakeyaml", "snakeyaml-engine", "2.1"},
	{"com.google.api.grpc", "proto-google-common-protos", "2.7.4"},
	{"com.google.jimfs", "jimfs", "1.1"},
	{"joda-time", "joda-time", "2.10.13"},
	{"ch.qos.logback", "logback-classic", "1.3.12"},
	{"ch.qos.logback", "logback-core", "1.3.12"},
	{"com.beust", "jcommander", "1.82"},
	{"commons-io", "commons-io", "2.11.0"},
	{"commons-logging", "commons-logging", "1.1.3"},
	{"org.apache.commons", "commons-collections4", "4.2"},
	{"commons-beanutils", "commons-beanutils", "1.9.4"},
	{"commons-collections", "commons-collections", "3.2.2"},
}
