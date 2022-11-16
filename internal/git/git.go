package git

import (
	"os/exec"
)

// Checks whether or not the path is a git repo
func IsRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	_, err := cmd.Output()
	if err == nil {
		return true
	} else {
		return false
	}
}

// Checks for changes in the remote git
// Process includes:
// 1. fetching remote repos
// 2. updating local repo with remote changes
func Update() {
	cmd := exec.Command("git", "fetch")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("git", "status", "--porcelain")
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

}
