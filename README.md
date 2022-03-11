# mongo-api

## TL;DR
```bash
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
```bash
$ docker build -t mongo-api .
$ docker run -p 8080:8080 -v $PWD/output:/output mongo-api [--serv=<Yourserv>] [--db=<yourDatabase>] [--coll=<yourCollection>]
```

If you start the docker-compose the mongodb listen on localhost

So you can use unit test in `./mongoclient/client_test.go` to test the code with a mock request 
