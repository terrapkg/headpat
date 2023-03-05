package main

import "github.com/go-martini/martini"

var g_anitya_ch = make(chan string)

// I know global vars are horrible but idk

func main() {
	go anitnya_conn()

	m := martini.Classic()
	m.Get("/", func() (int, string) {
		return 418, "にゃお〜"
	})
	m.Post("/gh", gh)
	m.Run()
}
