all: smallchat-server smallchat-client go-server go-client
CFLAGS=-O2 -Wall -W -std=c99

smallchat-server: smallchat-server.c chatlib.c
	$(CC) smallchat-server.c chatlib.c -o smallchat-server $(CFLAGS)

smallchat-client: smallchat-client.c chatlib.c
	$(CC) smallchat-client.c chatlib.c -o smallchat-client $(CFLAGS)

clean:
	rm -f smallchat-server
	rm -f smallchat-client
	rm -f go-server
	rm -f go-client

install:
	go mod tidy

go-server: smallchat-server.go chatlib.go
	go build -o go-server smallchat-server.go chatlib.go

go-client: smallchat-client.go chatlib.go
	go build -o go-client smallchat-client.go chatlib.go
