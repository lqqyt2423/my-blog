package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type PostList []*Post

// 文章排序
func (l PostList) Len() int { return len(l) }
func (l PostList) Less(i, j int) bool {
	if l[i].Date.Equal(l[j].Date) {
		return l[i].Title < l[j].Title
	}
	return l[i].Date.After(l[j].Date)
}
func (l PostList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

type Index struct {
	Title string
	Posts PostList
}

type postRes struct {
	post *Post
	err  error
}

var postsTemplate = template.Must(template.New("base.tmpl").
	Funcs(template.FuncMap{"format": timeFormat}).
	ParseFiles("template/base.tmpl", "template/index.tmpl"))

func getIndexHtml() ([]byte, error) {
	files, err := ioutil.ReadDir(conf.blogMdPath)
	if err != nil {
		return nil, err
	}

	var posts PostList
	var targetFiles []string

	for _, file := range files {
		filename := file.Name()
		if matched := postFileRegexp.MatchString(filename); matched {
			targetFiles = append(targetFiles, filename)
		}
	}

	c := make(chan postRes)

	// 生产者
	for _, file := range targetFiles {
		go getPost(file, c)
	}

	// 消费者
	for range targetFiles {
		if r := <-c; r.err == nil {
			posts = append(posts, r.post)
		}
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no posts")
	}

	sort.Sort(posts)

	buf := bytes.NewBuffer(nil)
	err = postsTemplate.Execute(buf, &Index{Title: "首页", Posts: posts})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getPost(filename string, c chan<- postRes) {
	path := filepath.Join(conf.blogMdPath, filename)
	file, err := os.Open(path)
	if err != nil {
		c <- postRes{post: nil, err: err}
		return
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		c <- postRes{post: nil, err: err}
		return
	}

	title := postTitleRegexp.ReplaceAll(line, []byte(""))
	post := &Post{
		Path:  mdToHTMLRegexp.ReplaceAllString(filename, ".html"),
		Title: string(title),
		Date:  getTimeFromPath(filename),
	}

	c <- postRes{post: post, err: nil}
}
