package models

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type Problem struct {
	Model
}

var pGetSelector = bson.M{"_id": 0}
var pListSelector = bson.M{"_id": 0, "pid": 1, "title": 1, "source": 1, "solve": 1, "submit": 1, "status": 1}

// POST /problem/insert
func (this *Problem) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Insert")
	this.Init(w, r)

	this.OpenDB()
	defer this.CloseDB()
	c := this.DB.C("problem")

	title := r.FormValue("title")
	source := r.FormValue("source")
	time := r.FormValue("time")
	memory := r.FormValue("memory")
	description := r.FormValue("description")
	input := r.FormValue("input")
	output := r.FormValue("output")
	sampleInput := r.FormValue("sampleInput")
	sampleOutput := r.FormValue("sampleOutput")
	hint := r.FormValue("hint")
	createTime := this.GetTime()
	pid, err := this.GetID("problem")
	if err != nil {
		http.Error(w, "pid error", 599)
		return
	}

	err = c.Insert(bson.M{
		"pid":           pid,
		"title":         title,
		"time":          time,
		"memory":        memory,
		"description":   description,
		"input":         input,
		"output":        output,
		"sample_input":  sampleInput,
		"sample_output": sampleOutput,
		"source":        source,
		"hint":          hint,
		"solve":         0,
		"submit":        0,
		"status":        1,
		"create_time":   createTime,
	})
	if err != nil {
		http.Error(w, "insert error", 599)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"pid":    pid,
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

// GET /problem/list/offset/<offset>/limit/<limit>/source/<source>
func (this *Problem) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem List")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	query, err := this.CheckQuery(args)
	if err != nil {
		http.Error(w, "args error", 400)
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	c := this.DB.C("problem")
	q := c.Find(query).Select(pListSelector).Sort("pid")
	if v, ok := args["limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "args error", 400)
			return
		}
		q = q.Limit(limit)
	}

	type problem struct {
		Pid    int
		Title  string
		Source string
		Solve  int
		Submit int
		Status int
	}
	var list []*problem
	err = q.All(&list)
	if err != nil {
		http.Error(w, "query error", 599)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"list": list,
	})
	if err != nil {
		http.Error(w, "json error", 599)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func (this *Problem) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)
	if v, ok := args["source"]; ok {
		query["source"] = v
	}
	if v, ok := args["offset"]; ok {
		var offset int
		offset, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		query["pid"] = bson.M{"$gte": offset}
	}
	return
}
