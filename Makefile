all: bin/shortlinks

bin:
	mkdir -p bin

bin/shortlinks: $(shell find . -name '*.go') bin go.mod go.sum
	cd cmd/shortlinks && go build -o ../../$@

clean:
	rm -rf bin

test:
	go test ./... -v

docker:
	docker build -t zerok/shortlinks:latest .

.PHONY: test clean all
