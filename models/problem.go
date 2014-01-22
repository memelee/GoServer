package models

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type Problem struct {
	Model
}

func (this *Problem) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Insert")
	this.Init(w, r)

	this.OpenDB()
	defer this.CloseDB()
	c := this.DB.C("problem")

	pid := this.GetID("problem")
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
	create_time := this.GetTime()

	err := c.Insert(bson.M{
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
		"create_time":   create_time,
	})
	if err != nil {
		log.Println(err)
	} else {
		b, _ := json.Marshal(map[string]interface{}{
			"pid":    pid,
			"ok":     1,
			"status": 1,
		})
		w.Write(b)
	}
}
