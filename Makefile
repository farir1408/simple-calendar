BIN?=./bin/calendar
GOOPENAPIV2_VERSION?=v2@v2.0.0-beta.4
GOOGLEAPIS_VERSION?=v1.15.0
PROJECT_VERSION?=0.0.1
LOCAL_GOOSE?=linux


.PHONY: build
build:
	GOOS=$(LOCAL_GOOSE) go build -o $(BIN) ./cmd/calendar-api


.PHONY: lint
lint:
	golangci-lint run --config=.golangci.yml ./...


.PHONY: .migrate
.migrate:
	goose -dir scripts/migrations postgres "user=calendar-api-user dbname=calendar-api-db password=calendar-api-password host=localhost port=5432 sslmode=disable" up
	sleep 1


.PHONY: .migrate-dry-run
.migrate-dry-run:
	goose -dir scripts/migrations postgres "user=calendar-api-user dbname=calendar-api-db password=calendar-api-password host=localhost port=5432 sslmode=disable" status


.PHONY: test-unit
test-unit:
	go test -race -v ./... -tags unit -count 1


.PHONY: test-integration
test-integration:
	go test -race -v ./... -tags integration -count 1


.PHONY: .run-test-db
.run-test-db:
	docker-compose -f build/docker-compose.yml up -d calendar_db


.PHONY: .stop-test-db
.stop-test-db:
	docker-compose -f build/docker-compose.yml down --rmi=all -v


.PHONY: run-test-integration
run-test-integration: .run-test-db .migrate test-integration .stop-test-db



.PHONY: install
install:
	go get \
    		github.com/golang/protobuf/protoc-gen-go \
    		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    		github.com/rakyll/statik


.PHONY: generate
generate:
	protoc -I. \
			-I$(GOPATH)/pkg/mod \
			-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GOOGLEAPIS_VERSION)/third_party/googleapis \
			-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/$(GOOPENAPIV2_VERSION) \
			--go_out=plugins=grpc:pkg \
			--grpc-gateway_out=logtostderr=true:pkg \
			--openapiv2_out=logtostderr=true:third_party/openapi \
			api/calendar/calendar.proto

	statik -src=./third_party/openapi/api/calendar -dest=./internal/pkg/ -ns openapi_definition


.PHONY: image-build
image-build:
	docker build -t farir/simple-calendar:$(PROJECT_VERSION) -f build/Dockerfile .


.PHONY: .docker-run-local
.docker-run-local:
	docker run -p 8091:8091 -p 8092:8092 --name=calendar-service farir/simple-calendar:$(PROJECT_VERSION) "--run-local"


.PHONY: .generate-statik
.generate-statik:
	statik -src=./third_party/openapi/api/calendar -dest=./internal/pkg/ -ns openapi_definition


.PHONY: update-swagger-ui
update-swagger-ui:
	go get -u github.com/rakyll/statik

	rm -rf /tmp/swagger-ui
	git clone https://github.com/swagger-api/swagger-ui.git /tmp/swagger-ui

	mkdir /tmp/swagger-ui/calendar; \
		cat /tmp/swagger-ui/dist/index.html | perl -pe 's/https?:\/\/petstore.swagger.io\/v2\///g' | sed 's/swagger.json/calendar.swagger.json/' > /tmp/swagger-ui/calendar/index.html; \
		cp /tmp/swagger-ui/dist/oauth2-redirect.html /tmp/swagger-ui/calendar; \
		cp /tmp/swagger-ui/dist/*.js /tmp/swagger-ui/calendar; \
		cp /tmp/swagger-ui/dist/*.css /tmp/swagger-ui/calendar; \
		cp /tmp/swagger-ui/dist/*.png /tmp/swagger-ui/calendar

	cp -a /tmp/swagger-ui/calendar/. ./third_party/openapi/api/calendar/
	statik -src=./third_party/openapi/api/calendar -dest=./internal/pkg/ -ns openapi_definition

	rm -rf /tmp/swagger-ui
