package runner

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Option func(*Runner)

func WithLogger(logger *log.Logger) Option {
	return func(runner *Runner) {
		runner.logger = logger
	}
}

func WithPathToJavaBinary(p string) Option {
	return func(runner *Runner) {
		runner.pathToJavaBinary = p
	}
}

func WithPathToJarFolder(p string) Option {
	return func(runner *Runner) {
		runner.pathToJarFolder = p
	}
}

func WithPathToPropertiesFile(p string) Option {
	return func(runner *Runner) {
		runner.pathToPropertiesFile = p
	}
}

type Runner struct {
	logger *log.Logger

	pathToJavaBinary     string
	pathToPropertiesFile string
	pathToJarFolder      string
}

func GetRunner(opts ...Option) (*Runner, error) {
	r := &Runner{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.pathToJavaBinary == "" {
		return nil, errors.New("missing path to java binary")
	}

	if r.pathToPropertiesFile == "" {
		return nil, errors.New("missing path to properties folder")
	}

	if r.pathToJarFolder == "" {
		return nil, errors.New("missing path to jar folder")
	}

	return r, nil
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

func (r *Runner) RunJavaDaemon(javaProperties ...string) (*exec.Cmd, error) {
	daemonClass := "software.amazon.kinesis.multilang.MultiLangDaemon"
	jarPaths, err := getJarPaths(r.pathToJarFolder)
	if err != nil {
		return nil, err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get present working directeory")
	}

	// TODO: check to make sure the properties file exists in currentDir

	paths := append(jarPaths, currentDir)
	classpath := strings.Join(paths, string(os.PathListSeparator))

	args := []string{
		"-cp",
		classpath,
		daemonClass,
		r.pathToPropertiesFile,
	}
	args = append(javaProperties, args...)

	cmd := exec.Command(r.pathToJavaBinary, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout

	r.logger.Println("Starting java daemon process.")
	if err = cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to run command to start java daemon")
	}

	return cmd, nil
}
