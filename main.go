package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/user/", userHandler)

	http.ListenAndServe(":8888", nil)
}
