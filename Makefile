vet:
	@echo go vet project
	@go list ./... |grep -v vendor | xargs go vet
fmt:
	@echo go fmt project
	@go list ./... |grep -v vendor | xargs go fmt

tidy:
	go mod tidy

vendor:
	go mod vendor

golint:
	@echo golint project
	@go list ./... |grep -v vendor | xargs golint

pretty: fmt vet golint tidy vendor

check: pretty
	git diff --exit-code

test:
	 go list ./... |grep -v vendor | xargs go test -v --cover --count=1

.PHONY: fmt vet tidy vendor pretty check test
