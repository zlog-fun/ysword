.PHONY: build
build: 
	go build -o ysword cmd/ysword/main.go

.PHONY: clean
clean:
	go clean