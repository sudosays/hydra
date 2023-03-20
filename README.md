# Hydra: A command-line manager for Hugo

The aim of hydra is to keep it simple and clean, but still provide a slick
blogging experience.

Currently, the interface is a very basic REPL with ideas to have a TUI in the future.

## Roadmap
* [x] Load and display posts
* [x] Edit posts from list
* [x] Manage multiple sites specified in config file
* [x] Create new posts
* [x] Delete posts
* [ ] Browse/sort by date, draft status

Future features:
* [ ] Publish drafts from post list
* [ ] Synchronise using Git
    * [ ] Github pages deploy
* [ ] Tag/category manager
* [ ] Interactive TUI 

## Prerequisites

Dependencies:
```
hugo >= v0.80.0
```

The now shelved TUI for hydra was built with the wonderful
[tcell](https://github.com/gdamore/tcell) package by [Garret
D'Amore](https://github.com/gdamore/tcell).

## Configuration

Currently, hydra does not have a configuration wizard, but it automatically
loads the configuration file at `~/.config/hydra.json`. The config file itself
is very straightforward:

``` json
{   
    "extension": "org",
    "editor":{
        "command": "/path/to/editor",
        "args": ""
    },
    "sites": [
        {"name":"Site one", "path":"/path/to/site/"},
        {"name":"Site two", "path":"/path/to/site2/"}
    ]
}
```

Note: the name of the site in the config file can be any name that you choose. It is there to help you distinguish between different sites.


## Installation

If you have added the Go install directory to your `$PATH` then installation is as easy as:

```
go install
```

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
