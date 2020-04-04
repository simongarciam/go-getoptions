.PHONY: test

test:
	go test -coverprofile=coverage.txt -covermode=atomic ./ ./completion/ ./option ./help

view: test
	go tool cover -html=coverage.txt
