package dao

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PearlsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "pearl"
)

func (m *PearlsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *PearlsDAO) FindAll() ([]Pearl, error) {
	var pearls []Pearl
	err := db.C(COLLECTION).Find(bson.M{}).All(&pearls)
	return pearls, err
}

// Find a Pearl by its id
func (m *PearlsDAO) FindById(id string) (Pearl, error) {
	var pearl Pearl
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&pearl)
	return pearl, err
}

func (m *PearlsDAO) Insert(Pearl pearl) error {
	err := db.C(COLLECTION).Insert(&pearl)
	return err
}

func (m *PearlsDAO) Delete(Pearl pearl) error {
	err := db.C(COLLECTION).Remove(&pearl)
	return err
}

func (m *PearlsDAO) Update(Pearl pearl) error {
	err := db.C(COLLECTION).UpdateId(Pearl.ID, &pearl)
	return err
}
