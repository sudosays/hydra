OUTDIR=bin
TARGET=pneuma

build:
	go build -o ./$(OUTDIR)/$(TARGET) $(TARGET).go

clean:
	rm -fr ./$(OUTDIR)/*

run:
	go run $(TARGET).go

test:
	go test ./...

verify:
	golint ./...

prep: test verify
	go mod tidy
	
	

