package main

type config struct {
	blogMdPath   string
	redisAddress string
}

var confs = map[string]*config{
	"dev": &config{
		blogMdPath:   "/Users/liqiang/Documents/code/programming_note",
		redisAddress: ":6379",
	},
	"prod": &config{
		blogMdPath:   "/root/programming_note",
		redisAddress: ":6379",
	},
}

var conf *config
