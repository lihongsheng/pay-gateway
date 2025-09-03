#GOHOSTOS:=$(shell go env GOHOSTOS)
#GOPATH:=$(shell go env GOPATH)
#VERSION=$(shell git describe --tags --always)
#
#ifeq ($(GOHOSTOS), windows)
#	#the `find.exe` is different from `find` in bash/shell.
#	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
#	#changed to use git-bash.exe to run find cron or other cron friendly, caused of every developer has a Git.
#	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
#	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
#	API_PROTO_FILES=$(shell $(Git_Bash) -c " cd ./$(NAME)/interface && find api -name *.proto")
#	ENUM_PROTO_FILES=$(shell $(Git_Bash) -c " cd ./shared && find enum -name *.proto")
#else
#	API_PROTO_FILES=$(shell cd ./interface && find proto -name *.proto)
#	ENUM_PROTO_FILES=$(shell find ./shared/enum -name *.proto)
#endif
#
#.PHONY: init
#init:
#	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
#	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
#	go install github.com/envoyproxy/protoc-gen-validate@latest
#	go install github.com/google/wire/cmd/wire@latest
#	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
#	go install github.com/golang/mock/mockgen@latest
#
#.PHONY: api
## generate api proto
#api:
#	cd ./interface && protoc --proto_path=./proto \
#	       --proto_path=./../../third_party \
#	       --proto_path=./../../ \
# 	       --go_out=paths=source_relative:./proto \
# 	       --go-http_out=paths=source_relative:./proto \
# 	       --go-grpc_out=paths=source_relative:./proto \
# 	       --validate_out="lang=go,paths=source_relative:./proto" \
#	       --openapi_out=fq_schema_naming=true,default_response=false,naming=proto:. \
#	       $(API_PROTO_FILES)
#
#.PHONY: errors
#errors:
#	protoc --proto_path=. --proto_path=./third_party --go_out=./ --go-errors_out=paths=./ ./errors/errors.proto
#.PHONY: enum
#enum:
#	protoc --proto_path=. \
#            --proto_path=./third_party \
#            --go_out=paths=source_relative:. \
#            --go-errors_out=paths=source_relative:. \
#             $(ENUM_PROTO_FILES)
#.PHONY: rpc
#rpc:
#	go build -o ./bin/server ./cmd/rpc/main.go ./cmd/rpc/wire_gen.go
#
#
#
##.PHONY: config
### generate internal proto
##config:
##	protoc --proto_path=. \
##	       --proto_path=./third_party \
## 	       --go_out=paths=source_relative:./shared \
##	       $(INTERNAL_PROTO_FILES)
#
## show help
#help:
#	@echo ''
#	@echo 'Usage:'
#	@echo ' make [target]'
#	@echo ''
#	@echo 'Targets:'
#	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
#	helpMessage = match(lastLine, /^# (.*)/); \
#		if (helpMessage) { \
#			helpCommand = substr($$1, 0, index($$1, ":")-1); \
#			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
#			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
#		} \
#	} \
#	{ lastLine = $$0 }' $(MAKEFILE_LIST)
#
#.DEFAULT_GOAL := help
#
##.PHONY: local-rpc
##local-rpc:
##	go run cmd/rpc/main.go cmd/rpc/wire_gen.go -conf configs/ -env local
##.PHONY: local-cron
##local-cron:
##	go run cmd/cron/main.go cmd/cron/wire_gen.go -conf configs/ -env local
#
#.PHONY: git
#git:
#	gofmt -l -w -s ./application
#	gofmt -l -w -s ./cmd
#	gofmt -l -w -s ./shared
#	gofmt -l -w -s ./interface
#	gofmt -l -w -s ./internal
#	git add ./application
#	git add ./cmd
#	git add ./shared
#	git add ./interface
#	git add ./internal
#	git add ./Makefile
#	git add ./README.md
#	git status
