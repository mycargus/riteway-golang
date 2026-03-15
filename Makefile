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
	./scripts/release.sh $(VERSION)
