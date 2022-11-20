# toto-server

## Build and Run
There is `Dockerfile` and `docker-compose.yaml` files to build required image and start docker container.

Container includes:
- toto-server - service with main business logic
- postgres-toto - service of Postgres DB
- redis-toto - service of Redis

```
$ cd toto-server

[//]: # (to create Docker image of toto-server)
$ docker build .

[//]: # (to start Docker container with all services)
$ docker-compose up -d
```

## App structure
```
. toto-server/
├── api/
│   └── swagger.yaml                swagger file for service
├── cmd/
│   ├── toto-server/
│   │   └── main.go                 starting point
├── config/                         config description
├── internal/
│   ├── app/                        main logic of starting service
│   ├── consts/                     constants of service
│   ├── entity/                     entities
│   ├── repository/                 repository layer logic
│   ├── response/                   responses of service
│   └── service/                    service layer logic
├── migrations/                     migration files
├── .env
├── config.yaml                     service config
├── docker-compose.yaml
├── Dockerfile
├── go.mod
│   └── go.sum
└── README.md
```