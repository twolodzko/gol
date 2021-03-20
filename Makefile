.DEFAULT_GOAL := run
.PHONY: test coverage clean run

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf *.out *.html

run:
	go run main.go
