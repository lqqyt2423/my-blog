package main

type config struct {
	blogMdPath string
}

var confs = map[string]*config{
	"dev": &config{
		blogMdPath: "/Users/liqiang/Documents/code/programming_note",
	},
	"prod": &config{
		blogMdPath: "/root/programming_note",
	},
}

var conf *config
