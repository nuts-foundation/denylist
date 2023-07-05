default:

denylist:
	go run cli/main.go | tee out/denylist.jws

test:
	go test -v .

