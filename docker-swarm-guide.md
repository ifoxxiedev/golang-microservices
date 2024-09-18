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

$ docker build -f caddy.production.dockerfile -t obededoreto/micro-caddy-production:1.0.0 .
$ docker push obededoreto/micro-caddy-production:1.0.0
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
  caddy:
    image: obededoreto/micro-caddy:1.0.0
    restart: on-failure
    ports:
      - "8079:80"
      - "80:80"
      - "443:443"
    deploy:
      mode: replicated
      replicas: 1
    volumes:
      - caddy_data:/data
      - caddy_config:/config

  authentication-service:
    image: obededoreto/authentication-service:1.0.0
    restart: on-failure
    environment:
      LOGGER_SERVICE_URL: "http://logger-service"
      DSN: "host=postgres port=5432 dbname=users user=postgres password=password sslmode=disable timezone=UTC connect_timeout=5"
    ports:
      - "8080:80"
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

  front-end:
    image: obededoreto/front-end:1.0.2
    restart: on-failure
    ports:
      - "8081:8081"
    environment:
      BROKER_URL: "http://broker"
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
    image: "postgres:14.2"
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

volumes:
  caddy_data:
    external: true
  caddy_config:
```

## Deploying in swarm cluster (LOCAL)

### Initializing Swarm Cluster

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

## Preparing Deploy in Production

### Preparing nodes in [linode](https://www.linode.com/)

1. Configuring VM Node-1

```sh
# Connecting in VM one
$ ssh root@ip-node-1

# Configure user
$ adduser tcs  # type the password
$ usermod -aG sudo tcs # add tcs user in sudo group

# Configure firewall
$ ufw allow ssh
$ ufw allow http
$ ufw allow https

$ ufw allow 2377/tcp
$ ufw allow 7946/tcp
$ ufw allow 7946/udp
$ ufw allow 4789/udp
$ ufw allow 8025/tcp

$ ufw enable
$ ufw status

# Configure hostname
sudo hostnamectl set-hostname node-1
# exit and log via ssh again

sudo vi /etc/hosts
# Copy the IP Address on Linode platform
172.105.9.56 node-1
172.105.9.57 node-2
```

2. Configuring VM Node-2

```sh
# Connecting in VM one
$ ssh root@ip-node-2

# Configure user
$ adduser tcs  # type the password
$ usermod -aG sudo tcs # add tcs user in sudo group

# Configure firewall
$ ufw allow ssh
$ ufw allow http
$ ufw allow https

$ ufw allow 2377/tcp
$ ufw allow 7946/tcp
$ ufw allow 7946/udp
$ ufw allow 4789/udp
$ ufw allow 8025/tcp

$ ufw enable
$ ufw status

# Configure hostname
sudo hostnamectl set-hostname node-2
# exit and log via ssh again

sudo vi /etc/hosts
# Copy the IP Address on Linode platform
172.105.9.56 node-1
172.105.9.57 node-2
```

3. Install docker in the server

```sh
# Connecting in VM one
$ ssh root@ip-node-1
$ sudo tsc #  type the password

# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
$ echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

$ sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

$ witch docker

$ sudo usermod -aG docker tcs
```

4. Configure DNS (Registro BR, GoDaddy, AWS Route530)

```
A   swarm.your-domain   172.105.9.56
A   swarm.your-domain   172.105.9.57
A   node-1.your-domain  172.105.9.56
A   node-2.your-domain  172.105.9.57

CNAME  broker  swarm.your-domain
```

### Initializing Swarm Cluster

1. Initializing cluster on node-1

```sh
# Connecting in VM one
$ ssh root@ip-node-1
$ sudo tcs # type the password
$ sudo docker swarm init --advertise-addr 172.105.9.56

# Copy the output command to register nodes in cluster
# or generate a new comand token
$ sudo docker swarm join-token manager
$ sudo docker service ls
```

2. Adding node-2 in swarm

```sh
# Connecting in VM one
$ ssh root@ip-node-2
$ sudo tcs # type the password

# Paste here the join token copied by node-1
$ sudo docker swarm join --token <TOKEN>
```

2. Applying stack

```sh
$ mkdir swarm
$ sudo chown tcs:tcs swarm/

$ cd swarm
$ mkdir caddy_data
$ mkdir caddy_config
$ mkdir db-data
$ mkdir db-data/mongo
$ mkdir db-data/postgres

# paste here the docker-stack.production.yaml content
$ vim docker-stack.yaml

# docker stack deploy -c <YAML> <PREFIX>
$ sudo docker stack deploy -c docker-stack.yaml myapp
```

3. Show services

```sh
$ docker service ls
```

### Guides

- Gluster, share volumes

```
https://www.gluster.org/install/
```

- Mount volumes with ssh

```
https://phoenixnap.com/kb/sshfs
```
