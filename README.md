# README

```
docker image build -t my-blog .
```

```
docker container run \
  -d \
  --rm \
  --name my-blog \
  -p 7000:8000 \
  --volume /Users/liqiang/Documents/code/programming_note:/root/programming_note \
  my-blog
```
