package models

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type user struct {
	Uid    string `json:"uid"bson:"uid"`
	Pwd    string `json:"pwd"bson:"pwd"`
	Nick   string `json:"nick"bson:"nick"`
	Mail   string `json:"mail"bson:"mail"`
	School string `json:"school"bson:"school"`

	Privilege int `json:"privilege"bson:"privilege"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:'create'`
}

type User struct {
	Model
}

var uDetailSelector = bson.M{"_id": 0}
var uListSelector = bson.M{"_id": 0, "uid": 1, "nick": 1, "status": 1}

// POST /user/login
func (this *User) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Login")
	this.Init(w, r)

	var ori user
	err := this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	var alt user
	c := this.DB.C("user")
	c.Find(bson.M{"uid": ori.Uid}).Select(uDetailSelector).One(&alt)

	var b []byte
	if ori.Pwd != "" && ori.Pwd == alt.Pwd {
		log.Println("Server User Login Successfully")
		b, err = json.Marshal(map[string]interface{}{
			"uid":       alt.Uid,
			"privilege": alt.Privilege,
			"status":    alt.Status,
		})
	} else {
		log.Println("Server User Login Failed")
		b, err = json.Marshal(map[string]interface{}{
			"uid":       "",
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

// POST /problem/logout
func (this *User) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Logout")
	this.Init(w, r)

	var one user
	err := this.LoadJson(r.Body, &one)
	if err != nil {
		http.Error(w, "load error", 400)
	}

	w.WriteHeader(200)
}

// POST /user/status/uid/<uid>
func (this *User) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Status")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid, err := strconv.Atoi(args["uid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	err = this.DB.C("user").Update(bson.M{"uid": uid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 400)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "status error", 599)
		return
	}

	w.WriteHeader(200)
}
