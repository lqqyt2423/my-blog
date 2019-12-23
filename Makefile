blog: main.go logger.go config.go index.go post.go search.go
	go build -o blog main.go logger.go config.go index.go post.go search.go

clean:
	rm blog
