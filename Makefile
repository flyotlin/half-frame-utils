.PHONY: run build test test_cleanup

run:
	go run main.go

build:
	go build main.go

test:
	./test.sh

test_cleanup:
	rm -rf *.jpg
