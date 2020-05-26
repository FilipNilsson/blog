package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	articles := make(map[string]Article)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Go to /articles")
	})

	r.HandleFunc("/articles/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := strings.ReplaceAll(mux.Vars(r)["key"], "+", " ")
		if _, ok := articles[key]; !ok {
			fmt.Fprintf(w, "No article with name %s found\n", key)
			return
		}
		if r.Method == "GET" {
			fmt.Fprintf(w, "Title: %s\n", key)
			fmt.Fprintf(w, "Author: %s\n", articles[key].Author)
			fmt.Fprintf(w, "Content: %s\n", articles[key].Content)
			return
		}
		if r.Method == "DELETE" {
			delete(articles, key)
			fmt.Fprintf(w, "Removed article %s\n", key)
			return
		}
	})

	r.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			for k, v := range articles {
				fmt.Fprintf(w, "Title: %s\n", k)
				fmt.Fprintf(w, "Author: %s\n", v.Author)
				fmt.Fprintf(w, "Content: %s\n", v.Content)
				fmt.Fprintf(w, "========================\n")
			}
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Invalid request received!")
		}
		fmt.Println("Body:", string(body))
		var body_map map[string]string
		json.Unmarshal([]byte(body), &body_map)

		if _, ok := body_map["author"]; !ok {
			fmt.Println("No author in data")
			return
		}

		if _, ok := body_map["title"]; !ok {
			fmt.Println("No title in data")
			return
		}

		if _, ok := body_map["content"]; !ok {
			fmt.Println("No content in data")
			return
		}
		fmt.Println("author:", body_map["author"])
		fmt.Println("title:", body_map["title"])
		fmt.Println("content:", body_map["content"])

		articles[body_map["title"]] = Article{
			Author:  body_map["author"],
			Content: body_map["content"],
		}
	})

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

type Article struct {
	Author  string
	Content string
}
