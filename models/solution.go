package models

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type solution struct {
	Sid int `json:"sid"bson:"sid"`

	Pid         int    `json:"pid"bson:"pid"`
	JudgeStatus int    `json:"judgeStatus"bson:"judgeStatus"`
	Time        int    `json:"time"bson:"time"`
	Memory      int    `json:"memory"bson:"memory"`
	CodeLen     int    `json:"codelen"bson:"codelen"`
	Language    string `json:"language"bson:"language"`
	Author      string `json:"author"bson:"author"`

	Code string `json:"code"bson:"code"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:"create"`
}

var sDetailSelector = bson.M{"_id": 0}
var sListSelector = bson.M{"_id": 0, "sid": 1, "pid": 1, "judgeStatus": 1, "time": 1, "memory": 1, "codelen": 1, "language": 1, "author": 1}

type Solution struct {
	Model
}

// POST /solution/detail/sid/<sid>
func (this *Solution) Detail(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Detail")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	sid, err := strconv.Atoi(args["sid"])
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

	var one solution
	err = this.DB.C("solution").Find(bson.M{"sid": sid}).Select(sDetailSelector).One(&one)
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

	w.WriterHeader(200)
	w.Write(b)
}

// POST /solution/delete/sid/<sid>
func (this *Solution) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Delete")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	sid, err := strconv.Atoi(args[2:])
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

	err = this.DB.C("solution").Remove(bson.M{"sid": sid})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "delete error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /solution/insert
func (this *Solution) Insert(w http.ResponseWriter, r http.Request) {
	log.Println("Server Solution Insert")
	this.Init(w, r)

	var one solution
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

	one.Status = 0
	one.Create = this.GetTime()
	one.Sid, err = this.GetID("problem")
	if err != nil {
		http.Error(w, "sid error", 500)
		return
	}

	err = this.DB.C("solution").Insert(&one)
	if err != nil {
		http.Error(w, "insert error", 500)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"sid":    one.Sid,
		"status": one.Status,
	})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /solution/update/sid/<sid>
func (this *Solution) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Update")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	sid, err := strconv.Atoi(args["sid"])
	if err != nil {
		http.Error(w, "args error", 400)
		return
	}

	var ori solution
	err = this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	var alt map[string]interface{}
	if ori.Pid != "" {
		alt["pid"] = ori.Pid
	}
	if ori.JudgeStatus != "" {
		alt["judgestatus"] = ori.JudgeStatus
	}
	if ori.Time != "" {
		alt["time"] = ori.Time
	}
	if ori.Memory != "" {
		alt["memory"] = ori.Memory
	}
	if ori.CodeLen != "" {
		alt["codeLen"] = ori.CodeLen
	}
	if ori.Language != "" {
		alt["language"] = ori.Language
	}
	if ori.Author != "" {
		alt["author"] = ori.Author
	}
	if ori.Code != "" {
		alt["code"] = ori.Code
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("solution").Update(bson.M{"sid": sid}, bson.M{"$set", alt})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /solution/status/sid/<sid>
func (this *Solution) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Status")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	sid, err := strconv.Atoi(args["sid"])
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

	err = this.DB.C("solution").Update(bson.M{"sid": sid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "status error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /solution/list/offest/<offest>/limit/<limit>/pid/<pid>/author/<author>/language/<language>/judgeStatus/<judgeStatus>/fromSid/<fromSid>
func (this *Solution) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution List")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	query, err := this.CheckQuery(args)

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	q := this.DB.C("solution").Find(bson.M{}).Select(sListSelector).Sort("-sid")

	if v, ok := args["offest"]; ok {
		offest, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "args error", 400)
			return
		}
		q = q.Skip(offest)
	}

	if v, ok := args["limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "args error", 400)
			return
		}
		q = q.Limit(limit)
	}

	var list []*solution
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

func (this *Solution) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)

	if v, ok := args["sid"]; ok {
		var sid int
		sid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["sid"] = sid
	}
	if v, ok := args["pid"]; ok {
		query["pid"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	if v, ok := args["author"]; ok {
		query["author"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	if v, ok := args["language"]; ok {
		query["language"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	if v, ok := args["judgeStatus"]; ok {
		query["judgeStatus"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	if v, ok := args["fromSid"]; ok {
		query["fromSid"] = bson.M{"$max": fromsid}
	}
	return
}
