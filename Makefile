.PHONY: docker-compose
docker-compose:
	docker compose up --detach

.PHONY: proto
proto:
	buf generate

.PHONY: test
test:
	go test ./...