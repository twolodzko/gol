.DEFAULT_GOAL := test
.PHONY: test cov cycl cogn fmt clean repl

test:
	go test ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

cycl:
	# go get github.com/fzipp/gocyclo/cmd/gocyclo
	gocyclo -top 10 .

cogn:
	# go get github.com/uudashr/gocognit/cmd/gocognit
	gocognit -top 10 .

fmt:
	go fmt ./...

clean:
	rm -rf *.out *.html

repl:
	go run main.go
