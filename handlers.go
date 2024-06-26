package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var LINKS = make(map[string]string)

type LinkRequest struct {
	Link string `json:"link"`
}

func parseLink(link string) (string, error) {
	if value, exist := LINKS[link]; exist {
		return value, nil
	} else {
		shortLink, err := coder(link)
		if err != nil {
			return "", err
		}
		LINKS[link] = shortLink
		return shortLink, nil
	}
}

func RegisterHandlers() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		var req LinkRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid jsob", http.StatusBadRequest)
			log.Println("Invalid json")
			return
		}

		link := req.Link

		if link == "" {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}
		result, err := parseLink(link)
		if err != nil {
			http.Error(w, "Error:", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Shorted link: http://127.0.0.1:8000/%s", result)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.Path, "/")

		var originalLink string
		for key, val := range LINKS {
			if val == token {
				originalLink = key
				break
			}
		}

		if originalLink != "" {
			http.Redirect(w, r, originalLink, http.StatusFound)
		} else {
			http.Error(w, "Link not found", http.StatusNotFound)
		}
	})
}
