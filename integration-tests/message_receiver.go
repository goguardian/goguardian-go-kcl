package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/goguardian/goguardian-go-kcl/kcl"
)

const ReceiverPort = ":8080"

type messageReceiver struct {
	processRecordsChan chan *kcl.ProcessRecordsInput
}

func GetMessageReceiver() *messageReceiver {
	m := &messageReceiver{
		processRecordsChan: make(chan *kcl.ProcessRecordsInput),
	}

	go m.startHTTPServer()
	return m
}

func (m *messageReceiver) startHTTPServer() {
	http.HandleFunc("/process_records", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		var processRecords kcl.ProcessRecordsInput
		err = json.Unmarshal(body, &processRecords)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		m.processRecordsChan <- &processRecords

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("message received"))
	})

	log.Fatal(http.ListenAndServe(ReceiverPort, nil))
}
