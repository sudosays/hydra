package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sudosays/hydra/internal/git"
	"github.com/sudosays/hydra/pkg/data/hugo"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type HydraConfig struct {
	Extension string        `json:"extension"`
	Editor    EditorCommand `json:"editor"`
	Sites     []HugoSite    `json:"sites"`
}

// HugoSite contains the information for a hugo site listed in the config
type HugoSite struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type EditorCommand struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var config HydraConfig

func init() {

	userHomePath, err := os.UserHomeDir()
	defaultConfigFilePath := path.Join(userHomePath, ".config", "hydra.json")

	// Parse command line flags if any
	configFilePath := flag.String("config", defaultConfigFilePath, "Path to a config file")

	config, err = readConfig(*configFilePath)
	check(err)
}

func main() {

	clearTerm()
	// Setup to parse args
	blog := hugo.Load(config.Sites[0].Path)

	if git.IsRepo() {
		go git.Update()
	}

	// main REPL
	for {
		printPostList(blog)
		command := promptUser("\nWhat would you like to do?\n" +
			"Commands: [n]ew [e]dit [d]elete [p]ublish [s]ync [q]uit\n> ")
		blog = parseCommand(command, blog)
		clearTerm()
	}
}

func readConfig(path string) (HydraConfig, error) {
	conf := HydraConfig{}
	configFile, err := os.Open(path)
	check(err)
	byteValue, err := ioutil.ReadAll(configFile)
	check(err)
	//fmt.Printf("%s", byteValue)
	err = json.Unmarshal(byteValue, &conf)
	check(err)
	//fmt.Printf("Config is: %+v\n", conf)
	configFile.Close()
	return conf, nil
}

func startEditor(path string) {
	editorCmd := exec.Command(config.Editor.Command, config.Editor.Args, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	err := editorCmd.Run()
	check(err)
}

func genPostList(blog hugo.Blog) ([]string, [][]string) {

	posts := blog.Posts
	headings := []string{"#", "Date", "Draft", "Title"}
	var postList [][]string

	for i, post := range posts {
		draftStatus := "False"
		if post.Draft {
			draftStatus = "True"
		}
		datetime, _ := time.Parse(time.RFC3339, post.Date)
		date := datetime.Format("2006/01/02")
		postList = append(
			postList,
			[]string{fmt.Sprintf("%d", i+1), date, draftStatus, post.Title})
	}

	return headings, postList
}

func parseCommand(cmd string, blog hugo.Blog) hugo.Blog {
	parts := strings.Split(cmd, " ")
	switch parts[0][0] {
	case 'e':
		index := -1
		if len(parts) > 1 {
			// startEditor
			i, err := strconv.Atoi(strings.Trim(parts[1], "\n\r "))
			check(err)
			index = i

		} else {
			i, err := strconv.Atoi(strings.Trim(promptUser("Enter a post number to edit:\n> "), "\n\r "))
			check(err)
			index = i
		}
		if index > 0 && index < len(blog.Posts) {
			startEditor(blog.Posts[index-1].Path)
		}
	case 'n':
		title := ""
		if len(parts) > 1 {
			// Call creation of new post with title given
			title = strings.Join(parts[1:], " ")
		} else {
			title = promptUser("Please enter a title for the post:\n> ")
		}
		fmt.Println("Attempting to create post with title: %s", title)
		postPath := blog.NewPost(title, config.Extension)
		startEditor(postPath)
	case 'd':
		// Ask for confirmation first
		// Possibly: add a hugo 'trash' folder in future for soft deletion.
		i, err := strconv.Atoi(strings.Trim(parts[1], "\n\r "))
		check(err)
		post := blog.Posts[i-1]
		warn := "WARNING: You are about to delete the post titled '%s'.\nThis action is irreversible, especially if the post has not been synchronised with git.\nProceed? [y/N] Default: N.\n> "
		ans := promptUser(fmt.Sprintf(warn, post.Title))
		if ans[0] == 'y' {
			blog.DeletePost(post.Path)
		}
	case 's':
		// Synchronise
		// Options:
		// 1) synchronise the repo only with the source code
		// 2) run "hugo" and publish the website code/repo found in ./public
		//
		log.Println("Attempting to sync repository")
		git.Sync()
	case 'q':
		clearTerm()
		fmt.Println("You have slain the hydra...")
		os.Exit(0)
	}

	return blog
}

func promptUser(prompt string) string {
	fmt.Print(prompt)
	input := bufio.NewReader(os.Stdin)
	ans, err := input.ReadString('\n')
	check(err)
	return ans
}

func printPostList(blog hugo.Blog) {
	// We might want to paginate the number of blog posts
	// Currently we will set the max to 10 posts per page

	maxItemsPerPage := 10.0

	header, list := genPostList(blog)

	numPages := int(math.Ceil(float64(len(list)) / float64(maxItemsPerPage)))

	for _, col := range header {
		fmt.Print(col + "\t")
	}

	fmt.Println()

	for _, post := range list[:10] {
		for _, col := range post {
			fmt.Print(col + "\t")
		}

		fmt.Print("\n")

	}

	fmt.Printf("Page 1 of %d\n", numPages)
}

func clearTerm() {
	fmt.Print("\033[H\033[2J")
}

// func printLogo() {

// 	logoText := `
//  /$$                       /$$
// | $$                      | $$
// | $$$$$$$  /$$   /$$  /$$$$$$$  /$$$$$$  /$$$$$$
// | $$__  $$| $$  | $$ /$$__  $$ /$$__  $$|____  $$
// | $$  \ $$| $$  | $$| $$  | $$| $$  \__/ /$$$$$$$
// | $$  | $$| $$  | $$| $$  | $$| $$      /$$__  $$
// | $$  | $$|  $$$$$$$|  $$$$$$$| $$     |  $$$$$$$
// |__/  |__/ \____  $$ \_______/|__/      \_______/
//            /$$  | $$
//           |  $$$$$$/
//            \______/
// `
// }
