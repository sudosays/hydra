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

var sites []Site
var ui libui.PneumaUI

func init() {
	home, err := os.UserHomeDir()
	configPath := path.Join(home, ".config", "pneuma.json")
	check(err)
	sites = readConfig(configPath)
	ui = libui.Init()

}

func main() {

	var site hugo.Blog
	if len(sites) > 1 {
		site = siteSelect()
	} else {
		site = hugo.Load(sites[0].Path)
	}

	siteOverview(site)

	for {
		ui.Tick()
	}

}

func siteSelect() hugo.Blog {
	ui.Reset()
	ui.AddLabel(0, 0, "Sites from config:")

	var sitesList [][]string
	sitesList = append(sitesList, []string{"Choice", "Name", "Path"})
	for i, site := range sites {
		sitesList = append(sitesList, []string{fmt.Sprintf("%d", i+1), site.Name, site.Path})
	}
	ui.AddTable(0, 2, sitesList)
	prompt := "Please choose a site (default=1): "
	ui.AddLabel(0, 3+len(sitesList), prompt)
	ui.MoveCursor(len(prompt), 3+len(sitesList))

	ui.Draw()

	siteSelect := getSelection(len(sitesList))

	site := hugo.Load(sites[siteSelect].Path)
	return site
}

func getSelection(max int) int {
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
	err := editorCmd.Start()
	err = editorCmd.Wait()
	check(err)
}

func siteOverview(site hugo.Blog) {
	ui.Reset()
	posts := site.Posts
	var postList [][]string
	postList = append(postList, []string{"#", "Date", "Title", "Draft"})
	for i, post := range posts {
		draftStatus := "False"
		if post.Draft {
			draftStatus = "True"
		}
		postList = append(postList, []string{fmt.Sprintf("%d", i+1), post.Date, post.Title, draftStatus})
	}

	ui.AddLabel(0, 0, "Posts from the blog:")
	ui.AddTable(0, 1, postList)

	prompt := "Select a post to edit: "
	ui.AddLabel(0, 3+len(postList), prompt)
	ui.Draw()

	ui.MoveCursor(len(prompt), 3+len(postList))
	postNum := getSelection(len(postList))
	startEditor(site.Posts[postNum].Path)
}
