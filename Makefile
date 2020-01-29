# Project params
PROJECT=profile


.PHONY: all
all: clean generate build run

.PHONY: run
run:
	@$(CURDIR)/bin/$(PROJECT)

.PHONY: build
build: 
	@go build -o $(CURDIR)/bin/$(PROJECT) $(CURDIR)/cmd/$(PROJECT)

.PHONY: generate
generate:
	@go generate ./...

.PHONY: test
test:
	@go test ./internal/$(PROJECT) -count=1 -cover

.PHONY: clean
clean:
	@rm -f $(CURDIR)/bin/$(PROJECT)
