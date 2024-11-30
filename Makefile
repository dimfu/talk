SERVER_BIN := bin/server
CLIENT_BIN := bin/client

.PHONY: all build run_server run_client clean

all: build

build: clean $(SERVER_BIN) $(CLIENT_BIN)

$(SERVER_BIN): server
	cd server && go build -o ../$(SERVER_BIN)

$(CLIENT_BIN): client
	cd client && go build -o ../$(CLIENT_BIN)

run_server:
	./$(SERVER_BIN)

run_client:
	./$(CLIENT_BIN)

clean:
	rm -rf bin
	