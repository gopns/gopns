all: build test clean

clean:	
	rm -rf *.out
	rm -rf gopns

build:
	go build

test:
	go test ./... 

coverage:
	./test-coverage.sh

	