package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
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

	err := download(dstFolder)
	if err != nil {
		fmt.Printf("failed to download due to error: %+v\n", err)
	}

	fmt.Println("Completed download")
}

// Download fetches each jar package from Maven and saves it in the specified
// dstPath.
func download(dstPath string) error {
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		err := os.Mkdir(dstPath, 0755)
		if err != nil {
			return errors.Wrap(err, "failed to make jar directory")
		}
	}

	for _, pkg := range packages {
		filename := path.Join(dstPath, pkg.Name())

		if err := downloadFileWithRetry(pkg.URL(), filename); err != nil {
			return err
		}
	}
	return nil
}

func downloadFileWithRetry(src string, dst string) error {
	// don't download if dst already xists
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	retryTimes := 3
	backoff := 500 * time.Millisecond
	var err error

	for i := 0; i < retryTimes; i++ {
		fmt.Printf("Downloading %s to %s\n", src, dst)
		err = downloadFile(src, dst)
		if err == nil {
			break
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	if err != nil {
		return err
	}

	return nil
}

func downloadFile(src string, dst string) error {
	resp, err := http.Get(src)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create destination file: %s", dst)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}

	return nil
}

type mavenPackageInfo struct {
	Group    string
	Artifact string
	Version  string
}

func (pkg *mavenPackageInfo) URL() string {
	paths := strings.Split(pkg.Group, ".")
	paths = append(paths, pkg.Artifact, pkg.Version, pkg.Name())
	return "http://search.maven.org/remotecontent?filepath=" + strings.Join(paths, "/")
}

func (pkg *mavenPackageInfo) Name() string {
	return fmt.Sprintf("%s-%s.jar", pkg.Artifact, pkg.Version)
}

// To update the packages to newer versions follow instructions here:
// https://github.com/awslabs/amazon-kinesis-client-python/blob/master/scripts/build_deps.py
var packages = [...]mavenPackageInfo{
	{"software.amazon.kinesis", "amazon-kinesis-client-multilang", "2.3.1"},
	{"software.amazon.kinesis", "amazon-kinesis-client", "2.3.1"},
	{"software.amazon.awssdk", "kinesis", "2.14.0"},
	{"software.amazon.awssdk", "aws-cbor-protocol", "2.14.0"},
	{"com.fasterxml.jackson.dataformat", "jackson-dataformat-cbor", "2.10.4"},
	{"software.amazon.awssdk", "aws-json-protocol", "2.14.0"},
	{"software.amazon.awssdk", "dynamodb", "2.14.0"},
	{"software.amazon.awssdk", "cloudwatch", "2.14.0"},
	{"software.amazon.awssdk", "netty-nio-client", "2.14.0"},
	{"io.netty", "netty-codec-http", "4.1.46.Final"},
	{"io.netty", "netty-codec-http2", "4.1.46.Final"},
	{"io.netty", "netty-codec", "4.1.46.Final"},
	{"io.netty", "netty-transport", "4.1.46.Final"},
	{"io.netty", "netty-resolver", "4.1.46.Final"},
	{"io.netty", "netty-common", "4.1.46.Final"},
	{"io.netty", "netty-buffer", "4.1.46.Final"},
	{"io.netty", "netty-handler", "4.1.46.Final"},
	{"io.netty", "netty-transport-native-epoll", "4.1.46.Final"},
	{"io.netty", "netty-transport-native-unix-common", "4.1.46.Final"},
	{"com.typesafe.netty", "netty-reactive-streams-http", "2.0.4"},
	{"com.typesafe.netty", "netty-reactive-streams", "2.0.4"},
	{"org.reactivestreams", "reactive-streams", "1.0.2"},
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
	{"software.amazon.awssdk", "sts", "2.14.0"},
	{"software.amazon.awssdk", "aws-query-protocol", "2.14.0"},
	{"software.amazon.awssdk", "protocol-core", "2.14.0"},
	{"software.amazon.awssdk", "profiles", "2.14.0"},
	{"software.amazon.awssdk", "sdk-core", "2.14.0"},
	{"com.fasterxml.jackson.core", "jackson-core", "2.10.4"},
	{"com.fasterxml.jackson.core", "jackson-databind", "2.10.4"},
	{"software.amazon.awssdk", "auth", "2.14.0"},
	{"software.amazon.eventstream", "eventstream", "1.0.1"},
	{"software.amazon.awssdk", "http-client-spi", "2.14.0"},
	{"software.amazon.awssdk", "regions", "2.14.0"},
	{"com.fasterxml.jackson.core", "jackson-annotations", "2.10.4"},
	{"software.amazon.awssdk", "annotations", "2.14.0"},
	{"software.amazon.awssdk", "utils", "2.14.0"},
	{"software.amazon.awssdk", "aws-core", "2.14.0"},
	{"software.amazon.awssdk", "metrics-spi", "2.14.0"},
	{"software.amazon.awssdk", "apache-client", "2.14.0"},
	{"org.apache.httpcomponents", "httpclient", "4.5.9"},
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
