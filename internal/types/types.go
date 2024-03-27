package types

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
