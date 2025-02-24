.PHONY: build test clean build-multiarch

build:
	go build -o aws-otp-auth ./cmd/aws-otp-auth

build-multiarch:
	@echo "Building multi-arch binaries for Linux..."
	mkdir -p dist
	  for GOARCH in amd64 arm64; do \
	    BIN=aws-otp-auth-linux-$$GOARCH; \
	    echo "Building $$BIN for linux/$$GOARCH"; \
	    env GOOS=linux GOARCH=$$GOARCH go build -o dist/$$BIN ./cmd/aws-otp-auth; \
	done

test:
	go test -v ./...

clean:
	rm -f aws-otp-auth
	rm -rf dist
