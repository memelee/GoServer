package model

import (
	"GoServer/class"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type contest struct {
	Cid      int         `json:"cid"bson:"cid"`
	Title    string      `json:"title"bson:"title"`
	Encrypt  int         `json:"encrypt"bson:"encrypt"`
	Argument interface{} `json:"argument"bson:"argument"`

	Start string `json:"start"bson:"start"`
	End   string `json:"end"bson:"end"`

	Status int    `json:"status"bson:"status"`
	Create string `'json:"create"bson:"create"`

	List []int `json:"list"bson:"list"`
}

var cDetailSelector = bson.M{"_id": 0}
var cListSelector = bson.M{"_id": 0, "cid": 1, "title": 1, "encrypt": 1, "argument": 1, "start": 1, "end": 1, "status": 1}

type Contest struct {
	class.Model
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

// POST /contest/update/cid/<cid>
func (this *Contest) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest Update")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	cid, err := strconv.Atoi(args["cid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var ori contest
	err = this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	alt := make(map[string]interface{})
	alt["title"] = ori.Title
	alt["start"] = ori.Start
	alt["end"] = ori.End
	alt["encrypt"] = ori.Encrypt
	alt["Argument"] = ori.Argument
	alt["list"] = ori.List

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("contest").Update(bson.M{"cid": cid}, bson.M{"$set": alt})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /contest/status/cid/<cid>
func (this *Contest) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest Status")
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

	err = this.DB.C("contest").Update(bson.M{"cid": cid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "status error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /contest/push/cid/<cid>
func (this *Contest) Push(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest Push")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	cid, err := strconv.Atoi(args["cid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var one contest
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

	err = this.DB.C("contest").Update(bson.M{"cid": cid}, bson.M{"$addToSet": bson.M{"list": bson.M{"$each": one.List}}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /contest/list/offset/<offset>/limit/<limit>/pid/<pid>/title/<title>
func (this *Contest) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Contest List")
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

	q := this.DB.C("contest").Find(query).Select(cListSelector).Sort("-cid")

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

	var list []*contest
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

func (this *Contest) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)

	if v, ok := args["cid"]; ok {
		var cid int
		cid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["cid"] = cid
	}
	if v, ok := args["title"]; ok {
		query["title"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	return
}
