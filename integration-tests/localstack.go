package integration_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func waitForLocalstack(timeout time.Duration) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return errors.New("timed out waiting for localstack to start")
		default:
		}

		resp, err := http.Get(fmt.Sprintf("%s/health", localstackEndpoint))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read localstack health")
		}

		health := struct {
			Services map[string]string `json:"services"`
		}{}
		err = json.Unmarshal(body, &health)
		if err != nil {
			return errors.Wrap(err, "failed to parse localstack health")
		}

		if health.Services["kinesis"] == "running" && health.Services["dynamodb"] == "running" {
			return nil
		}
	}
}
