.PHONY: build run

# ========= Vars definitions =========



# ========= Prepare commands =========

tidy:
	go mod tidy
	go clean

# ========= Compile commands =========

build:
	mkdir -p build
	for dir in cmd/*/ ; do \
		name=$$(basename $$dir); \
		go build -ldflags="-s -w" -o ./build/$$name -v ./cmd/$$name/main.go; \
	done

.DEFAULT_GOAL := run
