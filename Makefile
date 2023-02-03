PHONY: build run

build:
	go build $$(go list ./...)

run:
	go run $$(go list ./...)
