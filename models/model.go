package models

import (
	"GoServer/config"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"time"
)

type Model struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (this *Model) Init(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
}

func (this *Model) OpenDB() (err error) {
	this.Session, err = mgo.Dial(config.Host)
	if err != nil {
		return
	}
	this.DB = this.Session.DB(config.DB)
	return
}

func (this *Model) CloseDB() {
	if !config.Lasting {
		this.Session.Close()
	}
}

func (this *Model) GetID(c string) (id int, err error) {
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
	err = this.DB.Run(cmd, one)
	id = one.Value.Id
	return
}

func (this *Model) GetTime() string {
	t := time.Now().Unix()
	ft := time.Unix(t, 0).Format("2006-01-02 15:04:05")
	return ft
}

func (this *Model) ParseURL(url string) map[string]string {
	args := make(map[string]string)
	path := strings.Trim(url, "/")
	list := strings.Split(path, "/")

	for i := 1; i < len(list); i += 2 {
		args[list[i-1]] = list[i]
	}
	return args
}
