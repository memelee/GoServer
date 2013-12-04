package main

import (
	"GoServer/models"
	"net/http"
	"reflect"
	"strings"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	s := strings.Split(p, "/")
	if l := len(s); l >= 2 {
		c := &models.User{}
		m := strings.Title(s[1])
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

// Common

func callMethod(c interface{}, m string, rv []reflect.Value) {
	rc := reflect.ValueOf(c)
	rm := rc.MethodByName(m)
	rm.Call(rv)
}

func getReflectValue(w http.ResponseWriter, r *http.Request) []reflect.Value {
	rw := reflect.ValueOf(w)
	rr := reflect.ValueOf(r)
	return []reflect.Value{rw, rr}
}
