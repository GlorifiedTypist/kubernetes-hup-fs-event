TAG = latest
BIN = hup-fs-event

.PHONY: all
all: container

.PHONY: test
test: container 

.PHONY: build
build: main.go
	go build -o $(BIN) -v .

.PHONY: container
container: 
	docker build -t $(BIN) .

.PHONY: clean
clean:
	rm -f $(BIN)

