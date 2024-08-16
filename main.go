package main

import (
	"fmt"
	"net/http"

	"go-url-short/logic"
)

func main() {
	shortener := &logic.URLShortener{
		Urls: make(map[string]string),
	}

	http.HandleFunc("/", logic.RouteIndexGet)
	http.HandleFunc("/process", shortener.RouteSubmitPost)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	fmt.Println("URL Shortener is running on :8080")
	http.ListenAndServe(":8080", nil)
}
