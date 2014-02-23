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

var nDetailSelector = bson.M{"_id": 0}
var nListSelector = bson.M{"_id": 0, "nid": 1, "title": 1, "content": 1, "status": 1}

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

// POST /news/delete/nid/<nid>
func (this *News) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News Delete")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	nid, err := strconv.Atoi(args["nid"])
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

	err = this.DB.C("news").Remove(bson.M{"nid": nid})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "delete error", 599)
		return 
	}

	w.WriteHeader(200)
}

// POST /news/update/nid/<nid>
func (this *News) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News Update")
	this.Init(w, r)

	args := this.PraseURL(r.URL.Path[2:])
	nid, err := strconv.Atoi(args["nid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var ori news
	err = this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	var alt map[string]interface{}
	if ori.Title != "" {
		alt["title"] = ori.Title
	}
	if ori.Content != "" {
		alt["content"] = ori.Content
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	err = this.DB.C("news").Update(bson.M{"nid": nid}, bson.M{"$set": alt})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 599)
		return 
	}

	w.WriteHeader(200)
}

// POST /news/detail/nid/<nid>
func (this *News) Detail(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News Detail")
	this.Init(w, r)

	args := this.PraseURL(r.URL.Path[2:])
	nid, err := strconv.Atoi(args["nid"])
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

	var one news
	err = this.DB.C("news").Find(bson.M{"nid": nid}).Select(nDetailSelector).One(&one)
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "detail error", 599)
	}

	b, err := json.Marshal(&one)
	if err != nil {
		http.Error(w, "json error", 599)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /news/status/nid/<nid>
func (this *News) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News Status")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	nid, err := strconv.Atoi(args["nid"])
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

	err = this.DB.C("news").Update(bson.M{"nid": nid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "status error", 599)
		return
	}

	w.WriteHeader(200)
}

// POST /news/list/offest/<offest>/limit/<limit>
func (this *News) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server News List")
	this.Init(w, r)

	args := this.PraseURL(r.URL.Path[2:])
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

	q := this.DB.C("news").Find(1).Select(nListSelector).Sort({"pid": -1})

	if v, ok := args["offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "args error", 400)
			return
		}
		q = q.Skip(offset)
	}

	if v, ok := args["limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "args error". 400)
			return 
		}
		q = q.Limit(limit)
	}

	var list []*news
	err = q.All(&list)
	if err != nil {
		http.Error(w, "query error", 599)
		return 
	}

	b, err := json.Marshal(map[string]interface{}{"list": list})
	if err != nil {
		http.Error(w, "json error", 599)
		return 
	}

	w.WriteHeader(200)
	w.Write(b)
}