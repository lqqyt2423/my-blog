package main

import (
	"os"
)

type config struct {
	blogMdPath string
	listenAt   string
}

var confs = map[string]*config{
	"dev": &config{
		blogMdPath: "/Users/liqiang/Documents/_personal/code/programming_note",
		listenAt:   ":8000",
	},
	"prod": &config{
		blogMdPath: "/root/programming_note",
		listenAt:   ":7000",
	},
}

var conf *config

func loadConfFromEnv() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev"
	}
	conf = confs[env]
	if conf == nil {
		logger.Fatalln("GO_ENV is one of: dev, prod")
	}
	logger.Printf("env %s config loaded\n", env)
	conf.show()
}

func (c *config) show() {
	logger.Printf("[config] %s: %s\n", "blogMdPath", c.blogMdPath)
	logger.Printf("[config] %s: %s\n", "listenAt", c.listenAt)
}
