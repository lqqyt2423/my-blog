package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
)

var searchTemplate = template.Must(template.New("base.tmpl").
	Funcs(template.FuncMap{"format": timeFormat}).
	ParseFiles("template/base.tmpl", "template/search.tmpl"))

type Search struct {
	Title     string
	QueryWord string
	Posts     []*Post
	SiteView  int64
}

type SearchPosts []*Post

// 文章排序
func (l SearchPosts) Len() int { return len(l) }
func (l SearchPosts) Less(i, j int) bool {
	if l[i].score != l[j].score {
		return l[i].score > l[j].score
	}
	if !l[i].Date.Equal(l[j].Date) {
		return l[i].Date.After(l[j].Date)
	}
	return l[i].Title < l[j].Title
}
func (l SearchPosts) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func getSearchHtml(q string) ([]byte, error) {
	files, err := ioutil.ReadDir(conf.blogMdPath)
	if err != nil {
		return nil, err
	}

	var posts SearchPosts
	var targetFiles []string

	for _, file := range files {
		filename := file.Name()
		if matched := postFileRegexp.MatchString(filename); matched {
			targetFiles = append(targetFiles, filename)
		}
	}

	c := make(chan postRes)

	for _, file := range targetFiles {
		go getSearchPost(file, c, q)
	}

	for range targetFiles {
		if r := <-c; r.err == nil {
			posts = append(posts, r.post)
		}
	}

	sort.Sort(posts)

	buf := bytes.NewBuffer(nil)
	sv, _ := incView("")
	err = searchTemplate.Execute(buf, &Search{Title: "搜索页", QueryWord: q, Posts: posts, SiteView: sv})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getSearchPost(filename string, c chan<- postRes, q string) {
	path := filepath.Join(conf.blogMdPath, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		c <- postRes{err: err}
		return
	}
	bq := bytes.ToLower([]byte(q))
	title, content := getTitleAndContent(data)
	contentCount := bytes.Count(bytes.ToLower(content), bq)
	score := len(bq) * contentCount * 10000 / len(content)
	if bytes.Contains(bytes.ToLower(title), bq) {
		score += 1000
	}
	if score == 0 {
		c <- postRes{err: errors.New("not match")}
		return
	}
	post := &Post{
		Path:  mdToHTMLRegexp.ReplaceAllString(filename, ".html"),
		Title: string(title),
		Date:  getTimeFromPath(filename),
		score: score,
	}
	c <- postRes{post: post}
}
