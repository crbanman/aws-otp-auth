.PHONY: build test clean build-multiarch

build:
	go build -o aws-otp-auth ./cmd/aws-otp-auth

build-multiarch:
	@echo "Building multi-arch binaries for Linux and macOS..."
	mkdir -p dist
	for GOOS in linux darwin; do \
	  for GOARCH in amd64 arm64; do \
	    BIN=aws-otp-auth-$$GOOS-$$GOARCH; \
	    echo "Building $$BIN for $$GOOS/$$GOARCH"; \
	    env GOOS=$$GOOS GOARCH=$$GOARCH go build -o dist/$$BIN ./cmd/aws-otp-auth; \
	  done; \
	done

test:
	go test -v ./...

clean:
	rm -f aws-otp-auth
	rm -rf dist
