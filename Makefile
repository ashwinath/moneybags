commit=$(shell git rev-parse HEAD)

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
	@go test -v -coverprofile cover.out -race ./...
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: db
db:
	@docker run -d \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=very_secure \
		--name moneybags \
		--restart=always \
		postgres:16-alpine
	@timeout 30 bash -c "until docker exec moneybags pg_isready; do sleep 2; done"

.PHONY: remove-db
remove-db:
	@docker stop moneybags
	@docker rm moneybags

.PHONY: remake-db
remake-db: remove-db db

.PHONY: db-shell
db-shell:
	PGPASSWORD=very_secure psql -U postgres -d postgres -h 127.0.0.1 -p 5432

.PHONY: run
run:
	go run cmd/moneybags.go --config local_config.yaml

.PHONY: build
build:
	docker build -t $(REGISTRY)/moneybags:$(commit) -t $(REGISTRY)/moneybags:latest .

.PHONY: push
push:
	docker push $(REGISTRY)/moneybags:$(commit)
	docker push $(REGISTRY)/moneybags:latest

.PHONY: container
container: build push

.PHONY: helm-docs
helm-docs:
	@docker run --rm --volume "$$(pwd)/charts/moneybags:/helm-docs" -u $$(id -u) jnorwood/helm-docs:latest
