package article

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"gopkg.in/russross/blackfriday.v2"
)

// ==============================
// Detail
// ==============================

// 文章日期正则
var postDateRegexp = regexp.MustCompile(`2\d{7}`)

// 文章标题正则
var postTitleRegexp = regexp.MustCompile(`^\s*?#\s*`)

func timeFormat(t time.Time) string {
	return t.Format("2006年01月02日")
}

var postTemplate = template.Must(template.New("base.tmpl").
	Funcs(template.FuncMap{"format": timeFormat}).
	ParseFiles("template/base.tmpl", "template/post.tmpl"))

// Post struct
type Post struct {
	Path    string
	Title   string
	Content template.HTML
	Date    time.Time
	score   int
}

// PostWeb I don't know
type PostWeb struct {
	*Post
}

func getTimeFromPath(path string) time.Time {
	dateStr := postDateRegexp.FindString(path)
	date, _ := time.Parse("20060102", dateStr)
	return date
}

func getPostHTML(path string) ([]byte, error) {
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

// ==============================
// List
// ==============================

// 文章文件名正则
var postFileRegexp = regexp.MustCompile(`^2\d{7}(\w|-)+\.md$`)

// md 替换为 html
var mdToHTMLRegexp = regexp.MustCompile(`\.md$`)

// PostList type
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

// Index struct
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

func getPost(notePath string, filename string, c chan<- postRes) {
	path := filepath.Join(notePath, filename)
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

func getIndexHTML(notePath string) ([]byte, error) {
	files, err := ioutil.ReadDir(notePath)
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
		go getPost(notePath, file, c)
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

// ==============================
// Search
// ==============================

var searchTemplate = template.Must(template.New("base.tmpl").
	Funcs(template.FuncMap{"format": timeFormat}).
	ParseFiles("template/base.tmpl", "template/search.tmpl"))

// Search struct
type Search struct {
	Title     string
	QueryWord string
	Posts     []*Post
}

// SearchPosts type
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

func getSearchPost(notePath string, filename string, c chan<- postRes, q string) {
	path := filepath.Join(notePath, filename)
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

func getSearchHTML(notePath string, q string) ([]byte, error) {
	files, err := ioutil.ReadDir(notePath)
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
		go getSearchPost(notePath, file, c, q)
	}

	for range targetFiles {
		if r := <-c; r.err == nil {
			posts = append(posts, r.post)
		}
	}

	sort.Sort(posts)

	buf := bytes.NewBuffer(nil)
	err = searchTemplate.Execute(buf, &Search{Title: "搜索页", QueryWord: q, Posts: posts})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ==============================
// Entry
// ==============================

// Fs struct
type Fs struct {
	NotePath string
}

// Init fs config
func (fs *Fs) Init(notePath interface{}) {
	fs.NotePath = notePath.(string)
}

// Get article detail html
func (fs *Fs) Get(name string) ([]byte, error) {
	filename := name + ".md"
	filePath := filepath.Join(fs.NotePath, filename)

	return getPostHTML(filePath)
}

// GetAll function
func (fs *Fs) GetAll() ([]byte, error) {
	return getIndexHTML(fs.NotePath)
}

// Search function
func (fs *Fs) Search(q string) ([]byte, error) {
	return getSearchHTML(fs.NotePath, q)
}
