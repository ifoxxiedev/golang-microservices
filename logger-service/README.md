# Logger Service

## Working with gRrc + GoLang

### Install tools to work with GrpC

- Installing binaries

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

- Installing packages for service

```
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

- Download proto compiler, [releases](https://github.com/protocolbuffers/protobuf/releases)

```sh
$ curl -fsSL https://github.com/protocolbuffers/protobuf/releases/download/v28.0/protoc-28.0-linux-x86_64.zip -o protoc-28.0-linux-x86_64.zip
$ tar -xvf protoc-28.0-linux-x86_64.zip

$ cd protoc-28.0-linux-x86_64/bin
$ mv protoc ~/go/bin


$ whitch protoc
$ protoc
```

- Compile protofile

```
export GOPATH=~/go
export PATH=$PATH:/$GOPATH/bin
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
```
