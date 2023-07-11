default:

denylist:
	go run cli/main.go | tee denylist-out/denylist.jws

test:
	go test -v .

