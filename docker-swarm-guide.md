# Deploy - Docker Swarm Guide

## Building microservices (Docker image)

1. Building logger-service

```
$ cd logger-service
$ docker build -f logger-service.dockerfile -t obededoreto/logger-service:1.0.0 .
$ docker push obededoreto/logger-service:1.0.0
```

2. Building mail-service

```
$ cd mail-service
$ docker build -f mail-service.dockerfile -t obededoreto/mail-service:1.0.0 .
$ docker push obededoreto/mail-service:1.0.0
```

3. Building authenticator-service

```
$ cd authenticator-service
$ docker build -f authentication-service.dockerfile -t obededoreto/authentication-service:1.0.0 .
$ docker push obededoreto/authentication-service:1.0.0
```

4. Building broker-service

```
$ cd broker-service
$ docker build -f broker-service.dockerfile -t obededoreto/broker-service:1.0.1 .
$ docker push obededoreto/broker-service:1.0.1
```

5. Building listener-service

```
$ cd listener-service
$ docker build -f listener-service.dockerfile -t obededoreto/listener-service:1.0.0 .
$ docker push obededoreto/listener-service:1.0.0
```

## Building frontends (Docker image)

1. Building front end application

```
$ cd front-end
$ docker build -f front-end.dockerfile -t obededoreto/front-end:1.0.2 .
$ docker push obededoreto/front-end:1.0.2
```

2. Building caddy proxy

```
$ cd project
$ docker build -f caddy.dockerfile -t obededoreto/micro-caddy:1.0.0 .
$ docker push obededoreto/micro-caddy:1.0.0
```

3. Configure etc hosts

```sh
sudo vim /etc/hosts

127.0.0.1 localhost backend
::1       localhost backend
```

## Creating the stack file (docker-stack.yaml)

```yaml
services:
  authentication-service:
    image: obededoreto/authentication-service:1.0.0
    restart: on-failure
    environment:
      LOGGER_SERVICE_URL: "http://logger-service"
      DSN: "host=postgres port=5432 dbname=users user=postgres password=password sslmode=disable timezone=UTC connect_timeout=5"
    ports:
      - "8081:80"
    deploy:
      mode: replicated
      replicas: 1

  broker-service:
    image: obededoreto/broker-service:1.0.0
    restart: on-failure
    environment:
      AUTH_SERVICE_URL: "http://authentication-service"
      LOGGER_SERVICE_URL: "http://logger-service"
      MAIL_SERVICE_URL: "http://mail-service"
      RABBITMQ_URL: amqp://guest:guest@rabbitmq:5672/
    ports:
      - "8080:80"
    deploy:
      mode: replicated
      replicas: 1

  listener-service:
    image: obededoreto/listener-service:1.0.0
    restart: on-failure
    environment:
      AUTH_SERVICE_URL: "http://authentication-service"
      LOGGER_SERVICE_URL: "http://logger-service"
      MAIL_SERVICE_URL: "http://mail-service"
      RABBITMQ_URL: amqp://guest:guest@rabbitmq:5672/
    ports:
      - "8084:80"
    deploy:
      mode: replicated
      replicas: 1

  logger-service:
    image: obededoreto/logger-service:1.0.0
    restart: on-failure
    environment:
      MONGO_URL: "mongodb://mongo:27017"
    ports:
      - "8082:80"
    deploy:
      mode: replicated
      replicas: 1

  mail-service:
    image: obededoreto/mail-service:1.0.0
    restart: on-failure
    environment:
      MAIL_DOMAIN: "localhost"
      MAIL_HOST: "mail-client"
      MAIL_PORT: 1025
      MAIL_ENCRYPTION: "none"
      MAIL_USERNAME: ""
      MAIL_PASSWORD: ""
      MAIL_FROM_NAME: "Jhon Smith"
      MAIL_FROM_ADDRESS: "jhon.smith@example.com"
    ports:
      - "8083:80"
    deploy:
      mode: replicated
      replicas: 1

  # Third services
  rabbitmq:
    image: "rabbitmq:3.9-alpine"
    ports:
      - "5672:5672"
      - "15672:15672"
    deploy:
      mode: global

  pgadmin:
    image: dpage/pgadmin4:latest
    environment:
      PGADMIN_DEFAULT_EMAIL: dba@email.com
      PGADMIN_DEFAULT_PASSWORD: q1w2e3r4
    deploy:
      mode: global
    ports:
      - 8085:80

  postgres:
    image: "postgres:14.0"
    ports:
      - "5432:5432"
    deploy:
      mode: global
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: users
    volumes:
      - ./db-data/postgres/:/var/lib/postgresql/data/

  mongo:
    image: "mongo:4.2.17-bionic"
    ports:
      - "27017:27017"
    deploy:
      mode: global
    environment:
      MONGO_INITDB_DATABASE: logs
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
    volumes:
      - ./db-data/mongo/:/data/db/

  mail-client:
    image: mailhog/mailhog:latest
    deploy:
      mode: global
    ports:
      - "8025:8025"
      - "1025:1025"
```

## Initializing Swarm Cluster

1. Initializing cluster

```sh
$ docker swarm init
# docker swarm init --advertise-addr 192.168.1.12

# Copy the output command to register nodes in cluster
# or generate a new comand token
$ docker swarm join-token manager
$ docker swarm join-token worker
```

2. Applying stack

```sh
# docker stack deploy -c <YAML> <PREFIX>
$ docker stack deploy -c docker-stack.yaml myapp
```

3. Show services

```sh
$ docker service ls
```

3. Show tasks from service

```sh
$  docker service ps myapp_authentication-service
```

4. Scale services

```sh
# docker service scale <SERVICE_NAME> = <QNT>
$  docker service scale myapp_listener-service=3
```

5. Updating service image

```sh
# docker service update --image <IMAGE> <SERVICE_NAME>
$  docker service update --image obededoreto/logger-service:new-version myapp_logger-service
```

6. Stopping docker swarm

```sh
# $ docker service myapp_broker-service=0

# Removing stack apps
$ docker stack rm myapp

# Leaving from cluster
$ docker swarm leave --force
```
