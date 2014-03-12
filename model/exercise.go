package model

import (
	"GoServer/class"
	"GoServer/config"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type exercise struct {
	Eid      int         `json:"eid"bson:"eid"`
	Title    string      `json:"title"bson:"title"`
	Encrypt  int         `json:"encrypt"bson:"encrypt"`
	Argument interface{} `json:"argument"bson:"argument"`

	Start string `json:"start"bson:"start"`
	End   string `json:"end"bson:"end"`

	Status int    `json:"status"bson:"status"`
	Create string `'json:"create"bson:"create"`

	List []int `json:"list"bson:"list"`
}

var eDetailSelector = bson.M{"_id": 0}
var eListSelector = bson.M{"_id": 0, "eid": 1, "title": 1, "encrypt": 1, "argument": 1, "start": 1, "end": 1, "status": 1}

type Exercise struct {
	class.Model
}

// POST /exercise/detail/eid/<eid>
func (this *Exercise) Detail(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise Detail")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	eid, err := strconv.Atoi(args["eid"])
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

	var one exercise
	err = this.DB.C("exercise").Find(bson.M{"eid": eid}).Select(eDetailSelector).One(&one)
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

// POST /exercise/delete/eid/<eid>
func (this *Exercise) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise Delete")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	eid, err := strconv.Atoi(args["eid"])
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

	err = this.DB.C("exercise").Remove(bson.M{"eid": eid})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "delete error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /exercise/insert
func (this *Exercise) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise Insert")
	this.Init(w, r)

	var one exercise
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
	one.Eid, err = this.GetID("exercise")
	if err != nil {
		http.Error(w, "eid error", 500)
		return
	}

	err = this.DB.C("exercise").Insert(&one)
	if err != nil {
		http.Error(w, "insert error", 500)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"eid":    one.Eid,
		"status": one.Status,
	})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /exercise/update/eid/<eid>
func (this *Exercise) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise Update")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	eid, err := strconv.Atoi(args["eid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var ori exercise
	err = this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	alt := make(map[string]interface{})
	if ori.Title != "" {
		alt["title"] = ori.Title
	}
	if ori.Start != "" {
		alt["start"] = ori.Start
	}
	if ori.End != "" {
		alt["end"] = ori.End
	}
	if ori.Encrypt > config.EncryptNA {
		alt["encrypt"] = ori.Encrypt
		alt["Argument"] = ori.Argument
	}
	if ori.List != nil {
		alt["list"] = ori.List
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("exercise").Update(bson.M{"eid": eid}, bson.M{"$set": alt})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /exercise/status/eid/<eid>
func (this *Exercise) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise Status")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	eid, err := strconv.Atoi(args["eid"])
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

	err = this.DB.C("exercise").Update(bson.M{"eid": eid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "status error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /exercise/push/eid/<eid>
func (this *Exercise) Push(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise Push")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	eid, err := strconv.Atoi(args["eid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var one exercise
	err = this.LoadJson(r.Body, &one)
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

	err = this.DB.C("exercise").Update(bson.M{"eid": eid}, bson.M{"$addToSet": bson.M{"list": bson.M{"$each": one.List}}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /exercise/list/offset/<offset>/limit/<limit>/pid/<pid>/title/<title>
func (this *Exercise) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Exercise List")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	query, err := this.CheckQuery(args)
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

	q := this.DB.C("exercise").Find(query).Select(eListSelector).Sort("eid")

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
			http.Error(w, "args error", 400)
			return
		}
		q = q.Limit(limit)
	}

	var list []*exercise
	err = q.All(&list)
	if err != nil {
		http.Error(w, "query error", 500)
		return
	}

	b, err := json.Marshal(map[string]interface{}{"list": list})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func (this *Exercise) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)

	if v, ok := args["eid"]; ok {
		var eid int
		eid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["eid"] = eid
	}
	if v, ok := args["title"]; ok {
		query["title"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	return
}
