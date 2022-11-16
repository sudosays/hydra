package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sudosays/hydra/internal/git"
	"github.com/sudosays/hydra/pkg/data/hugo"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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
	home, err := os.UserHomeDir()
	configPath := path.Join(home, ".config", "hydra.json")
	check(err)
	config, err = readConfig(configPath)
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
		fmt.Println("\nWhat would you like to do?\nCommands: [n]ew [e]dit [p]ublish [s]ync [q]uit")
		fmt.Print("> ")
		input := bufio.NewReader(os.Stdin)
		command, err := input.ReadString('\n')
		check(err)
		parseCommand(command)
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

func parseCommand(cmd string) {
	parts := strings.Split(cmd, " ")
	switch parts[0][0] {
	case 'e':
		fmt.Println("Trying to edit post #", parts[1])
	case 'q':
		os.Exit(0)
	}
}

func printPostList(blog hugo.Blog) {
	header, list := genPostList(blog)

	for _, col := range header {
		fmt.Print(col + "\t")
	}

	fmt.Println()

	for _, post := range list {
		for _, col := range post {
			fmt.Print(col + "\t")
		}

		fmt.Print("\n")

	}
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
