# Battlesloop

This project aims to bring the classic tabletop game Battleship to the CLI using Go. For a(n incomplete) description of the messaging protocol, see [PROTOCOL.md](PROTOCOL.md)

## Usage

1. If you haven't already, [install Go](https://go.dev/doc/install) on your system.

2. Clone the repository.  
`$ git clone https://github.com/scrmbld/battlesloop-go.git`

### Running the server

This step is not necessary if you wish to connect to an existing server.

```
$ cd battlesloop-go/server
$ go run .
```

To connect to a server, see [Running the client](#running-the-client)

### Running the client

```
$ cd battlesloop-go/client
$ go run .
```

Once the client is running, you will be prompted to enter the address of the server. If you wish to run a server yourself, see [Running the server](#running-the-server).

The following demonstrates connecting to a server running on the local machine.  
```
Enter the ip address of the server you would like to connect to:
> localhost
```
