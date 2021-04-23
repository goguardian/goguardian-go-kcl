package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"github.com/goguardian/goguardian-go-kcl/runner"
)

const (
	defaultTimeout = 30 * time.Second
)

type propertiesVar struct {
	StreamName string
	AppName    string
}

func TestRecordsReceived(t *testing.T) {
	testStreamName := "stream_" + uuid.Must(uuid.NewV4()).String()
	testAppName := "app_" + uuid.Must(uuid.NewV4()).String()

	t.Log("Getting local kinesis client")
	tClient, err := GetLocalKinesisClient()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("deleting kinesis stream if present")
	err = tClient.DeleteStream(testStreamName, defaultTimeout)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("creating new kinesis stream")
	err = tClient.CreateStream(testStreamName, 4, defaultTimeout)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("putting records in kinesis stream")
	err = tClient.PutRecords(testStreamName, []string{"alice", "bob", "charlie"})
	if err != nil {
		t.Fatal(err)
	}

	// Create properties file for this test consumer
	tmpl, err := template.ParseFiles("test_app_properties.tmpl")
	if err != nil {
		t.Fatal("failed to parse properties template file")
	}

	propertiesFile, err := ioutil.TempFile("", "test_app_properties")
	if err != nil {
		t.Fatal("failed to create properties file")
	}

	err = tmpl.Execute(propertiesFile, propertiesVar{
		AppName:    testAppName,
		StreamName: testStreamName,
	})
	if err != nil {
		t.Fatal("failed to populate properties file")
	}

	propertiesFile.Close()

	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		t.Fatal("JAVA_HOME environment variable not specified")
	}
	r, err := runner.GetRunner(
		runner.WithPathToJarFolder("../jar"),
		runner.WithPathToPropertiesFile(propertiesFile.Name()),
		runner.WithPathToJavaBinary(javaHome+"/bin/java"),
	)
	if err != nil {
		t.Fatal("failed to get runner")
	}

	testPassed := make(chan bool)
	receiver := GetMessageReceiver()
	go func() {
		receivedRecords := map[string]bool{}
		for {
			req := <-receiver.processRecordsChan
			fmt.Println("GOT REQUEST")
			for _, record := range req.Records {
				fmt.Println(record.Data)
				receivedRecords[string(record.Data)] = true
			}

			if receivedRecords["alice"] && receivedRecords["bob"] && receivedRecords["charlie"] {
				fmt.Println("found all the records")
				testPassed <- true
			}
		}
	}()

	go func() {
		t.Fatal(r.RunJavaDaemon())
	}()

	<-testPassed
	fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	os.Exit(0)
	// TODO: shutdown the java daemon
}
