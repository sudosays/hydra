package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
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
	Extension string              `toml:"extension"`
	Editor    EditorCommand       `toml:"editor"`
	Sites     map[string]HugoSite `toml:"sites"`
}

// HugoSite contains the information for a hugo site listed in the config
type HugoSite struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
}

type EditorCommand struct {
	Command string `toml:"command"`
	Args    string `toml:"args"`
}

type FilterArgs struct {
	Sections  []string
	Draft     bool
	Published bool
	// DateBefore, DateAfter
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var config HydraConfig

var currentPageIndex int = 1 // For pagination purposes
var numPages int = 0

const maxItemsPerPage int = 10

func init() {

	userHomePath, err := os.UserHomeDir()
	defaultConfigFilePath := path.Join(userHomePath, ".config", "hydra.json")

	// Parse command line flags if any
	configFilePath := flag.String("config", defaultConfigFilePath, "Path to a config file")
	flag.Parse()

	config, err = readConfig(*configFilePath)
	check(err)

	fmt.Printf("[Config editor] Command: %v; Args: %v;\n", config.Editor.Command, config.Editor.Args)

}

func main() {

	clearTerm()
	// Setup to parse args
	blog := hugo.Load(config.Sites["main"].Path)

	// Calculate the number of pages
	numPages = int(math.Ceil(float64(len(blog.Posts)) / float64(maxItemsPerPage)))

	posts := blog.Posts
	// main REPL
	for {
		printPostList(posts)
		command := promptUser("\nWhat would you like to do?\n" +
			"Commands: [a]dd [e]dit [d]elete [n]ext/[p]rev page [q]uit\n> ")
		posts = parseCommand(command, blog)
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
	if _, err = toml.Decode(string(byteValue), &conf); err != nil {
		log.Fatal(err)
	}
	check(err)
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

func parseCommand(cmd string, blog hugo.Blog) []hugo.Post {
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
	case 'a':
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
	case 'n':
		if currentPageIndex < numPages {
			currentPageIndex++
		}
	case 'p':
		if currentPageIndex > 1 {
			currentPageIndex--
		}
	case 'f':
		// filtering
		// format/commands?
		// 'd' : toggle drafts-only
		// 'p' : toggle published-only
		// 's <section name>' : filter section by name (show only section), 's all' or 's' restore

	case 'q':
		clearTerm()
		fmt.Println("You have slain the hydra...")
		os.Exit(0)
	}
	return blog.Posts
}

func promptUser(prompt string) string {
	fmt.Print(prompt)
	input := bufio.NewReader(os.Stdin)
	ans, err := input.ReadString('\n')
	check(err)
	return ans
}

func printPostList(posts []hugo.Post) {
	// We might want to paginate the number of blog posts
	// Currently we will set the max to 10 posts per page

	fmtString := "%3s %-12s %-6s %-15s"

	fmt.Printf(fmtString, "#", "Date", "Draft", "Title")
	fmt.Print("\n")
	fmt.Println(strings.Repeat("-", 50))

	startPostIndex := (currentPageIndex - 1) * maxItemsPerPage
	endPostIndex := startPostIndex + maxItemsPerPage // Possibly longer than len(list)
	if endPostIndex > len(posts) {
		endPostIndex = len(posts)
	}

	pageList := posts[startPostIndex:endPostIndex]

	for i, post := range pageList {
		draftStatus := ""
		if post.Draft {
			draftStatus = "D"
		}
		datetime, _ := time.Parse(time.RFC3339, post.Date)
		date := datetime.Format("2006/01/02")

		title := post.Title
		if len(title) > 30 {
			title = title[:26]
			title += "..."
		}

		fmt.Printf(fmtString, strconv.Itoa(startPostIndex+i+1), date, draftStatus, title)
		fmt.Print("\n")

	}

	fmt.Printf("( Showing [%d-%d]", startPostIndex+1, endPostIndex)
	fmt.Printf(" | Page %d of %d )\n", currentPageIndex, numPages)
}

func clearTerm() {
	fmt.Print("\033[H\033[2J")
}
