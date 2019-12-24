package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"time"

	"gopkg.in/russross/blackfriday.v2"
)

var postTemplate = template.Must(template.New("base.tmpl").
	Funcs(template.FuncMap{"format": timeFormat}).
	ParseFiles("template/base.tmpl", "template/post.tmpl"))

type Post struct {
	Path    string
	Title   string
	Content template.HTML
	Date    time.Time
	score   int
}

type PostWeb struct {
	*Post
}

func getPostHtml(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	post := handleContent(data)
	post.Date = getTimeFromPath(path)

	buf := bytes.NewBuffer(nil)
	err = postTemplate.Execute(buf, &PostWeb{post})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func handleContent(data []byte) *Post {
	title, content := getTitleAndContent(data)
	content = blackfriday.Run(content)
	return &Post{Title: string(title), Content: template.HTML(content)}
}

func getTitleAndContent(data []byte) ([]byte, []byte) {
	data = bytes.TrimSpace(data)
	arr := bytes.Split(data, []byte("\n"))
	title := arr[0]
	title = postTitleRegexp.ReplaceAll(title, []byte(""))
	content := bytes.Join(arr[1:], []byte("\n"))
	return title, content
}
