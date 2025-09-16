.PHONY: build run

# ========= Vars definitions =========

go_ldflags = "-s -w"

# ========= Prepare commands =========

clear-build:
	rm -rf ./build || echo "build folder did not exists yet!"
	mkdir ./build

tidy:
	go mod tidy
	go clean

# ========= Compile commands =========

build: clear-build
	mkdir -p build

	GOOS=linux
	for dir in cmd/*/ ; do \
		name=$$(basename $$dir); \
		go build -ldflags=$(go_ldflags) -o ./build/linux/$$name -v ./cmd/$$name/main.go; \
	done

govulncheck-scan:
	@for binary in ./build/linux/*; do \
		echo "checking binary go: $$binary"; \
		govulncheck -mode binary $$binary; \
	done

gosec-scan:
	@for subproj in ./cmd/*; do \
		echo "checking sub project go: $$subproj"; \
		gosec $$subproj/...; \
	done

build-prod: gosec-scan build govulncheck-scan

lic:
	./set_lic.sh

.DEFAULT_GOAL := run
