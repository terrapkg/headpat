package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func dl(fullURLFile string) chan string {
	l := log.New(os.Stdout, "[dwnload] "+fullURLFile+": ", 0)
	r := make(chan string)
	go func() {
		//? https://golangdocs.com/golang-download-files

		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := client.Get(fullURLFile)
		if err != nil {
			l.Fatal(err)
		}
		l.Printf("got %d", resp.StatusCode)
		defer resp.Body.Close()
		s, err := io.ReadAll(resp.Body)
		if err != nil {
			l.Fatal(err)
		}
		r <- string(s)
	}()
	return r
}

type GhCommit struct {
	Modified []string
}
type GhResp struct {
	Ref     string
	Commits []GhCommit
}

func gh(req *http.Request) int {
	if req.Header.Get("X-Hub-Signature") != os.Getenv("GH_WEBHOOK_SECRET") {
		return 403
	}
	if req.Header.Get("X-GitHub-Event") != "push" {
		return 400
	}
	decoder := json.NewDecoder(req.Body)
	var resp GhResp
	decoder.Decode(&resp)
	if resp.Ref != "refs/heads/main" { // for the time being
		return 204
	}
	var files []string
	for _, commit := range resp.Commits {
		files = append(files, commit.Modified...)
	}

	var chls []chan string
	for _, file := range files {
		if !strings.HasSuffix(file, "/pat") {
			continue
		}
		chls = append(chls, dl("https://raw.githubusercontent.com/terrapkg/packages/"+resp.Ref+"/"+file))
	}
	for _, chl := range chls {
		content := <-chl
		g_anitya_ch <- content
	}

	return 204
}
