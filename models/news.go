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

	this.OpenDB()
	defer this.CloseDB()
	c := this.DB.C("news")

	nid := this.GetID("news")
	title := r.FormValue("title")
	news := r.FormValue("news")
	create_time := this.GetTime()

	err := c.Insert(bson.M{
		"nid":         nid,
		"title":       title,
		"news":        news,
		"status":      1,
		"create_time": create_time,
	})
	if err != nil {
		log.Println(err)
	} else {
		b, _ := json.Marshal(map[string]interface{}{
			"nid":    nid,
			"ok":     1,
			"status": 1,
		})
		w.Write(b)
	}

}
