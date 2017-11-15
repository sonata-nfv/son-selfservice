package database

import (
	"fmt"

	r "github.com/GoRethink/gorethink"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/zanetworker/son-selfservice/selfservice-backend/models"
	"time"
)

//Database is the struct for connecting to the database
type Database struct {
	Connection *r.Session
	DBChannel  chan r.ChangeResponse
}

//DeleteAll deletes all entries in the datq
func (db *Database) DeleteAll(dbName, tableName string) {
	r.DB(dbName).Table("fsm_psa").Delete().Run(db.Connection)
}

//AddFSM asdasd
func (db *Database) AddFSM(dbName, tableName string, dataToAdd interface{}) error {
	var fsmToAdd models.FSM
	if err := mapstructure.Decode(dataToAdd, &fsmToAdd); err != nil {
		log.Error(err.Error())
		return err
	}
	_, err := r.DB(dbName).Table(tableName).Insert(dataToAdd).RunWrite(db.Connection)
	if err != nil {
		return err
	}

	return nil
}

//NewDB initialize DB
func NewDB() *Database {
     	var err error
	var sess *r.Session

	for i := 1; ; i++ {
		//dbURL := configuration.AppConfig.DBHostIP + ":" + configuration.AppConfig.DBHostPort
		session, err := r.Connect(r.ConnectOpts{
			 //TODO fetch from Config File
			 Address: "selfservice-db:28015",
			 //Address:  "localhost:28015",
			 Database: "fsms",
		})
		if err == nil {
		   sess = session
		   break
	   	}
		if i >= 10 {
		   log.Fatal(err.Error())
		}
		time.Sleep(2*time.Second)
	}
	
	resp, err := r.DBCreate("fsms").RunWrite(sess)
	if err != nil {
		fmt.Print(err)
	}
	log.Infof("%d DB created", resp.DBsCreated)

	response, err := r.DB("fsms").TableCreate("fsm_psa").RunWrite(sess)
	if err != nil {
		log.Errorf("Error creating table: %s", err)
	}

	log.Infof("%d table created", response.TablesCreated)

	return &Database{
		Connection: sess,
		DBChannel:  make(chan r.ChangeResponse),
	}
}
