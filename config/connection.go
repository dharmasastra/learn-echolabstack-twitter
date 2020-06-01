package config

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"
)

func conn() *mgo.Session {

	// Database connection
	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// Create indices
	if err = db.Copy().DB("twitter").C("users").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	return db
}
