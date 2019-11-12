

default: build

build:
	rm -rf ./bin && mkdir -p ./bin
	go build -o ./bin/blackjack ./src/main.go

run:
	./bin/blackjack COMPARE 10 6

test:
	go test ./src/cards
	# go test ./src/game
	# go test ./src/player

clean:
	reset
	rm ./out/*.csv