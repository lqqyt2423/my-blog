# README

## go mod

```bash
go mod init lqqyt2423/go_blog
go list -m all
export GOPROXY=https://goproxy.io
go get github.com/gomodule/redigo/redis

go mod tidy

go help mod
```

go mod 在本地的缓存文件夹为 GOPATH/pkg/mod

## 本地测试

```bash
go run *.go

make
./blog
```

## docker-compose

```bash
docker-compose build
docker-compose up
docker-compose down
docker-compose -f docker-compose-prod.yml up
docker-compose -f docker-compose-prod.yml up -d
```

## build

```bash
docker image build -t my-blog .
```

## run

### dev

```bash
docker container run \
  -d \
  --rm \
  --name my-blog \
  -p 7000:8000 \
  --volume /Users/liqiang/Documents/code/programming_note:/root/programming_note \
  my-blog
```

### prod test

```bash
docker container run \
  -e GO_ENV=prod \
  -d \
  --rm \
  --name my-blog \
  -p 7000:8000 \
  --volume /root/programming_note:/root/programming_note \
  my-blog
```

### prod

```bash
docker container run \
  -e GO_ENV=prod \
  -d \
  --name my-blog \
  -p 7000:8000 \
  --volume /root/programming_note:/root/programming_note \
  --restart=always \
  my-blog
```
