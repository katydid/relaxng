build:
	go build ./...

test:
	go test .

gofmt:
	gofmt -l -s -w .

vet:
	go vet ./...

travis:
	make test
	make gofmt
	make vet
	git diff --exit-code
