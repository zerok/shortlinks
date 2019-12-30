all: bin/shortlinks

bin:
	mkdir -p bin

bin/shortlinks: $(shell find . -name '*.go') bin go.mod go.sum
	cd cmd/shortlinks && go build -o ../../$@

cover.out: $(shell find . -name '*.go')
	go test --coverprofile=cover.out ./...

cover.html: cover.out
	go tool cover --html=$< -o $@

clean:
	rm -rf bin cover.html cover.out

test:
	go test ./... -v

docker:
	docker build -t zerok/shortlinks:latest .

.PHONY: test clean all
