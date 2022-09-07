.DEFAULT_GOAL := help

.PHONY: help

help:				## Show all commands and their usage
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

server:				## Run the server
	@go run main.go

test:			## Run all the tests
	@go test -v -cover ./...