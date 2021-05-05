JAR_DEPENDENCIES_FOLDER=./jar

install_jars:
	go run jar-download/main.go $(JAR_DEPENDENCIES_FOLDER)

build_runner:
	go build -o ./runner/cmd/runner ./runner/cmd

build_sample_app:
	go build -o ./sample/sample ./sample

build_integration_processor:
	go build -o ./integration-tests/test-app/test_app ./integration-tests/test-app

clean: 
	rm -rf *log* && rm -rf $(JAR_DEPENDENCIES_FOLDER)

run_sample: install_jars build_runner build_sample_app
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties

start_localstack:
	docker-compose up -d

stop_localstack:
	docker-compose stop

run_integ_test: install_jars build_integration_processor
	go test -count=1 -v ./integration-tests
