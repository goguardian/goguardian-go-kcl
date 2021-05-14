JAR_DEPENDENCIES_FOLDER=./jar

download_jars:
	rm -rf $(JAR_DEPENDENCIES_FOLDER) && \
	go build -o ./jar-download/jar-download ./jar-download && \
	./jar-download/jar-download $(JAR_DEPENDENCIES_FOLDER)

build_sample_app:
	go build -o ./sample/sample ./sample

build_runner:
	go build -o ./runner/cmd/runner ./runner/cmd

build: build_runner build_sample_app

clean: 
	rm -rf *log*

run_sample:
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties

build_and_run_sample: build run_sample

download_build_and_run_sample: download_jars build_and_run_sample
