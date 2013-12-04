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

func (this *User) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Login")
	this.Init(w, r)

	uid := r.FormValue("uid")
	pwd := r.FormValue("pwd")

	log.Println(uid, pwd)

	this.OpenDB()
	c := this.DB.C("user")
	defer this.CloseDB()

	type result struct {
		Uid       string
		Pwd       string
		Privilege int
		Status    int
	}

	one := result{}
	c.Find(bson.M{"uid": uid}).One(&one)

	var out map[string]interface{}
	if pwd != "" && pwd == one.Pwd {
		log.Println("Server User Login Successfully")
		out = map[string]interface{}{
			"uid":       one.Uid,
			"ok":        1,
			"privilege": one.Pwd,
			"status":    one.Status,
		}
	} else {
		log.Println("Server User Login Failed")
		out = map[string]interface{}{
			"uid":       one.Uid,
			"ok":        0,
			"privilege": 0,
			"status":    0,
		}
	}
	b, _ := json.Marshal(out)
	w.Write(b)
}
