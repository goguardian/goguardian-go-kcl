package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "redacts password from URL with basic auth credentials",
			input:    "https://aws:supersecrettoken@example.d.codeartifact.us-west-2.amazonaws.com/maven/repo/",
			expected: "https://aws:xxxxx@example.d.codeartifact.us-west-2.amazonaws.com/maven/repo/",
		},
		{
			name:     "leaves URL without credentials unchanged",
			input:    "https://repo1.maven.org/maven2/",
			expected: "https://repo1.maven.org/maven2/",
		},
		{
			name:     "returns invalid URL unchanged",
			input:    "not a url ://bad",
			expected: "not a url ://bad",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redactURL(tt.input)
			if got != tt.expected {
				t.Errorf("redactURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	// Setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("packageData"))
	}))
	defer server.Close()

	d := getDownloader(3, server.URL+"/")
	d.packages = []mavenPackage{
		{
			Artifact: "some.artifact.path",
			Group:    "some-package-group",
			Version:  "1.2.3",
		},
	}

	tempDir, err := os.MkdirTemp("", "someDir")
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
	data, err := os.ReadFile(downloadedFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "packageData" {
		t.Errorf("expected jar file to contain 'packageData', but instead it contained '%s'", string(data))
	}
}
