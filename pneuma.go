package main

import (
	//"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"io"
	"os"
	"os/exec"
	//"path"
	//"strconv"
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
	/*home, err := os.UserHomeDir()
	configPath := path.Join(home, ".config", "pneuma.json")
	check(err)
	os.Chdir(home)
	cwd, err := os.Getwd()
	check(err)
	fmt.Println("Current working dir: ", cwd)
	fmt.Println("Config path: ", configPath)
	fmt.Println("Sites from config:")
	fmt.Println("------------------")
	sites := readConfig(configPath)
	if len(sites) > 0 {
		for i, s := range sites {
			fmt.Printf("%d) %s\n", i+1, s.Name)
			fmt.Printf("Path: %s\n\n", s.Path)
		}
	}
	fmt.Printf("Please choose a site (default=1): ")
	inputReader := bufio.NewReader(os.Stdin)
	inputString, err := inputReader.ReadString('\n')
	check(err)
	choice, err := strconv.Atoi(inputString[:len(inputString)-1])
	choice--
	check(err)
	if choice < len(sites) {
		os.Chdir(sites[choice].Path)
		//siteDir, err := os.Getwd()
		check(err)
		//siteConfig := hugoConfig()
		//fmt.Println(siteConfig)
		posts := hugoGetPosts()
		printPosts(posts)
		//startEditor(posts[0].Path)
	}*/
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)
	screen.Clear()
	for {
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
