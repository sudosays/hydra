package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Site struct {
	Name, Path string
}

type Post struct {
	Title, Date, Path string
	Draft             bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printPosts(posts []Post) {
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

func hugoConfig() string {
	s := ""
	hugoConfigCmd := exec.Command("hugo", "config")
	hugoConfig, err := hugoConfigCmd.Output()
	s = string(hugoConfig)
	check(err)

	return s
}

func hugoGetPosts() []Post {
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
	screen := initScreen()
	sites := readConfig(configPath)
	line := 0
	putStr(0, line, screen, fmt.Sprintf("Current working dir: %s", cwd))
	line++
	putStr(0, line, screen, fmt.Sprintf("Config path: %s", configPath))
	line++
	putStr(0, line, screen, fmt.Sprintf("Sites from config:"))
	line++
	showSiteList(0, line, screen, sites)
	line += len(sites)
	putStr(0, line, screen, fmt.Sprintf("Please choose a site (default=1): "))
	choice := 0
	var posts []Post
	if choice < len(sites) {
		os.Chdir(sites[choice].Path)
		//siteDir, err := os.Getwd()
		check(err)
		//siteConfig := hugoConfig()
		//fmt.Println(siteConfig)
		posts = hugoGetPosts()
		//printPosts(posts)
		//startEditor(posts[0].Path)
	}
	var postList []string
	for _, p := range posts {
		indicator := ""
		if p.Draft {
			indicator = "*"
		}
		postList = append(postList, fmt.Sprintf("%s%s\n", indicator, p.Title))
	}

	for {
		printList(0, line+1, screen, postList)
		screen.Sync()
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune {
				if ev.Rune() == 'q' {
					screen.Fini()
					os.Exit(0)
				}
			}
		}

	}
}

func putRune(x, y int, s tcell.Screen, r rune) {
	s.SetContent(x, y, r, []rune{}, tcell.StyleDefault)
}

func putStr(x, y int, s tcell.Screen, str string) {
	for _, c := range str {
		putRune(x, y, s, rune(c))
		x++
	}
}

func printList(x, y int, s tcell.Screen, list []string) {
	for num, item := range list {
		outStr := fmt.Sprintf("%d) %s", num+1, item)
		putStr(x, y, s, outStr)
		y++
	}
}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)
	screen.Clear()
	return screen
}

func showSiteList(x, y int, screen tcell.Screen, sites []Site) {
	line := y
	var siteList []string
	if len(sites) > 0 {
		for _, s := range sites {
			siteList = append(siteList, s.Name)
		}
	}
	printList(0, line, screen, siteList)
}
