JAR_DEPENDENCIES_FOLDER=./jar

install_jars: build_jar_download
	jar-download/jar-download $(JAR_DEPENDENCIES_FOLDER)

build_runner:
	go build -o ./runner/cmd/runner ./runner/cmd

build_sample_app:
	go build -o ./sample/sample ./sample

build_jar_download:
	go build -o ./jar-download/jar-download ./jar-download

build_integration_processor:
	go build -o ./integration-tests/test-app/test_app ./integration-tests/test-app

clean: 
	rm -rf $(JAR_DEPENDENCIES_FOLDER)
	rm -f *log*
	rm -f ./runner/cmd/runner
	rm -f ./sample/sample
	rm -f ./integration-tests/test-app/test_app

run_sample: install_jars build_runner build_sample_app
	./runner/cmd/runner -jar jar -java `which java` -properties sample/sample.properties

run_integ_test: install_jars build_runner build_integration_processor
	docker-compose up -d && \
	go test -count=1 -v ./integration-tests && \
	docker-compose down
