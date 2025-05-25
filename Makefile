build:
		@go build -o bin/icsMaker

run: build
		@./bin/icsMaker

test: 
		@go test -v ./...