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

.PHONY: clean
clean:
	@rm -f $(CURDIR)/bin/$(PROJECT)
