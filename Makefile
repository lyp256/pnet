vet:
	go vet ./...
fmt:
	go fmt ./...

check: fmt vet
	git diff --exit-code

test:
	go test -v --count=1 ./...
