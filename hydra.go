package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sudosays/hydra/internal/ui"
	"github.com/sudosays/hydra/pkg/data/hugo"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
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
	defaultConfigFilePath := path.Join(userHomePath, ".config", "hydra.toml")

	// Parse command line flags if any
	configFilePath := flag.String("config", defaultConfigFilePath, "Path to a config file")
	flag.Parse()

	config, err = readConfig(*configFilePath)
	check(err)

	fmt.Printf("[Config editor] Command: %v; Args: %v;\n", config.Editor.Command, config.Editor.Args)

}

func main() {

	clearTerm()

	blog := hugo.Load(config.Sites["main"].Path)

	posts := blog.Posts

	p := tea.NewProgram(ui.InitialModel(posts))
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	clearTerm()

	fmt.Println("You have slain the hydra...")
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

func clearTerm() {
	fmt.Print("\033[H\033[2J")
}
