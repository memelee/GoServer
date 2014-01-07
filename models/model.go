package models

import (
	"GoServer/config"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type Result struct {
	Id int
}

type Model struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (this *Model) Init(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
}

func (this *Model) OpenDB() {
	this.Session, _ = mgo.Dial(config.Host)
	this.DB = this.Session.DB(config.DB)
}

func (this *Model) CloseDB() {
	if !config.Lasting {
		this.Session.Close()
	}
}

func (this *Model) GetID(c string) int {
	inc := bson.M{
		"$inc": bson.M{
			"id": 1,
		},
	}
	q := bson.M{
		"name": c,
	}
	cmd := bson.M{
		"findAndModify": "ids",
		"query":         q,
		"update":        inc,
		"upsert":        true,
	}
	one := &Result{}
	_ = this.DB.Run(cmd, one)
	log.Println(one.Id)
	return one.Id
}
