all: build test clean


clean:	
	rm -rf *.out
	rm -rf gopns

build:
	go build

test:
	go test ./... 

run:
	go run main.go

coverage:
	./test-coverage.sh

	
