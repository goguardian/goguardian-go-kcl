package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var pathToPropertiesFile string
	flag.StringVar(&pathToPropertiesFile, "properties", "", "The path to the properties file")

	var pathToJavaBinary string
	flag.StringVar(&pathToJavaBinary, "java", "", "The path to the java executable e.g. <path>/jdk/bin/java")

	var pathToJarFolder string
	flag.StringVar(&pathToJarFolder, "jar", "", "The path to the jar dependencies")

	flag.Parse()

	checkStringNotEmpty("properties file", pathToPropertiesFile)
	checkStringNotEmpty("java path", pathToJavaBinary)
	checkStringNotEmpty("jar path", pathToJarFolder)

	r := &runner{
		pathToJavaBinary:     pathToJavaBinary,
		pathToPropertiesFile: pathToPropertiesFile,
		pathToJarFolder:      pathToJarFolder,
		logger:               log.New(os.Stdout, "", log.LstdFlags),
	}
	r.runJavaDaemon()
}

func checkStringNotEmpty(name, str string) {
	if str == "" {
		fmt.Printf("Must provide %s\n", name)
		flag.Usage()
		os.Exit(1)
	}
}
