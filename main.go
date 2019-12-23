package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	loadConfFromEnv()

	// redis
	c, err := redis.Dial("tcp", conf.redisAddress)
	if err != nil {
		log.Fatal(err)
	}
	redisConn = c
	defer c.Close()

	http.HandleFunc("/", router)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

var mu sync.Mutex

// redis connection
var redisConn redis.Conn

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

// 根据 path 递增浏览量并返回之后的浏览量
// 如果 path 传空则只递增全站浏览量
func incView(path string) (site, page int64) {
	key := ""
	if path != "" {
		key = "blog:" + path
	}
	mu.Lock()
	redisConn.Send("INCR", "blog:site-view")
	if key != "" {
		redisConn.Send("INCR", key)
	}
	res, err := redisConn.Do("")
	mu.Unlock()
	if err != nil {
		log.Println("incView error:", err)
		return
	}
	site = res.([]interface{})[0].(int64)
	if key != "" {
		page = res.([]interface{})[1].(int64)
	}
	return
}
