mock:
	rm -rf mocks
	find . -name "mock_*.go" -type f -delete
	mockery --all --inpackage --quiet --testonly
test:
	go test -race -v ./...
lint:
	staticcheck ./...