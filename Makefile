.PHONY: denylist test

default:

denylist:
	go run cli/main.go | tee denylist/denylist.jws

test:
	go test -v .

