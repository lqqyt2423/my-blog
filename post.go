package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"time"

	"gopkg.in/russross/blackfriday.v2"
)

var postTemplate = template.Must(template.New("post.gtpl").
	Funcs(template.FuncMap{"format": timeFormat}).
	ParseFiles("post.gtpl"))

type Post struct {
	Path    string
	Title   string
	Content template.HTML
	Date    time.Time
}

func getPostHtml(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	post := handleContent(data)
	post.Date = getTimeFromPath(path)

	buf := bytes.NewBuffer(nil)
	err = postTemplate.Execute(buf, post)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func handleContent(data []byte) *Post {
	data = bytes.TrimSpace(data)
	arr := bytes.Split(data, []byte("\n"))
	title := arr[0]
	title = postTitleRegexp.ReplaceAll(title, []byte(""))
	content := bytes.Join(arr[1:], []byte("\n"))
	content = blackfriday.Run(content)
	return &Post{Title: string(title), Content: template.HTML(content)}
}
