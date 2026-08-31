
.DEFAULT_GOAL := build

.PHONY: assets
# STYLE grades the sheet; see scripts/build-tileset64 (none|harbor|retro|chart).
STYLE ?= harbor
assets:
	go run ./scripts/build-tileset64 -style $(STYLE)
	go run ./scripts/bundle-tileset assets/pirate-wars-tileset-64.png cmd/resources/resources.go

.PHONY: icons
icons:
	for s in 32 64 256; do rsvg-convert -w $$s -h $$s assets/icon.svg -o assets/icon-$$s.png; done

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
