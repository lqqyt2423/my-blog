FROM golang:1.11-alpine

WORKDIR /go/src/blog

COPY . .

RUN apk add --no-cache git mercurial \
  && go get -d -v gopkg.in/russross/blackfriday.v2 \
  && apk del git mercurial

RUN go install -v ./

ENV GO_ENV=prod

CMD ["blog"]

EXPOSE 8000
