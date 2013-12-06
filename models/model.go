package models

import (
	"GoServer/config"
	"labix.org/v2/mgo"
	"net/http"
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
