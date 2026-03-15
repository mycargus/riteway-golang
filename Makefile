.PHONY: test fmt vet check

test:
	go test -count=1 -race ./...

fmt:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "gofmt: unformatted files:"; \
		echo "$$out"; \
		exit 1; \
	fi

vet:
	go vet ./...

check: fmt vet test
