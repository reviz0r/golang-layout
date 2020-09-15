# Project params
PROJECT=profile

all: clean build run

run:
	@$(CURDIR)/bin/$(PROJECT)

build: 
	@go build -o $(CURDIR)/bin/$(PROJECT) $(CURDIR)/cmd/$(PROJECT)

test:
	@go test ./internal/$(PROJECT) -count=1 -cover

lint:
	@go vet ./...

generate:
	@go generate ./...

clean:
	@rm -f $(CURDIR)/bin/$(PROJECT)

migrate-up:
	@migrate -source file://./migrations -database 'postgres://postgres@localhost:5432/golang-layout?sslmode=disable' up

migrate-down:
	@migrate -source file://./migrations -database 'postgres://postgres@localhost:5432/golang-layout?sslmode=disable' down

.PHONY: all run build generate test lint clean migrate-up migrate-down
