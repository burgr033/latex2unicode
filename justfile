@default:
    just --list
test:
    go test -v -race -cover ./...
test-coverage:
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
view-test-coverage: test-coverage
    xdg-open coverage.html &

