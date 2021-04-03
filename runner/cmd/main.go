package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goguardian/goguardian-go-kcl/runner"
)

const (
	// command line flag names
	propertiesKey = "properties"
	javaKey       = "java"
	jarKey        = "jar"
)

func main() {
	var pathToJavaBinary string
	flag.StringVar(&pathToJavaBinary, javaKey, "", "The path to the java executable e.g. <path>/jdk/bin/java")

	var pathToPropertiesFile string
	flag.StringVar(&pathToPropertiesFile, propertiesKey, "", "The path to the properties file")

	var pathToJarFolder string
	flag.StringVar(&pathToJarFolder, jarKey, "", "The path to the jar dependencies")

	flag.Parse()

	pathToJavaBinary = getVariable(javaKey, pathToJavaBinary)
	pathToPropertiesFile = getVariable(propertiesKey, pathToPropertiesFile)
	pathToJarFolder = getVariable(jarKey, pathToJarFolder)

	r, err := runner.GetRunner(
		runner.WithPathToJavaBinary(pathToJavaBinary),
		runner.WithPathToPropertiesFile(pathToPropertiesFile),
		runner.WithPathToJarFolder(pathToJarFolder),
	)
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	r.RunJavaDaemon()
}

// getVariable is used to fallback to environment variables if the flag was not
// passed in through the command line.
func getVariable(flagKey, value string) string {
	if value == "" {
		value = os.Getenv(strings.ToUpper(flagKey))
	}

	if value == "" {
		fmt.Printf("Must provide %s\n", flagKey)
		flag.Usage()
		os.Exit(1)
	}

	return value
}
