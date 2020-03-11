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

.PHONY: all run build generate test lint clean
