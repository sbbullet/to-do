.DEFAULT_GOAL := help

.PHONY: help

help:				## Show all commands and their usage
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

server:				## Run the server
	@go run main.go

target01:			## This message will also show up when typing 'make help'
	@echo does something

target02:			## This message will show up too!!!
target02: target00 target01
	@echo does even more