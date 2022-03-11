# mongo-api

## TL;DR
```bash
$ git clone https://github.com/sokukata/mongo-api.git
$ cd mongo-api

$ docker-compose build
$ docker compose up

# Go to http://localhost:8080

# For login, credentials are soku:mypwd

# Do Some Test

$ docker-compose down
```

## To Test More

You can start only the Api server to connect it to your own mongo
```bash
$ docker build 
$ docker compose up

# Go to http://localhost:8080
$ docker build -t mongo-api .
$ docker run -p 8080:8080 -v $PWD/output:/output mongo-api [--serv=<Yourserv>] [--db=<yourDatabase>] [--coll=<yourCollection>]
```

If you start the docker-compose the mongodb listen on localhost

I tested on a mac so mount voluume are not straight forward (bind to the docker VM and not the Host) but you can check the container file system:
```bash
$ docker exec -it <container_ID> /bin/sh
>\# ls output
```

So you can use unit test in `./mongoclient/client_test.go` to test the code with a mock request 

If you want to build and run go directly without docker:
```bash
$ go mod download
$ go build
./mongo-api [--serv=<Yourserv>] [--db=<yourDatabase>] [--coll=<yourCollection>]

```