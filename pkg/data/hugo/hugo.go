package hugo

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Blog struct {
	Title, Path string
	Posts       []Post
}

type Post struct {
	Title, Date, Path string
	Draft             bool
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Config() string {
	s := ""
	hugoConfigCmd := exec.Command("hugo", "config")
	hugoConfig, err := hugoConfigCmd.Output()
	s = string(hugoConfig)
	check(err)

	return s
}

func Load(path string) Blog {
	os.Chdir(path)
	return Blog{Title: "blog", Path: path, Posts: loadPosts()}
}

func (blog *Blog) NewPost(title string) string {
	os.Chdir(blog.Path)

	ext := ".md"
	filename := strings.ToLower(title)
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, ":", "")
	filePath := fmt.Sprintf("%s/%s%s", "blog", filename, ext)

	newPostCmd := exec.Command("hugo", "new", filePath)
	_, err := newPostCmd.Output()
	check(err)

	blog.Posts = loadPosts()

	return filePath
}

func loadPosts() []Post {
	var posts []Post
	hugoListCmd := exec.Command("hugo", "list", "all")
	rawPostList, err := hugoListCmd.Output()
	csvReader := csv.NewReader(strings.NewReader(string(rawPostList)))
	isHeader := true
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else {
			if isHeader {
				isHeader = !isHeader
				continue
			} else {
				post := Post{Path: record[0],
					Date:  record[3],
					Title: record[2],
					Draft: (record[6] == "true"),
				}
				posts = append(posts, post)
			}
		}
	}
	check(err)
	return posts
}
