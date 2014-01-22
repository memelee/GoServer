package models

import (
	"GoServer/config"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

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
	type result struct {
		Value struct {
			Id int
		}
	}

	cmd := bson.M{
		"findAndModify": "ids",
		"query": bson.M{
			"name": c,
		},
		"update": bson.M{
			"$inc": bson.M{
				"id": 1,
			},
		},
		"upsert": true,
	}
	one := &result{}
	_ = this.DB.Run(cmd, one)
	return one.Value.Id
}

func (this *Model) GetTime() string {
	t := time.Now().Unix()
	ft := time.Unix(t, 0).Format("2006-01-02 15:04:05")
	return ft
}
