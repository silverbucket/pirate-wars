
.DEFAULT_GOAL := build

.PHONY: test
test:
	go test -v -cover -coverprofile=c.out ./...

.PHONY: coverage-report-html
coverage-report-html: test
	go tool cover -html=c.out

.PHONY: coverage-report-text
coverage-report-text: test
	go tool cover -func=c.out

.PHONY: build
build:
	go build

.PHONY: clean
clean:
	go clean
	rm -f c.out