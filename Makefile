.PHONY: start
start:
	docker-compose up -d --build --remove-orphans

.PHONY: integration-test
e2e:
	go clean -testcache && go test -p 1 ./... -v -tags integration-test
