.PHONY: proto
proto:
	@rm -rf ./pbgo
	@mkdir pbgo
	@protoc -I=./proto --go_out=. ./proto/*

.PHONY: vendor
vendor:
	@go mod tidy
	@go mod vendor

.PHONY: test
test:
	@go clean -testcache
	@go test -v -cover -race ./...

.PHONY: db
db:
	@docker run -d \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=very_secure \
		--name moneybags \
		--restart=always \
		postgres:15-alpine
	@timeout 30 bash -c "until docker exec moneybags pg_isready; do sleep 2; done"

.PHONY: run
run:
	go run cmd/moneybags.go --config local_config.yaml
