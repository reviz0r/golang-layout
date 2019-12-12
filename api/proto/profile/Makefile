# System params
GOPATH=/Users/$(shell whoami)/go
GATEWAY_VERSION=1.12.1
GATEWAY_PATH=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v$(GATEWAY_VERSION)

# Project params
PROJECT=profile
TARGET_DIR=../../../pkg/$(PROJECT)/

# File lists
PROTO_GO_IN=$(wildcard *.proto)
PROTO_GW_IN=$(wildcard *.api.proto)
PROTO_GO_OUT=$(join $(addsuffix $(TARGET_DIR), $(dir $(PROTO_GO_IN))), $(notdir $(PROTO_GO_IN:.proto=.pb.go)))
PROTO_GW_OUT=$(join $(addsuffix $(TARGET_DIR), $(dir $(PROTO_GW_IN))), $(notdir $(PROTO_GW_IN:.proto=.pb.gw.go)))
PROTO_SW_OUT=$(notdir $(PROTO_GW_IN:.proto=.swagger.json))


.PHONY: all
all: $(PROTO_GO_OUT) $(PROTO_GW_OUT) $(PROTO_SW_OUT)

.PHONY: clean
clean:
	# $(info Cleaning files generated from $(PROTO_GO_IN))
	@rm -f $(PROTO_GO_OUT) $(PROTO_GW_OUT) $(PROTO_SW_OUT)

# Rule for compiling protobuf + grpc
$(TARGET_DIR)%.pb.go : %.proto
	$(info Generating proto + grpc from $<)
	@protoc \
		--proto_path=/usr/local/include \
		--proto_path=$(GATEWAY_PATH) \
		--proto_path=$(GATEWAY_PATH)/third_party/googleapis \
		--proto_path=. \
		--go_out=plugins=grpc:$(GOPATH)/src \
		$<

# Rule for compiling grpc-gateway
$(TARGET_DIR)%.pb.gw.go : %.proto
	$(info Generating grpc-gateway from $<)
	@protoc \
		--proto_path=/usr/local/include \
		--proto_path=$(GATEWAY_PATH) \
		--proto_path=$(GATEWAY_PATH)/third_party/googleapis \
		--proto_path=. \
		--grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
		$<

# Rule for compiling swagger
%.swagger.json: %.proto
	$(info Generating swagger from $(PROTO_GW_IN))
	@protoc \
		--proto_path=/usr/local/include \
		--proto_path=$(GATEWAY_PATH) \
		--proto_path=$(GATEWAY_PATH)/third_party/googleapis \
		--proto_path=. \
		--swagger_out=logtostderr=true:. \
		$(PROTO_GW_IN)
