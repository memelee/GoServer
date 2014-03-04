package models

import (
	"GoServer/config"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type solution struct {
	Sid int `json:"sid"bson:"sid"`

	Pid      int    `json:"pid"bson:"pid"`
	Uid      string `json:"uid"bson:"uid"`
	Judge    int    `json:"judge"bson:"judge"`
	Time     int    `json:"time"bson:"time"`
	Memory   int    `json:"memory"bson:"memory"`
	Length   int    `json:"length"bson:"length"`
	Language int    `json:"language"bson:"language"`

	Module int `json:"module"bson:"module"`
	Mid    int `json:"mid"bson:"mid"`

	Code string `json:"code"bson:"code"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:"create"`
}

var sDetailSelector = bson.M{"_id": 0}
var sAchieveSelector = bson.M{"_id": 0, "pid": 1}
var sListSelector = bson.M{"_id": 0, "sid": 1, "pid": 1, "uid": 1, "judge": 1, "time": 1, "memory": 1, "length": 1, "language": 1, "create": 1}

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

	w.WriteHeader(200)
	w.Write(b)
}

// POST /solution/delete/sid/<sid>
func (this *Solution) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Delete")
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
func (this *Solution) Insert(w http.ResponseWriter, r *http.Request) {
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

	one.Status = 1
	one.Create = this.GetTime()
	one.Sid, err = this.GetID("solution")
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

	alt := make(map[string]interface{})
	if ori.Judge > 0 {
		alt["judge"] = ori.Judge
	}
	if ori.Time >= 0 {
		alt["time"] = ori.Time
	}
	if ori.Memory >= 0 {
		alt["memory"] = ori.Memory
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("solution").Update(bson.M{"sid": sid}, bson.M{"$set": alt})
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

// POST /solution/count/pid/<pid>/uid/<uid>/action/<accept/solve/submit>
func (this *Solution) Count(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Count")
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

	var count int
	c := this.DB.C("solution")
	switch v := args["action"]; v {
	case "submit":
		count, err = c.Find(query).Count()
		if err != nil {
			http.Error(w, "query error", 500)
			return
		}
	case "accept":
		query["judge"] = config.JudgeAC
		count, err = c.Find(query).Count()
		if err != nil {
			http.Error(w, "query error", 500)
			return
		}
	case "solve":
		var list []int
		query["judge"] = config.JudgeAC
		err = c.Find(query).Distinct("pid", &list)
		if err != nil {
			http.Error(w, "query error", 500)
			return
		}
		count = len(list)
	default:
		http.Error(w, "args error", 400)
		return
	}

	b, err := json.Marshal(map[string]interface{}{"count": count})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /solution/achieve/uid/<uid>
func (this *Solution) Achieve(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution Achieve")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid := args["uid"]

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	var list []int
	err = this.DB.C("solution").Find(bson.M{"uid": uid, "judge": config.JudgeAC}).Sort("pid").Distinct("pid", &list)
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

// POST /solution/list/offset/<offset>/limit/<limit>/sid/<sid>/pid/<pid>/uid/<uid>/language/<language>/judge/<judge>/module/<module>/mid/<mid>/from/<from>
func (this *Solution) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Solution List")
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

	q := this.DB.C("solution").Find(query).Select(sListSelector).Sort("-sid")

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
		var pid int
		pid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["pid"] = pid
	}
	if v, ok := args["uid"]; ok {
		query["uid"] = v
	}
	if v, ok := args["language"]; ok {
		var language int
		language, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["language"] = language
	}
	if v, ok := args["judge"]; ok {
		var judge int
		judge, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["judge"] = judge
	}
	if v, ok := args["module"]; ok {
		var module int
		module, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["module"] = module
	}
	if v, ok := args["mid"]; ok {
		var mid int
		mid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["mid"] = mid
	}
	if v, ok := args["from"]; ok {
		var from int
		from, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["sid"] = bson.M{"$gte": from}
	}
	return
}
