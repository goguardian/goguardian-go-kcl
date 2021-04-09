JAR_DEPENDENCIES_FOLDER=./jar

install_jars:
	rm -rf $(JAR_DEPENDENCIES_FOLDER) && \
	go run jar-download/main.go $(JAR_DEPENDENCIES_FOLDER)

build_sample_app:
	go build -o ./sample/sample ./sample

build_runner:
	go build -o ./runner/cmd/runner ./runner/cmd

build: build_runner build_sample_app

clean: 
	rm -rf log*

build_and_run_sample: install_jars build
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties

run_sample:
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties