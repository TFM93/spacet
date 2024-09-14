.PHONY: docker-compose
docker-compose:
	docker compose up --detach --build

.PHONY: proto
proto:
	buf generate

.PHONY: test
test:
	go test ./...

.PHONY: mocks
mocks:
	mockery