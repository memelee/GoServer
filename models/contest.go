package models

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type contest struct {
	Cid     int    `json:"cid"bson:"cid"`
	Title   string `json:"title"bson:"title"`
	Encrypt int    `json:"encrypt"bson:"encrypt"`
	Start   string `json:"start"bson:"start"`
	End     string `json:"end"bson:"end"`

	Status int    `json:"status"bson:"status"`
	Create string `'json:"create"bson:"create"`

	List []int `json:"list"bson:"list"`
}

var cDetailSelector = bson.M{"_id": 0}
var cListSelector = bson.M{"_id": 0, "cid": 1, "title": 1, "encrypt": 1, "start": 1, "end": 1, "status": 1}

type Contest struct {
	Model
}

// POST /contest/detail/cid/<cid>
func (this *Contest) Detail(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest Detail")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	cid, err := strconv.Atoi(args["cid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	var one contest
	err = this.DB.C("contest").Find(bson.M{"cid": cid}).Select(cDetailSelector).One(&one)
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "detail error", 500)
		return
	}

	b, err := json.Marshal(&one)
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /contest/delete/cid/<cid>
func (this *Contest) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest Delete")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	cid, err := strconv.Atoi(args["cid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("contest").Remove(bson.M{"cid": cid})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "delete error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /contest/insert
func (this *Contest) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest Insert")
	this.Init(w, r)

	var one contest
	err := this.LoadJson(r.Body, &one)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	one.Status = 1
	one.Create = this.GetTime()
	one.Cid, err = this.GetID("contest")
	if err != nil {
		http.Error(w, "cid error", 500)
		return
	}

	err = this.DB.C("contest").Insert(&one)
	if err != nil {
		http.Error(w, "insert error", 500)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"cid":    one.Cid,
		"status": one.Status,
	})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}
