package git

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

type RepoStatus struct {
	ModifiedFiles  []string
	UntrackedFiles []string
	CanPull        bool
	CanPush        bool
}

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

// Get a list of all of the files that have been added (untracked) or modified
// as well as if there are any pulls that can be made
// Returns a RepoStatus struct
func Status() string {
	cmd := exec.Command("git", "status", "--porcelain")
	result, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(result)
}

func Sync() {
	status := Status()
	if len(status) > 0 {
		//cmd := exec.Command("git", "add", "--all")
		msg := fmt.Sprintf("hydra sync at %d\n\n%s", time.Now().Unix(), status)
		log.Println(msg)
	}
	// if len(status.ModifiedFiles) > 0 || len(status.UntrackedFiles) > 0 {
	// }

}
