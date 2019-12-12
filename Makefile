# Project params
PROJECT=profile


.PHONY: all
all: clean build
	@$(CURDIR)/bin/$(PROJECT)

.PHONY: build
build: generate
	@go build -o $(CURDIR)/bin/$(PROJECT) $(CURDIR)/cmd/$(PROJECT)

.PHONY: generate
generate:
	@go generate ./...

.PHONY: clean
clean:
	@rm -f $(CURDIR)/bin/$(PROJECT)
