package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev"
	}
	conf = confs[env]
	if conf == nil {
		log.Fatal("无效的环境变量")
	}

	http.HandleFunc("/", router)
	log.Fatal(http.ListenAndServe(":8000", nil))
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
	path := req.URL.Path
	status := http.StatusOK

	// logger
	defer func(start time.Time) {
		cost := time.Since(start)
		cost = cost.Round(time.Millisecond)
		log.Printf("%v %v %v %v\n", req.Method, req.URL, status, cost)
	}(time.Now())

	// 首页
	if path == "/" {
		content, err := getIndexHtml()
		if err == nil {
			w.Write(content)
			return
		}
		log.Println(err)
	}

	// 文章详情页
	if match := postRegexp.MatchString(path); match {
		name := postRegexp.FindStringSubmatch(path)[1]
		filename := name + ".md"
		fPath := filepath.Join(conf.blogMdPath, filename)
		content, err := getPostHtml(fPath)
		if err == nil {
			w.Write(content)
			return
		}
		log.Println(err)
	}

	// search
	if path == "/search" {
		if q := req.URL.Query().Get("q"); q != "" {
			content, err := getSearchHtml(q)
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
