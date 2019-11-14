
default: build

build:
	rm -rf ./bin && mkdir -p ./bin
	go build -o ./bin/blackjack ./src/main.go

test:
	go test ./src/cards
	# go test ./src/game
	# go test ./src/player

clean:
	reset
	rm ./out/*.csv

# run modes

run: human

learn:
	./bin/blackjack LEARN 10 6

story:
	./bin/blackjack STORY 10 6 > ./out/log.txt

compare:
	./bin/blackjack COMPARE 10 6

human:
	./bin/blackjack HUMAN 10 1


