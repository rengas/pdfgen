.PHONY: start
start:
	docker-compose up -d --build --remove-orphans

.PHONY: integration-test
e2e:
	go clean -testcache && go test -p 1 ./... -v -tags integration-test

.PHONY: generate-docs
generate-docs:
	swag init -dir cmd/api --o docs/api --ot yaml --instanceName api --pd true

.PHONY: generate-model
generate-model: generate-docs
	sh ./scripts/generate-model

.PHONY: serve-docs
serve-docs:
	python3 -m http.server 9000 --directory docs/api


.PHONY: e2e
e2e:
	go clean -testcache && go test -p 1 ./... -v -tags e2e