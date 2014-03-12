package model

import (
	"GoServer/class"
	"GoServer/config"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

type user struct {
	Uid string `json:"uid"bson:"uid"`
	Pwd string `json:"pwd"bson:"pwd"`

	Nick   string `json:"nick"bson:"nick"`
	Mail   string `json:"mail"bson:"mail"`
	School string `json:"school"bson:"school"`

	Privilege int `json:"privilege"bson:"privilege"`

	Solve  int `json:"solve"bson:"solve"`
	Submit int `json:"submit"bson:"submit"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:'create'`
}

type User struct {
	class.Model
}

var uDetailSelector = bson.M{"_id": 0}
var uListSelector = bson.M{"_id": 0, "uid": 1, "nick": 1, "solve": 1, "submit": 1, "status": 1}

// POST /user/login
func (this *User) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Login")
	this.Init(w, r)

	var ori user
	err := this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	ori.Pwd, err = class.EncryptPassword(ori.Pwd)
	if err != nil {
		http.Error(w, "encrypt error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	var alt user
	err = this.DB.C("user").Find(bson.M{"uid": ori.Uid}).Select(uDetailSelector).One(&alt)
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "login error", 500)
		return
	}

	var b []byte
	if ori.Pwd == alt.Pwd {
		log.Println("Server User Login Successfully")
		b, err = json.Marshal(map[string]interface{}{
			"uid":       alt.Uid,
			"privilege": alt.Privilege,
			"status":    alt.Status,
		})
	} else {
		log.Println("Server User Login Failed")
		b, err = json.Marshal(map[string]interface{}{
			"uid":       "",
			"privilege": 0,
			"status":    0,
		})
	}
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /user/logout
func (this *User) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Logout")
	this.Init(w, r)

	var one user
	err := this.LoadJson(r.Body, &one)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	w.WriteHeader(200)
}

// POST /user/detail/uid/<uid>
func (this *User) Detail(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem Detail")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid := args["uid"]

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	var one user
	err = this.DB.C("user").Find(bson.M{"uid": uid}).Select(uDetailSelector).One(&one)
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

// POST /user/delete/uid/<uid>
func (this *User) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Delete")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid := args["uid"]

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("user").Remove(bson.M{"uid": uid})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "delete error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /user/insert
func (this *User) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Insert")
	this.Init(w, r)

	var one user
	err := this.LoadJson(r.Body, &one)
	if err != nil {
		http.Error(w, "load errpr", 400)
		return
	}

	one.Pwd, err = class.EncryptPassword(one.Pwd)
	if err != nil {
		http.Error(w, "encrypt error", 400)
		return
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	one.Privilege = config.PrivilegePU
	one.Solve = 0
	one.Submit = 0
	one.Status = 1
	one.Create = this.GetTime()

	err = this.DB.C("user").Insert(&one)
	if err != nil {
		http.Error(w, "insert error", 500)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"uid":       one.Uid,
		"privilege": one.Privilege,
		"status":    one.Status,
	})
	if err != nil {
		http.Error(w, "json error", 500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

// POST /user/update/uid/<uid>
func (this *User) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Update")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid := args["uid"]

	var ori user
	err := this.LoadJson(r.Body, &ori)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	alt := make(map[string]interface{})
	if ori.Pwd != "" {
		alt["pwd"] = ori.Pwd
	}
	if ori.Nick != "" {
		alt["nick"] = ori.Nick
	}
	if ori.Mail != "" {
		alt["mail"] = ori.Mail
	}
	if ori.School != "" {
		alt["school"] = ori.School
	}
	if ori.Privilege > config.PrivilegeNA {
		alt["privilege"] = ori.Privilege
	}

	err = this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("user").Update(bson.M{"uid": uid}, bson.M{"$set": alt})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "update error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /user/status/uid/<uid>
func (this *User) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Status")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid := args["uid"]

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("user").Update(bson.M{"uid": uid}, bson.M{"$inc": bson.M{"status": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "status error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /user/record/uid/<uid>/action/<solve/submit>
func (this *User) Record(w http.ResponseWriter, r *http.Request) {
	log.Println("Server User Submit")
	this.Init(w, r)

	args := this.ParseURL(r.URL.Path[2:])
	uid := args["uid"]

	var inc int
	switch v := args["action"]; v {
	case "solve":
		inc = 1
	case "submit":
		inc = 0
	default:
		http.Error(w, "args error", 400)
		return
	}

	err := this.OpenDB()
	defer this.CloseDB()
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	err = this.DB.C("user").Update(bson.M{"uid": uid}, bson.M{"$inc": bson.M{"solve": inc, "submit": 1}})
	if err == mgo.ErrNotFound {
		http.Error(w, "not found", 404)
		return
	} else if err != nil {
		http.Error(w, "record error", 500)
		return
	}

	w.WriteHeader(200)
}

// POST /user/list/offset/<offset>/limit/<limit>/uid/<uid>/nick/<nick>
func (this *User) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Server Problem List")
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

	q := this.DB.C("user").Find(query).Select(uListSelector).Sort("-solve", "submit")

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

	var list []*user
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

func (this *User) CheckQuery(args map[string]string) (query bson.M, err error) {
	query = make(bson.M)

	if v, ok := args["uid"]; ok {
		query["uid"] = v
	}
	if v, ok := args["nick"]; ok {
		query["nick"] = bson.M{"$regex": bson.RegEx{v, "i"}}
	}
	return
}
