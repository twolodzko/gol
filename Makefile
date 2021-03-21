.DEFAULT_GOAL := run
.PHONY: test cov fmt clean run

test:
	go test ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

clean:
	rm -rf *.out *.html

run:
	go run main.go
