# Deploy - Kubernetes Guide

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

## Deploying in kubnernetes cluster (LOCAL)

### Spinning k8s cluster

1. Creating a cluster (local) with minikube

```sh
# Creating cluster with minikube
$ minikube start --nodes=2

# Show cluster status
$ minuke status

# Stop cluster
$ minikube stop

# Start cluster
$ minikube start

# Enable minikube dashboard
$ minikube dashboard
```

1. Creating a cluster (local) with k3d

```sh
# Creating cluster with k3d
$ k3d cluster create myapps --servers 1 -p "30000:30000@loadbalancer"
```
