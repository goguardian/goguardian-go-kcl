JAR_DEPENDENCIES_FOLDER=./jar

install_jars:
	go run jar-download/main.go $(JAR_DEPENDENCIES_FOLDER)
