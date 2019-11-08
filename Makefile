

default: build

build:
	rm -rf ./bin && mkdir -p ./bin
	go build -o ./bin/blackjack ./src/main.go

run:
	./bin/blackjack

test:
	go test ./src/cards
	# go test ./src/game
	# go test ./src/player

clean:
	reset
	rm -rf ./out