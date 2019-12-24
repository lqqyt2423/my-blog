package post

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// PostFileRegexp 文章文件名正则
var PostFileRegexp = regexp.MustCompile(`^2\d{7}(\w|-)+\.md$`)

// MdToHTMLRegexp md 替换为 html
var MdToHTMLRegexp = regexp.MustCompile(`\.md$`)

// PostTitleRegexp 文章标题正则
var PostTitleRegexp = regexp.MustCompile(`^\s*?#\s*`)

// PostDateRegexp 文章日期正则
var PostDateRegexp = regexp.MustCompile(`2\d{7}`)

// Post struct
type Post struct {
	Path    string
	Title   string
	Date    time.Time
	Content string
	HTML    string
}

type postchan struct {
	post *Post
	err  error
}

func loadPost(dirname string, fileinfo os.FileInfo, pc chan<- postchan) {
	filename := fileinfo.Name()
	path := filepath.Join(dirname, filename)
	body, err := ioutil.ReadFile(path)
	if err != nil {
		pc <- postchan{nil, err}
		return
	}

	content := string(body)
	title := PostTitleRegexp.FindString(content)
	post := &Post{
		Path:    MdToHTMLRegexp.ReplaceAllString(filename, ".html"),
		Title:   title,
		Content: content,
	}
	post.ParseDate()
	pc <- postchan{post, nil}
}

// InitPosts 初始化文章列表
func InitPosts(dirname string) ([]*Post, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var targetFiles []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if matched := PostFileRegexp.MatchString(file.Name()); matched {
			targetFiles = append(targetFiles, file)
		}
	}
	if len(targetFiles) == 0 {
		return nil, fmt.Errorf("no post files")
	}

	var posts []*Post
	return posts, nil
}

// ParseDate add Date field to Post
func (p *Post) ParseDate() {
	_, filename := filepath.Split(p.Path)
	dateStr := PostDateRegexp.FindString(filename)
	date, _ := time.Parse("20060102", dateStr)
	p.Date = date
}
