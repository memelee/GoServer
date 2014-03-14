package class

import (
	"GoServer/config"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"time"
)

type ids struct {
	Name string `json:"name"bson:"name"`
	Id   int    `json:"id"bson:"id"`
}

type Model struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (this *Model) Init(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
}

func (this *Model) OpenDB() (err error) {
	this.Session, err = mgo.Dial(config.DBHost)
	if err != nil {
		return
	}

	this.DB = this.Session.DB(config.DBName)
	return
}

func (this *Model) CloseDB() {
	if !config.DBLasting {
		this.Session.Close()
	}
}

func (this *Model) GetID(c string) (id int, err error) {
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	var one ids
	_, err = this.DB.C("ids").Find(bson.M{"name": c}).Apply(change, &one)
	id = one.Id
	return
}

func (this *Model) GetTime() (ft string) {
	t := time.Now().Unix()
	ft = time.Unix(t, 0).Format("2006-01-02 15:04:05")
	return
}

func (this *Model) LoadJson(r io.Reader, v interface{}) (err error) {
	err = json.NewDecoder(r).Decode(v)
	return
}

func (this *Model) ParseURL(url string) (args map[string]string) {
	args = make(map[string]string)
	path := strings.Trim(url, "/")
	list := strings.Split(path, "/")

	for i := 1; i < len(list); i += 2 {
		args[list[i-1]] = list[i]
	}
	return
}

func (this *Model) EncryptPassword(str string) (pwd string, err error) {
	m := md5.New()
	_, err = io.WriteString(m, str)
	p1 := fmt.Sprintf("%x", m.Sum(nil))
	s := sha1.New()
	_, err = io.WriteString(s, str)
	p2 := fmt.Sprintf("%x", s.Sum(nil))
	pwd = p1 + p2
	return
}
