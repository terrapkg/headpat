package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
)

func gh(req *http.Request) int {
	if req.Header.Get("X-GitHub-Event") != "push" {
		return 400
	}
	type GhResp struct {
		ref string
		compare string
	}
	decoder := json.NewDecoder(req.Body)
	var resp GhResp
	decoder.Decode(&resp)
	// resp.compare + ".diff"
	// parse diff to find out changed files?
	// or use git directly? dunno
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
