# Project params
PROJECT=profile


.PHONY: all
all: clean build
	@$(CURDIR)/bin/$(PROJECT)

.PHONY: build
build:
	@go build -o $(CURDIR)/bin/$(PROJECT) $(CURDIR)/cmd/$(PROJECT)

.PHONY: clean
clean:
	@rm -f $(CURDIR)/bin/$(PROJECT)
