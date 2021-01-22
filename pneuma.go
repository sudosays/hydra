package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	libui "github.com/sudosays/pneuma/internal/ui"
	"github.com/sudosays/pneuma/pkg/data/hugo"
	"io"
	//"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

// hugoSite contains the information for a hugo site listed in the config
type hugoSite struct {
	Name, Path string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var sites []hugoSite
var ui libui.PneumaUI

func init() {
	home, err := os.UserHomeDir()
	configPath := path.Join(home, ".config", "pneuma.json")
	check(err)
	sites = readConfig(configPath)
	ui = libui.Init()

}

func main() {
	quitCmdEventKey := libui.CommandKey{Key: tcell.KeyRune, Rune: 'q', Mod: tcell.ModNone}
	cmds := make(map[libui.CommandKey]func())
	cmds[quitCmdEventKey] = quit
	ui.SetCommands(cmds)

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

	headings := []string{"Choice", "Name", "Path"}

	var sitesList [][]string
	for i, site := range sites {
		sitesList = append(sitesList, []string{fmt.Sprintf("%d", i+1), site.Name, site.Path})
	}
	ui.AddTable(0, 2, headings, sitesList)
	prompt := "Please choose a site (default=1): "
	ui.AddLabel(0, 4+len(sitesList), prompt)
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

func readConfig(path string) []hugoSite {
	var sites []hugoSite
	configFile, err := os.Open(path)
	check(err)
	decoder := json.NewDecoder(configFile)
	for {
		var site hugoSite
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
	headings := []string{"#", "Date", "Title", "Draft"}
	var postList [][]string

	for i, post := range posts {
		draftStatus := "False"
		if post.Draft {
			draftStatus = "True"
		}
		datetime, _ := time.Parse(time.RFC3339, post.Date)
		date := datetime.Format("2006/01/02")
		postList = append(postList, []string{fmt.Sprintf("%d", i+1), date, post.Title, draftStatus})
	}

	ui.AddLabel(0, 0, "Posts from the blog:")
	table := ui.AddTable(0, 1, headings, postList)

	editPost := func() {
		path := posts[table.Index].Path
		startEditor(path)
	}

	nextCmdEventKey := libui.CommandKey{Key: tcell.KeyRune, Rune: 'j', Mod: tcell.ModNone}
	prevCmdEventKey := libui.CommandKey{Key: tcell.KeyRune, Rune: 'k', Mod: tcell.ModNone}
	quitCmdEventKey := libui.CommandKey{Key: tcell.KeyRune, Rune: 'q', Mod: tcell.ModNone}
	enterCmdEventKey := libui.CommandKey{Key: tcell.KeyEnter, Rune: rune(13), Mod: tcell.ModNone}

	cmds := make(map[libui.CommandKey]func())
	cmds[quitCmdEventKey] = quit
	cmds[nextCmdEventKey] = table.NextItem
	cmds[prevCmdEventKey] = table.PreviousItem
	cmds[enterCmdEventKey] = editPost

	ui.SetCommands(cmds)

	prompt := "Select a post to edit: "
	ui.AddLabel(0, 4+len(postList), prompt)
	ui.Draw()

	ui.MoveCursor(len(prompt), 3+len(postList))
}

func quit() {
	ui.Close()
}
