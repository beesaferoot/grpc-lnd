# LND gRPC Service 


## Requirements 

Before setting up the server/client code, we first need to install the required dependencies.

- install dependencies using go mod ```go mod tidy```

## server setup 

- Before starting the server, we need to have a postgreSQL database up and running
- Export the database URL using `export DATAB_URL="postgres db url"`
- start the server using `go run server/server.go`

## client setup 
- To run the client we simply execute the command `go run client/client.go` on a separate terminal window/tab. Note: ensure the server is *_up and runing_* before runing the client code.  