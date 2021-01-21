
TARGET=pneuma

test:
	go test ./...

verify:
	golint ./...

build:
	go build -o bin/$(TARGET) $(TARGET).go

run:
	go run $(TARGET).go
