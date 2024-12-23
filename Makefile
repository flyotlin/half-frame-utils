.PHONY: run build test test_cleanup

run:
	go run main.go

build:
	./build.sh

test:
	./test.sh

test_cleanup:
	rm -rf *.jpg
