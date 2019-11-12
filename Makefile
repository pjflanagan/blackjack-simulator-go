

default: build

build:
	rm -rf ./bin && mkdir -p ./bin
	go build -o ./bin/blackjack ./src/main.go

run: human

compare:
	./bin/blackjack COMPARE 10 6 > out/log.txt

human:
	./bin/blackjack COMPETE 10 1

test:
	go test ./src/cards
	# go test ./src/game
	# go test ./src/player

clean:
	reset
	rm ./out/*.csv