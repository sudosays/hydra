package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sudosays/hydra/internal/types"
	"github.com/sudosays/hydra/internal/ui"
	"github.com/sudosays/hydra/pkg/data/hugo"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var config types.HydraConfig

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

	p := tea.NewProgram(ui.InitialModel(blog, config.Editor))
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	clearTerm()

	fmt.Println("You have slain the hydra...")
}

func readConfig(path string) (types.HydraConfig, error) {
	conf := types.HydraConfig{}
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
