build:
	go build -o ./build/bar ./version_bar
	go build -o ./build/foo ./version_foo

run: build
	cd ./build && ./foo

.PHONY: build
