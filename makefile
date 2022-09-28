CMD = ./cmd

GO = go

compile:
	$(GO) build -o ./bin/app $(CMD)
	cp -r ./res ./bin/

./bin:
	mkdir -p bin

run:
	$(GO) run $(CMD)

clean:
	rm -r ./bin/*

all:
	@echo compile, run, clean
