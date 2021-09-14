package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestDownload(t *testing.T) {
	// Setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("packageData"))
	}))
	defer server.Close()

	d := getDownloader()
	d.packages = []mavenPackage{
		{
			Artifact: "some.artifact.path",
			Group:    "some-package-group",
			Version:  "1.2.3",
		},
	}
	d.mavenBaseURL = server.URL + "/?filename="

	tempDir, err := ioutil.TempDir("", "someDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test function
	err = d.download(tempDir)
	if err != nil {
		t.Errorf("failed to download jar package: %+v", err)
	}

	// Validate file is saved properly
	downloadedFile := path.Join(tempDir, "some.artifact.path-1.2.3.jar")
	file, err := os.Open(downloadedFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "packageData" {
		t.Errorf("expected jar file to contain 'packageData', but instead it contained '%s'", string(data))
	}
}
