# gochat

Simple command-line chat application implemented in Go using the `net` package for network communication.

## Features

* Server: Listens for incoming connections and broadcasts messages to all connected clients.
* Client: Connects to the server and sends messages to the server, receiving broadcasts from other clients.

## Prerequisites

* Go (version 1.20 or higher recommended)

## Usage

1. Clone the project

```
$ git clone https://github.com/danielmourad/gochat
$ cd gochat
```

2. Run the server

```
$ cd cmd/server
$ go run .
```

The server will start listening for connections on port `4700`.

3. Run the client

```
$ cd cmd/client
$ go run . <name>
```

The client will connect to the server running on `127.0.0.1:4700`.

4. Start Chatting
    * Type messages in the client terminal and press `enter` to send them to the server.
    * Messages sent by clients will be broadcasted to all connected clients by the server.
    * type `exit` and hit `enter` to stop the client.
