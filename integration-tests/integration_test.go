package integration_tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"text/template"
	"time"

	"github.com/goguardian/goguardian-go-kcl/runner"
)

const defaultTimeout = 30 * time.Second

type propertiesVar struct {
	StreamName string
	AppName    string
}

func TestRecordsReceived(t *testing.T) {
	// comment this line to run integration tests
	// TODO: make this cleaner
	t.Skip("Skipping testing in CI environment")

	now := time.Now()
	testStreamName := fmt.Sprintf("stream_%d", now.Unix())
	testAppName := fmt.Sprintf("app_%d", now.Unix())

	fmt.Println("Getting local kinesis client")
	tClient, err := GetLocalKinesisClient()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("deleting kinesis stream if present")
	err = tClient.DeleteStream(testStreamName, defaultTimeout)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("creating new kinesis stream")
	err = tClient.CreateStream(testStreamName, 4, defaultTimeout)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("putting records in kinesis stream")
	err = tClient.PutRecords(testStreamName, []string{"alice", "bob", "charlie"})
	if err != nil {
		t.Fatal(err)
	}

	// Create properties file for this test consumer
	tmpl, err := template.ParseFiles("test-app/test_app_properties.tmpl")
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
			for _, record := range req.Records {
				receivedRecords[string(record.Data)] = true
			}

			if receivedRecords["alice"] && receivedRecords["bob"] && receivedRecords["charlie"] {
				fmt.Println("found all the records")
				testPassed <- true
			}
		}
	}()

	var cmd *exec.Cmd
	go func() {
		var err error
		cmd, err = r.RunJavaDaemon("-Daws.accessKeyId=some_key", "-Daws.secretKey=some_secret_key")
		if err != nil {
			t.Fatal(err)
		}
	}()

	<-testPassed
	cmd.Process.Kill()
}
