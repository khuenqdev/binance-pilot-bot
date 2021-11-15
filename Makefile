.PHONY: install gotest build serve clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
DATE?=$(shell date -u +%s)

export TAG
export DATE

test:
	go build -o bin/binance .
	./bin/binance test

dev:
	go get .
	go build -ldflags "-X github.com/khuenqdev/binance-pilot-bot/cmd.BuiltAt=$(DATE) -X github.com/khuenqdev/binance-pilot-bot/cmd.GitCommitHash=$(TAG)" -o ./bin/binance .
	docker build -t localhost:5000/binance .
	docker push localhost:5000/binance