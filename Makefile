.PHONY: build test cov staticcheck vet cycl cogn benchmarks fmt clean repl

build:
	go build

test: staticcheck vet
	go test ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

staticcheck:
	# go get honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

vet:
	go vet ./...

cycl:
	# go get github.com/fzipp/gocyclo/cmd/gocyclo
	gocyclo -top 10 .

cogn:
	# go get github.com/uudashr/gocognit/cmd/gocognit
	gocognit -top 10 .

benchmarks: build
	go test -bench=. ./...

fmt:
	go fmt ./...

clean:
	rm -rf *.out *.html ./gol

repl: build
	./gol
