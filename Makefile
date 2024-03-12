setup:
	go mod tidy

run:
	go run . ./test/test.txt

build:
	go build -ldflags "-s -w" -o bin/ .

clean:
	rm -rf bin data