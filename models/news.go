package models

import (
	"encoding/json"
	// "labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type news struct {
	Nid     int    `json:"nid"bson:"nid"`
	Title   string `json:"title"bson:"title"`
	Content string `json:"content"bson:"content"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:'create'`
}

type News struct {
	Model
}

// POST /news/insert
func (this *News) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News Insert")
	this.Init(w, r)

	var one news
	err := this.LoadJson(r.Body, &one)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	one.Status = 1
	one.Create = this.GetTime()
	one.Nid, err = this.GetID("news")
	if err != nil {
		http.Error(w, "nid error", 599)
		return
	}

	err = this.DB.C("news").Insert(&one)
	if err != nil {
		http.Error(w, "insert error", 599)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"nid":    one.Nid,
		"status": one.Status,
	})
	if err != nil {
		http.Error(w, "json error", 599)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}
