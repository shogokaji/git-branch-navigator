.PHONY: build run

# Go build command
build:
	go build -o bin/app cmd/main.go

# Run the application
run:
	go run main.go
