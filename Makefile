.PHONY: build image clean

build:
	cd cmd && GOOS=linux GOARCH=amd64 go build -o wave-generator main.go

image: build
	docker build -t wave-generator ./cmd

clean:
	rm -f cmd/wave-generator 