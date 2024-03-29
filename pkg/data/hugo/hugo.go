package hugo

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// A Blog contains all the data of a Hugo blog. The Path represents the
// working directory for the site.
type Blog struct {
	Title, Path string
	Posts       []Post
}

// A Post contains all the metadata related to a hugo post, but not the content
// of the post itself
type Post struct {
	Title, Date, Path string
	Draft             bool
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() string {
	s := ""
	hugoConfigCmd := exec.Command("hugo", "config")
	hugoConfig, err := hugoConfigCmd.Output()
	s = string(hugoConfig)
	check(err)

	return s
}

// Load takes a path to a hugo site working directory and returns a Blog.
func Load(path string) Blog {
	os.Chdir(path)
	return Blog{Title: "blog", Path: path, Posts: loadPosts()}
}

// NewPost takes a title string and creates a new blog post before returning
// the path to the created file. Important: for now, the default is to make a
// new blog post in `$SITE_PATH/content/blog`
func (blog *Blog) NewPost(title string, extension string) string {
	log.Printf("Attempting to add post with title: %s, and extension %s\n", title, extension)
	os.Chdir(blog.Path)
	log.Printf("Path changed to: %s\n", blog.Path)

	filename := strings.ToLower(title)
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, ":", "")
	filePath := fmt.Sprintf("%s/%s.%s", "blog", filename, extension)
	log.Printf("Post file path is: %s\n", filePath)

	newPostCmd := exec.Command("hugo", "new", filePath)
	postPathRaw, err := newPostCmd.Output()
	postPath := string(postPathRaw)
	postPath = strings.Split(postPath, " ")[0]
	log.Printf("Post path from command result is: %s\n", postPath)
	check(err)

	blog.Posts = loadPosts()

	return blog.Posts[0].Path
}

func loadPosts() []Post {
	/// TODO find a way to partially reload the posts in a blog?
	/// Times when effective:
	/// 	a) deleting a post
	/// 	b) creating a post
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

// Synchronise removes all the files in a blog's public directory and then
// builds the site with Hugo.
func (blog Blog) Synchronise() {
	os.Chdir(blog.Path)

	buildCmd := exec.Command("hugo")

	err := buildCmd.Run()
	check(err)
}

// DeletePost removes the post file from the site's content/section directory.
// Warning: this method is destructive and should use user confirmation.
func (blog *Blog) DeletePost(deletePath string) error {

	os.Chdir(blog.Path)
	postPath := path.Join(blog.Path, deletePath)
	err := os.Remove(postPath)
	blog.Posts = loadPosts()
	return err
}
