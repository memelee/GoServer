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
var pListSelector = bson.M{"_id": 0, "pid": 1, "title": 1, "source": 1}

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
	_, err := this.CheckQuery(args)
	if err != nil {
		http.Error(w, "args error", 400)
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 599)
		return
	}

	// c := this.DB.C("problem")

}

func (this *Problem) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)
	if args["source"] != "" {
		query["source"] = args["source"]
	}
	if args["offset"] != "" {
		var offset int
		offset, err = strconv.Atoi(args["offset"])
		if err != nil {
			return
		}
		query["pid"] = bson.M{"$gte": offset}
	}
	if args["limit"] != "" {
		_, err = strconv.Atoi(args["limit"])
		if err != nil {
			return
		}
	}
	return
}
