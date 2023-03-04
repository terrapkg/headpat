package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-martini/martini"
)

var fs = memfs.New() // I need help on how to actually interact with the fs aaa
var gitstore = memory.NewStorage()

func init() {
	if _, err := os.Stat("./terra/"); os.IsNotExist(err) {
		log.Println("Cloning repo...")
		_, e := git.Clone(gitstore, fs, &git.CloneOptions{
			URL: "git@github.com:terrapkg/terra.git",
		})
		if e != nil {
			log.Fatal(e);
			os.Exit(1)
		}
	} else {
		r, e := git.Open(gitstore, fs)
		if e != nil {
			log.Fatal(e)
			os.Exit(1)
		}
		wt, e := r.Worktree()
		if e != nil {
			log.Fatal(e)
			os.Exit(1)
		}
		switch wt.Pull(&git.PullOptions{}) {
			case nil, git.NoErrAlreadyUpToDate: break
			default: {
				log.Fatal(e)
				os.Exit(1)
			}
		}
	}
}

func gh(req *http.Request) int {
	if req.Header.Get("X-GitHub-Event") != "push" {
		return 400
	}

	type GhCommit struct {
		modified []string
	}
	type GhResp struct {
		ref     string
		commits []GhCommit
	}
	decoder := json.NewDecoder(req.Body)
	var resp GhResp
	decoder.Decode(&resp)
	return 204
}

func main() {
	m := martini.Classic()
	m.Get("/", func() (int, string) {
		return 418, "にゃお〜"
	})
	m.Post("/gh", gh)
	m.Run()
}
