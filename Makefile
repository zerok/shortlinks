all: bin/shortlinks

bin:
	mkdir -p bin

pkged.go: $(shell find ./migrations/ -name '*.sql')
	pkger

bin/shortlinks: $(shell find . -name '*.go') bin go.mod go.sum pkged.go
	cd cmd/shortlinks && go build -o ../../$@

clean:
	rm -rf bin

test:
	go test ./... -v

.PHONY: test clean all
