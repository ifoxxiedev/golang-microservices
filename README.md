### Dependencies

> Check if your environment have this dependecies (infra)
> \
> Note: `docker-compose.yaml` contains this dependecies defined.

<br />

## GoLang CLI Commands

- Initializing a module
```sh
go env # Mostrar variaveis de ambiente (de configuração)

go mod init <micro_service>
```

- Buiding Application 
```sh
go build
```

- Download for [package](https://pkg.go.dev/)
```sh

go get -u github.com/badoux/checkmail
```

- Install, Remove and Update Packages
```sh
go mod tidy # Install dependecies or Remove unused Packages (it's needs go.sum/go.mod)
go install # Install dependencies (it's needs go.sum/go.mod)
```

- Execute a builded package
```sh
go run main.go
```

- Go Envs
```sh
go env
go env GOTPATH
```


- Getting help with CLI
```sh
go help
go help get
```

## How To Use

To clone and run this application, you'll need [Git](https://git-scm.com) and [GoLang](https://go.dev/)

- Clone Repository

```bash
$ git clone https://github.com/ifoxxiedev/golang-microservices.git golang-microservices
$ cd golang-microservices
```

### Development Way

- Install Depths
```sh
$ go download # OR go install
```

- Configuring Live Reload
```
# If don't have local live reload
$ go install https://github.com/air-verse/air

$ air init
```

> With config base
```
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main.exe"
  cmd = "go build -o ./tmp/main.exe ./cmd/api"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[proxy]
  app_port = 0
  enabled = false
  proxy_port = 0

[screen]
  clear_on_rebuild = false
  keep_scroll = true

```


#### Running applications

- Preparing environment variables

```sh
$ cp .env.example .env
```

##### Running Local
```
$ air -d
$ 
```

##### Running With Docker (Build Wait)

```bash
$ cd project
$ docker-compose -f docker-compose.yml up -d
```

- Show logs

```bash
$ cd project
$ docker-compose -f docker-compose.yml logs -f <micro_service>
```

## Utilities

This application uses the following packages and patterns:

- [GoLang](https://go.dev/)
- [GoPackages](https://pkg.go.dev/)
- [Chi](https://go-chi.io/#/)
- [Design Patterns](https://refactoring.guru/design-patterns)

---

<br />
<br />

## License

ISC