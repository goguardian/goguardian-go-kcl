package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type runner struct {
	logger *log.Logger

	pathToJavaBinary     string
	pathToPropertiesFile string
	pathToJarFolder      string
}

func getJarPaths(jarFolder string) ([]string, error) {
	files, err := ioutil.ReadDir(jarFolder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read jar folder")
	}

	if len(files) == 0 {
		return nil, errors.New("empty jar folder")
	}

	jarPaths := []string{}
	for _, file := range files {
		jarPath, err := filepath.Abs(filepath.Join(jarFolder, file.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get jar path")
		}
		jarPaths = append(jarPaths, jarPath)
	}
	return jarPaths, nil
}

func (r *runner) runJavaDaemon() error {
	daemonClass := "software.amazon.kinesis.multilang.MultiLangDaemon"
	jarPaths, err := getJarPaths(r.pathToJarFolder)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get present working directeory")
	}

	//TODO: check to make sure the properties file exists in currentDir

	paths := append(jarPaths, currentDir)
	classpath := strings.Join(paths, string(os.PathListSeparator))

	args := []string{
		"-cp",
		classpath,
		daemonClass,
		r.pathToPropertiesFile,
	}

	cmd := exec.Command(r.pathToJavaBinary, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout

	r.logger.Println("Starting java daemon process.")
	if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to run command to start java daemon")
	}

	// we need to be able to catch any system exits from the java daemon so we wait for the command
	if err = cmd.Wait(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to run command to wait for java daemon with stderr: %s", stderr.String()))
	}

	r.logger.Println("Java daemon has exited.")
	return nil
}
