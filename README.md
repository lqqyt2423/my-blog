# README

## build

```
docker image build -t my-blog .
```

## run

#### dev

```
docker container run \
  -d \
  --rm \
  --name my-blog \
  -p 7000:8000 \
  --volume /Users/liqiang/Documents/code/programming_note:/root/programming_note \
  my-blog
```

#### prod test

```
docker container run \
  -d \
  --rm \
  --name my-blog \
  -p 7000:8000 \
  --volume /root/programming_note:/root/programming_note \
  my-blog
```

#### prod

```
docker container run \
  -d \
  --name my-blog \
  -p 7000:8000 \
  --volume /root/programming_note:/root/programming_note \
  --restart=always \
  my-blog
```
