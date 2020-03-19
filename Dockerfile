FROM golang:1.11
RUN go get -d -v gopkg.in/russross/blackfriday.v2
WORKDIR /go/src/lqqyt2423/go_blog
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM alpine:latest
WORKDIR /root
COPY . .
COPY --from=0 /go/src/lqqyt2423/go_blog/go_blog .
CMD ["./go_blog"]
EXPOSE 7000
