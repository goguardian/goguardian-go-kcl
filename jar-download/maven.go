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
// TODO: Be careful when updating the aws java sdk dependencies:
// Do not use AWS SDK for Java version 2.27.19 to 2.27.23 with KCL 3.x.
// These versions include an issue that causes an exception error related to KCL's DynamoDB usage.
// We recommend that you use the AWS SDK for Java version 2.28.0 or later to avoid this issue.
// See: https://docs.aws.amazon.com/streams/latest/dev/kcl-migration-from-2-3.html#kcl-migration-from-2-3-prerequisites
var mavenPackages = []mavenPackage{
	{"software.amazon.kinesis", "amazon-kinesis-client-multilang", "3.1.1"},
	{"software.amazon.kinesis", "amazon-kinesis-client", "3.1.1"},
	{"software.amazon.awssdk", "kinesis", "2.31.62"},
	{"software.amazon.awssdk", "aws-cbor-protocol", "2.31.62"},
	{"software.amazon.awssdk", "third-party-jackson-dataformat-cbor", "2.31.62"},
	{"software.amazon.awssdk", "aws-json-protocol", "2.31.62"},
	{"software.amazon.awssdk", "dynamodb", "2.31.62"},
	{"software.amazon.awssdk", "dynamodb-enhanced", "2.31.62"},
	{"com.amazonaws", "dynamodb-lock-client", "1.3.0"},
	{"software.amazon.awssdk", "cloudwatch", "2.31.62"},
	{"software.amazon.awssdk", "netty-nio-client", "2.31.62"},
	{"io.netty", "netty-codec-http", "4.1.118.Final"},
	{"io.netty", "netty-codec-http2", "4.1.118.Final"},
	{"io.netty", "netty-codec", "4.1.118.Final"},
	{"io.netty", "netty-transport", "4.1.118.Final"},
	{"io.netty", "netty-common", "4.1.118.Final"},
	{"io.netty", "netty-buffer", "4.1.118.Final"},
	{"io.netty", "netty-transport-classes-epoll", "4.1.118.Final"},
	{"io.netty", "netty-resolver", "4.1.118.Final"},
	{"software.amazon.awssdk", "sdk-core", "2.31.62"},
	{"software.amazon.awssdk", "checksums-spi", "2.31.62"},
	{"software.amazon.awssdk", "checksums", "2.31.62"},
	{"software.amazon.awssdk", "retries", "2.31.62"},
	{"software.amazon.awssdk", "aws-core", "2.31.62"},
	{"software.amazon.eventstream", "eventstream", "1.0.1"},
	{"software.amazon.awssdk", "arns", "2.31.62"},
	{"software.amazon.awssdk", "regions", "2.31.62"},
	{"software.amazon.awssdk", "utils", "2.31.62"},
	{"software.amazon.awssdk", "http-client-spi", "2.31.62"},
	{"software.amazon.glue", "schema-registry-serde", "1.1.24"},
	{"org.apache.kafka", "kafka-clients", "3.6.1"},
	{"com.github.luben", "zstd-jni", "1.5.5-1"},
	{"org.lz4", "lz4-java", "1.8.0"},
	{"org.xerial.snappy", "snappy-java", "1.1.10.5"},
	{"com.kjetland", "mbknor-jackson-jsonschema_2.12", "1.0.39"},
	{"org.scala-lang", "scala-library", "2.12.10"},
	{"javax.validation", "validation-api", "2.0.1.Final"},
	{"io.github.classgraph", "classgraph", "4.8.120"},
	{"com.github.erosb", "everit-json-schema", "1.14.5"},
	{"org.json", "json", "20250107"},
	{"commons-validator", "commons-validator", "1.9.0"},
	{"commons-digester", "commons-digester", "2.1"},
	{"com.damnhandy", "handy-uri-templates", "2.1.8"},
	{"joda-time", "joda-time", "2.10.2"},
	{"com.google.re2j", "re2j", "1.8"},
	{"org.jetbrains.kotlin", "kotlin-stdlib", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-stdlib-jdk8", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-stdlib-jdk7", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-reflect", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-scripting-compiler-impl-embeddable", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-scripting-common", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-scripting-jvm", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-script-runtime", "1.9.25"},
	{"org.jetbrains.kotlin", "kotlin-scripting-compiler-embeddable", "1.9.25"},
	{"com.squareup.okio", "okio", "3.4.0"},
	{"com.squareup.okio", "okio-jvm", "3.4.0"},
	{"com.squareup.okio", "okio-fakefilesystem", "3.2.0"},
	{"com.squareup.okio", "okio-fakefilesystem-jvm", "3.2.0"},
	{"org.jetbrains.kotlinx", "kotlinx-datetime-jvm", "0.3.2"},
	{"org.jetbrains.kotlinx", "kotlinx-serialization-core-jvm", "1.4.0"},
	{"org.jetbrains.kotlin", "kotlin-stdlib-common", "1.7.10"},
	{"com.squareup.wire", "wire-schema", "5.2.0"},
	{"com.squareup.wire", "wire-runtime", "5.2.0"},
	{"com.squareup.wire", "wire-compiler", "5.2.0"},
	{"com.squareup.wire", "wire-schema-jvm", "5.2.0"},
	{"com.palantir.javapoet", "javapoet", "0.6.0"},
	{"com.squareup", "kotlinpoet-jvm", "2.0.0"},
	{"com.squareup.wire", "wire-runtime-jvm", "5.2.0"},
	{"com.squareup.wire", "wire-kotlin-generator", "5.2.0"},
	{"com.squareup.wire", "wire-grpc-client-jvm", "5.2.0"},
	{"com.squareup.okhttp3", "okhttp", "5.0.0-alpha.14"},
	{"org.jetbrains.kotlinx", "kotlinx-coroutines-core-jvm", "1.9.0"},
	{"com.squareup.wire", "wire-java-generator", "5.2.0"},
	{"com.squareup.wire", "wire-swift-generator", "5.2.0"},
	{"io.outfoxx", "swiftpoet", "1.6.5"},
	{"com.charleskorn.kaml", "kaml-jvm", "0.67.0"},
	{"it.krzeminski", "snakeyaml-engine-kmp-jvm", "3.0.3"},
	{"net.thauvin.erik.urlencoder", "urlencoder-lib-jvm", "1.6.0"},
	{"com.google.api.grpc", "proto-google-common-protos", "2.7.4"},
	{"com.google.jimfs", "jimfs", "1.1"},
	{"software.amazon.glue", "schema-registry-common", "1.1.24"},
	{"software.amazon.awssdk", "glue", "2.22.12"},
	{"software.amazon.glue", "schema-registry-build-tools", "1.1.24"},
	{"software.amazon.awssdk", "url-connection-client", "2.22.12"},
	{"org.apache.avro", "avro", "1.11.4"},
	{"com.google.guava", "guava", "32.1.1-jre"},
	{"com.google.guava", "failureaccess", "1.0.1"},
	{"com.google.guava", "listenablefuture", "9999.0-empty-to-avoid-conflict-with-guava"},
	{"org.checkerframework", "checker-qual", "3.33.0"},
	{"com.google.errorprone", "error_prone_annotations", "2.18.0"},
	{"com.google.j2objc", "j2objc-annotations", "2.8"},
	{"com.google.protobuf", "protobuf-java", "4.27.5"},
	{"org.apache.commons", "commons-lang3", "3.14.0"},
	{"commons-collections", "commons-collections", "3.2.2"},
	{"io.netty", "netty-handler", "4.1.118.Final"},
	{"io.netty", "netty-transport-native-unix-common", "4.1.118.Final"},
	{"com.google.code.findbugs", "jsr305", "3.0.2"},
	{"com.fasterxml.jackson.core", "jackson-databind", "2.12.7.1"},
	{"com.fasterxml.jackson.core", "jackson-annotations", "2.12.7"},
	{"com.fasterxml.jackson.core", "jackson-core", "2.12.7"},
	{"org.reactivestreams", "reactive-streams", "1.0.4"},
	{"software.amazon.awssdk", "annotations", "2.31.62"},
	{"org.slf4j", "slf4j-api", "2.0.13"},
	{"org.jetbrains", "annotations", "26.0.1"},
	{"io.reactivex.rxjava3", "rxjava", "3.1.11"},
	{"software.amazon.awssdk", "sts", "2.31.62"},
	{"software.amazon.awssdk", "aws-query-protocol", "2.31.62"},
	{"software.amazon.awssdk", "protocol-core", "2.31.62"},
	{"software.amazon.awssdk", "profiles", "2.31.62"},
	{"software.amazon.awssdk", "http-auth-aws", "2.31.62"},
	{"software.amazon.awssdk", "auth", "2.31.62"},
	{"software.amazon.awssdk", "http-auth-aws-eventstream", "2.31.62"},
	{"software.amazon.awssdk", "http-auth-spi", "2.31.62"},
	{"software.amazon.awssdk", "http-auth", "2.31.62"},
	{"software.amazon.awssdk", "identity-spi", "2.31.62"},
	{"software.amazon.awssdk", "metrics-spi", "2.31.62"},
	{"software.amazon.awssdk", "json-utils", "2.31.62"},
	{"software.amazon.awssdk", "third-party-jackson-core", "2.31.62"},
	{"software.amazon.awssdk", "endpoints-spi", "2.31.62"},
	{"software.amazon.awssdk", "retries-spi", "2.31.62"},
	{"software.amazon.awssdk", "apache-client", "2.31.62"},
	{"org.apache.httpcomponents", "httpclient", "4.5.13"},
	{"org.apache.httpcomponents", "httpcore", "4.4.16"},
	{"commons-codec", "commons-codec", "1.17.1"},
	{"ch.qos.logback", "logback-classic", "1.3.14"},
	{"ch.qos.logback", "logback-core", "1.3.14"},
	{"com.beust", "jcommander", "1.82"},
	{"commons-io", "commons-io", "2.16.1"},
	{"org.apache.commons", "commons-collections4", "4.4"},
	{"commons-beanutils", "commons-beanutils", "1.11.0"},
	{"commons-logging", "commons-logging", "1.3.5"},
}
