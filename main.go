package main

import (
	"html/template"
	"lqqyt2423/go_blog/article"
	"net/http"
	"regexp"
	"time"
)

var source article.Article

func main() {
	loadConfFromEnv()

	source = &article.Cache{}
	source.Init(conf.blogMdPath)

	http.HandleFunc("/", router)
	logger.Fatal(http.ListenAndServe(conf.listenAt, nil))
}

// 文章详情页路由正则
var postRegexp = regexp.MustCompile(`^/post/(2\d{7}(\w|-)+)\.html$`)

// 文章文件名正则
var postFileRegexp = regexp.MustCompile(`^2\d{7}(\w|-)+\.md$`)

// md 替换为 html
var mdToHTMLRegexp = regexp.MustCompile(`\.md$`)

// 文章标题正则
var postTitleRegexp = regexp.MustCompile(`^\s*?#\s*`)

// 文章日期正则
var postDateRegexp = regexp.MustCompile(`2\d{7}`)

var notFoundTemplate = template.Must(template.ParseFiles("template/base.tmpl", "template/404.tmpl"))

func router(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	status := http.StatusOK

	// logger
	defer func(start time.Time) {
		cost := time.Since(start)
		cost = cost.Round(time.Millisecond)
		logger.Printf("%v %v %v %v\n", method, req.URL, status, cost)
	}(time.Now())

	// 首页
	if method == "GET" && path == "/" {
		content, err := source.GetAll()
		if err == nil {
			w.Write(content)
			return
		}
		logger.Println(err)
	}

	// 文章详情页
	if method == "GET" {
		if match := postRegexp.MatchString(path); match {
			name := postRegexp.FindStringSubmatch(path)[1]
			content, err := source.Get(name)
			if err == nil {
				w.Write(content)
				return
			}
			logger.Println(err)
		}
	}

	// search
	if method == "GET" && path == "/search" {
		if q := req.URL.Query().Get("q"); q != "" {
			content, err := source.Search(q)
			if err == nil {
				w.Write(content)
				return
			}
		}
	}

	// 404
	status = http.StatusNotFound
	w.WriteHeader(status)
	notFoundTemplate.Execute(w, &struct{ Title string }{Title: "404"})
}

func timeFormat(t time.Time) string {
	return t.Format("2006年01月02日")
}

func getTimeFromPath(path string) time.Time {
	dateStr := postDateRegexp.FindString(path)
	date, _ := time.Parse("20060102", dateStr)
	return date
}
