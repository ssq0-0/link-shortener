package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type SafeWriter struct {
	mu    sync.RWMutex
	Links map[string]string
}

var globalSafeWriter = SafeWriter{Links: make(map[string]string)}

func (s *SafeWriter) Parser(link string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	strData := string(link)
	fmt.Println(strData)
	if value, exists := s.Links[strData]; exists {
		return value, nil
	} else {
		shortLink, err := coder(strData)
		if err != nil {
			return "", err
		}
		s.Links[strData] = shortLink
		return shortLink, nil
	}
}

func (s *SafeWriter) CheckerInMap(token string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Links == nil {
		return "", errors.New("Links map is not initialized")
	}
	for key, val := range s.Links {
		if val == token {
			return key, nil
		}
	}

	return "", errors.New("token not found in map")
}

func RegisterHandlers() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
		}
		defer r.Body.Close()

		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
			return
		}

		link, ok := requestData["link"].(string)
		if !ok {
			http.Error(w, "Link field is missing or not a string", http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			response, err := globalSafeWriter.Parser(link)
			if err != nil {
				http.Error(w, "Can't short this link", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Shorted link: http://127.0.0.1:8000/%s", response)
		}(link)
		wg.Wait()
		fmt.Println(globalSafeWriter.Links)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.Path, "/")
		var wg sync.WaitGroup

		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			response, err := globalSafeWriter.CheckerInMap(token)
			if err != nil || response == "" {
				fmt.Println(response)
				http.Error(w, "ShortLink time expired", http.StatusBadRequest)
				return
			} else {
				http.Redirect(w, r, response, http.StatusFound)
			}
		}(token)
		wg.Wait()
	})
}
