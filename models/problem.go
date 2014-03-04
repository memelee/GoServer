package models

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type problem struct {
	Pid int `json:"pid"bson:"pid"`

	Time   int `json:"time"bson:"time"`
	Memory int `json:"memory"bson:"memory"`

	Title       string `json:"title"bson:"title"`
	Description string `json:"description"bson:"description"`
	Input       string `json:"input"bson:"input"`
	Output      string `json:"output"bson:"output"`
	Source      string `json:"source"bson:"source"`
	Hint        string `json:"hint"bson:"hint"`

	In  string `json:"in"bson:"in"`
	Out string `json:"out"bson:"out"`

	Solve  int `json:"solve"bson:"solve"`
	Submit int `json:"submit"bson:"submit"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:"create"`
}

var pDetailSelector = bson.M{"_id": 0}
var pListSelector = bson.M{"_id": 0, "pid": 1, "title": 1, "source": 1, "solve": 1, "submit": 1, "status": 1}

type Problem struct {
	Model
}

// POST /problem/detail/pid/<pid>

func (this *Problem) Detail(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Detail")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	pid, err := strconv.Atoi(args["pid"])
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

	var one problem
	err = this.DB.C("problem").Find(bson.M{"pid": pid}).Select(pDetailSelector).One(&one)
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

// POST /problem/delete/pid/<pid>
func (this *Problem) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Delete")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	pid, err := strconv.Atoi(args["pid"])
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

	err = this.DB.C("problem").Remove(bson.M{"pid": pid})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "delete error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /problem/insert
func (this *Problem) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Insert")
	this.Init(w, r)

	var one problem
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

	one.Solve = 0
	one.Submit = 0
	one.Status = 0
	one.Create = this.GetTime()
	one.Pid, err = this.GetID("problem")
	if err != nil {
		http.Error(w, "pid error", 500)
		return
	}

	err = this.DB.C("problem").Insert(&one)
	if err != nil {
		http.Error(w, "insert error", 500)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"pid":    one.Pid,
		"status": one.Status,
	})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /problem/update/pid/<pid>
func (this *Problem) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Update")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	pid, err := strconv.Atoi(args["pid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var ori problem
	err = this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	alt := make(map[string]interface{})
	if ori.Title != "" {
		alt["title"] = ori.Title
	}
	if ori.Description != "" {
		alt["description"] = ori.Description
	}
	if ori.Input != "" {
		alt["input"] = ori.Input
	}
	if ori.Output != "" {
		alt["output"] = ori.Output
	}
	if ori.Source != "" {
		alt["source"] = ori.Source
	}
	if ori.Hint != "" {
		alt["hint"] = ori.Hint
	}
	if ori.In != "" {
		alt["in"] = ori.In
	}
	if ori.Out != "" {
		alt["out"] = ori.Out
	}
	if ori.Time > 0 {
		alt["time"] = ori.Time
	}
	if ori.Memory > 0 {
		alt["memory"] = ori.Memory
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("problem").Update(bson.M{"pid": pid}, bson.M{"$set": alt})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /problem/status/pid/<pid>
func (this *Problem) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Status")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	pid, err := strconv.Atoi(args["pid"])
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

	err = this.DB.C("problem").Update(bson.M{"pid": pid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "status error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /problem/record/pid/<pid>/action/<solve/submit>
func (this *Problem) Record(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Prblem Record")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	pid, err := strconv.Atoi(args["pid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var inc int
	switch v := args["action"]; v {
	case "solve":
		inc = 1
	case "submit":
		inc = 0
	default:
		http.Error(w, "args error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("problem").Update(bson.M{"pid": pid}, bson.M{"$inc": bson.M{"solve": inc, "submit": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "record error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /problem/list/offset/<offset>/limit/<limit>/pid/<pid>/title/<title>/source/<source>
func (this *Problem) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem List")
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

	q := this.DB.C("problem").Find(query).Select(pListSelector).Sort("pid")

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

	var list []*problem
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

func (this *Problem) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)

	if v, ok := args["pid"]; ok {
		var pid int
		pid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["pid"] = pid
	}
	if v, ok := args["title"]; ok {
		query["title"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	if v, ok := args["source"]; ok {
		query["source"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	return
}
