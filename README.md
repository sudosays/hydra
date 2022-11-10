# Hydra: A command-line manager for Hugo

The aim of hydra is to keep it simple and clean, but still provide a slick
blogging experience.

## Roadmap
* [X] Load and display posts
* [X] Edit posts from list
* [X] Manage multiple sites specified in config file
* [X] Create new posts
* [ ] Publish drafts from post list
* [ ] Github pages deploy
* [ ] Browse/sort by date, draft status
* [ ] Tag/category manager

## Prerequisites

Dependencies:
```
hugo >= v0.80.0
```

The now shelved TUI for hydra was built with the wonderful
[tcell](https://github.com/gdamore/tcell) package by [Garret
D'Amore](https://github.com/gdamore/tcell).

## Installation

_TBC_

### Configuration

Currently, hydra does not have a configuration wizard, but it automatically
loads the configuration file at `~/.config/hydra.json`. The config file itself
is very straightforward:

``` json
{   "extension": "org",
    "editor":"/path/to/editor",
    "sites": [
        {"name":"Site one", "path":"/path/to/site/"},
        {"name":"Site two", "path":"/path/to/site2/"}
    ]
}
```

Note: the name of the site in the config file can be any name that you choose. It is there to help you distinguish between different sites.

## Usage

```
make run
```
or 
```
make build
./bin/hydra
```

### Running tests

If you are hacking on hydra, there are some useful make rules to know about:

* `make test`: runs `go test` for every file
* `make verify`: runs `golint` for the project

## Contributing

No PRs will be accepted at this time, but you are more than welcome to open issues and muck about.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](./LICENSE) file for details.


## Default vs Interactive mode

By default, hydra lists posts to stdout and parses arguments to perform actions.

Interactive-mode waits for user prompts until quit.

The TUI mode is shelved for now, in favor of focusing on end-to-end usability
without getting bogged down into developing a TUI with good UX.
