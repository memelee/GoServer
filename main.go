package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/news/", newsHandler)
	http.HandleFunc("/problem/", problemHandler)
	http.HandleFunc("/contest/", contestHandler)
	http.HandleFunc("/exercise/", exerciseHandler)

	http.ListenAndServe(":8888", nil)
}
