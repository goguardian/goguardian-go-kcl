JAR_DEPENDENCIES_FOLDER=./jar

install_jars:
	go run jar-download/main.go $(JAR_DEPENDENCIES_FOLDER)

build_sample_app:
	go build -o ./sample/sample ./sample

build_runner:
	go build -o ./runner/cmd/runner ./runner/cmd

build: build_runner build_sample_app

run_sample: install_jars build
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties
