FROM golang:1.13
WORKDIR /root
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM alpine:latest
WORKDIR /root
COPY template ./template
COPY --from=0 /root/go_blog .
CMD ["./go_blog"]
EXPOSE 7000
