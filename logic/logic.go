package logic

import (
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	Urls map[string]string
}

func RouteIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func (us *URLShortener) RouteSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("view.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "URL parameter is missing", http.StatusBadRequest)
			return
		}

		shortKey := generateShortKey()
		us.Urls[shortKey] = originalURL

		shortUrl := "http://localhost:8080/short/" + shortKey

		var data = map[string]string{"originalURL": originalURL, "shortUrl": shortUrl}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func (us URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "short key is missing", http.StatusBadRequest)
		return
	}

	originalURL, exists := us.Urls[shortKey]
	if !exists {
		http.Error(w, "short key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
