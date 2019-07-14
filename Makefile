.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	gox -output=build/vte_{{.OS}}_{{.Arch}}
