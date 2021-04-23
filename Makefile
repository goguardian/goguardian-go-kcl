JAR_DEPENDENCIES_FOLDER=./jar

install_jars:
	rm -rf $(JAR_DEPENDENCIES_FOLDER) && \
	go run jar-download/main.go $(JAR_DEPENDENCIES_FOLDER)

build_sample_app:
	go build -o ./sample/sample ./sample

build_runner:
	go build -o ./runner/cmd/runner ./runner/cmd

build_integration_processor:
	go build -o ./integration_tests/test_app ./integration_tests

build: build_runner build_sample_app build_integration_processor

clean: 
	rm -rf *log*

run_sample:
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties

build_and_run_sample: build run_sample

download_build_and_run_sample: install_jars build_and_run_sample

start_localstack:
	docker-compose up -d

stop_localstack:
	docker-compose stop

run_integ_test: build start_localstack
	go test -count=1 -v ./integration_tests
