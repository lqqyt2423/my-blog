FROM golang:1.11
RUN go get -d -v gopkg.in/russross/blackfriday.v2
RUN go get -d -v github.com/gomodule/redigo/redis
WORKDIR /go/src/blog
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY . .
COPY --from=0 /go/src/blog/app .
CMD ["./app"]
EXPOSE 8000
