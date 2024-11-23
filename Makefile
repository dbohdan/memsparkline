TEST_BINARIES := test/sleep

.PHONY: all
all: memsparkline

.PHONY: clean
clean:
	-rm memsparkline $(TEST_BINARIES)

memsparkline: main.go
	CGO_ENABLED=0 go build

.PHONY: release
release:
	go run script/release.go

.PHONY: test
test: memsparkline $(TEST_BINARIES)
	go test

test/sleep: test/sleep.go
	go build -o $@ test/sleep.go
