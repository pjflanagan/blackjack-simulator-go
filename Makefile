

default: build

build:
	rm -rf ./bin && mkdir -p ./bin
	go build -o ./bin/blackjack ./src/main.go

run:
	./bin/blackjack