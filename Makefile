.DEFAULT_GOAL := test
.PHONY: test cov fmt clean repl

test:
	go test ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

clean:
	rm -rf *.out *.html

repl:
	go run main.go
