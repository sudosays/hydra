package main

import (
	"encoding/json"
	"fmt"
	libui "github.com/sudosays/pneuma/internal/ui"
	"github.com/sudosays/pneuma/pkg/data/hugo"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
)

type Site struct {
	Name, Path string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printPosts(posts []hugo.Post) {
	for _, p := range posts {
		indicator := ""
		if p.Draft {
			indicator = "*"
		}
		fmt.Printf("%s%s\n", indicator, p.Title)
	}
}

func readConfig(path string) []Site {
	var sites []Site
	configFile, err := os.Open(path)
	check(err)
	decoder := json.NewDecoder(configFile)
	for {
		var site Site
		if err := decoder.Decode(&site); err == io.EOF {
			break
		} else {
			check(err)
			sites = append(sites, site)
		}
	}
	configFile.Close()
	return sites
}

func startEditor(path string) {
	editorCmd := exec.Command("vim", path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	err := editorCmd.Run()
	check(err)
}

func main() {
	home, err := os.UserHomeDir()
	configPath := path.Join(home, ".config", "pneuma.json")
	check(err)
	os.Chdir(home)
	cwd, err := os.Getwd()
	check(err)

	sites := readConfig(configPath)

	ui := libui.Init()
	ui.MoveCursor(0, 0)
	ui.PutStr(fmt.Sprintf("Current working dir: %s", cwd))
	ui.MoveCursor(0, 1)
	ui.PutStr(fmt.Sprintf("Config path: %s", configPath))
	ui.MoveCursor(0, 2)
	ui.PutStr(fmt.Sprintf("Sites from config:"))
	var sitesList []string
	for _, site := range sites {
		sitesList = append(sitesList, fmt.Sprintf("%s (%s)", site.Name, site.Path))
	}
	ui.MoveCursor(0, 3)
	ui.PrintList(sitesList)
	ui.MoveCursor(0, 3+len(sitesList))
	prompt := "Please choose a site (default=1): "
	ui.PutStr(prompt)
	ui.MoveCursor(len(prompt), 3+len(sitesList))

	siteSelect := getSelection(ui, len(sitesList))

	ui.Reset()
	site := hugo.Load(sites[siteSelect].Path)
	posts := site.Posts
	var postList []string
	for _, post := range posts {
		indicator := ""
		if post.Draft {
			indicator = "*"
		}
		postList = append(postList, fmt.Sprintf("%s%s", indicator, post.Title))
	}

	ui.PutStr("Posts from the blog:")
	ui.MoveCursor(0, 1)
	ui.PrintList(postList)

	ui.MoveCursor(0, 2+len(postList))
	prompt = "Select a post to edit: "
	ui.PutStr(prompt)
	ui.MoveCursor(len(prompt), 2+len(postList))

	postNum := getSelection(ui, len(postList))
	startEditor(path.Join(sites[siteSelect].Path, posts[postNum].Path))

	for {
		ui.Tick()
	}

}

func getSelection(ui libui.PneumaUI, max int) int {
	selection := 0
	input := ui.WaitForInput()
	choice, err := strconv.Atoi(input)
	if err == nil {
		choice--
		if choice > 0 && choice <= max {
			selection = choice
		} else {
			check(err)
		}
	}
	return selection
}
