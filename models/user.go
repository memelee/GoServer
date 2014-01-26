package models

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type User struct {
	Model
}

var uGetSelector = bson.M{"_id": 0}

func (this *User) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Login")
	this.Init(w, r)

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db errpr", 599)
		return
	}

	uid := r.FormValue("uid")
	pwd := r.FormValue("pwd")

	type result struct {
		Uid       string
		Pwd       string
		Privilege int
		Status    int
	}

	one := &result{}
	c := this.DB.C("user")
	c.Find(bson.M{"uid": uid}).Select(uGetSelector).One(one)

	var b []byte
	if pwd != "" && pwd == one.Pwd {
		log.Println("Server User Login Successfully")
		b, err = json.Marshal(map[string]interface{}{
			"uid":       one.Uid,
			"ok":        1,
			"privilege": one.Privilege,
			"status":    one.Status,
		})
	} else {
		log.Println("Server User Login Failed")
		b, err = json.Marshal(map[string]interface{}{
			"uid":       one.Uid,
			"ok":        0,
			"privilege": 0,
			"status":    0,
		})
	}
	if err != nil {
		http.Error(w, "json error", 599)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func (this *User) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Logout")
	this.Init(w, r)

	uid := r.FormValue("uid")
	log.Println("Server User Logout Successfully")

	b, err := json.Marshal(map[string]interface{}{
		"uid": uid,
		"ok":  1,
	})
	if err != nil {
		http.Error(w, "json error", 599)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}
