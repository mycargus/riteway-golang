.PHONY: test test-scripts fmt vet check release

test:
	go test -count=1 -race ./...

test-scripts:
	bash scripts/test-release.sh

fmt:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "gofmt: unformatted files:"; \
		echo "$$out"; \
		exit 1; \
	fi

vet:
	go vet ./...

check: fmt vet test test-scripts

release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is not set. Usage: make release VERSION=v0.1.0"; \
		exit 1; \
	fi
	./scripts/release.sh $(VERSION)
