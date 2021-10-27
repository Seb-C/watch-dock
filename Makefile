.PHONY: test run

build:
	go build -o ./bin/watch-dock ./cmd/main.go

run: build
	export WDBIN="$$(realpath ./bin/watch-dock)" && cd ./example/hello-world.go && $$WDBIN
