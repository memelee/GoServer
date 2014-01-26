package models

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type News struct {
	Model
}

func (this *News) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News Insert")
	this.Init(w, r)

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	title := r.FormValue("title")
	news := r.FormValue("news")
	createTime := this.GetTime()
	nid, err := this.GetID("news")
	if err != nil {
		http.Error(w, "nid error", 599)
		return
	}

	c := this.DB.C("news")
	err = c.Insert(bson.M{
		"nid":         nid,
		"title":       title,
		"news":        news,
		"status":      1,
		"create_time": createTime,
	})
	if err != nil {
		http.Error(w, "insert error", 599)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"nid":    nid,
		"ok":     1,
		"status": 1,
	})
	if err != nil {
		http.Error(w, "json error", 599)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}
