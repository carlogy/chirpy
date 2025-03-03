run:
	GOCACHE=of && go build -o out && ./out

test:
	go test ./...

testAll:

	go test ./... -v
