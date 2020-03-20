# README

## go mod

```bash
go mod init lqqyt2423/go_blog
go list -m all
export GOPROXY=https://goproxy.io
go get gopkg.in/russross/blackfriday.v2

go mod tidy

go help mod
```

go mod 在本地的缓存文件夹为 GOPATH/pkg/mod

## 本地测试

```bash
go run *.go

make
./go_blog
```

## docker

```bash
docker build -t go_blog .
docker run --rm go_blog
```

## docker-compose

```bash
docker-compose build
docker-compose up
docker-compose down
```
