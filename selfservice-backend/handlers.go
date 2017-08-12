package main

import (
	"fmt"

	r "github.com/GoRethink/gorethink"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/zanetworker/son-selfservice/selfservice-backend/models"
)

//StartFSM command to send the SSM to start a specific FSM
func StartFSM(client *Client, fsmInputData interface{}) {
	var fsmData models.FSMAction
	var fsmDataReply models.Message
	if err := mapstructure.Decode(fsmInputData, &fsmData); err != nil {
		log.Error(err)
	}

	//TODO send command to the SSM
	fsmDataReply.Name = "fsm start"
	fsmDataReply.Data = fsmData
	client.send <- fsmDataReply

	log.Infof("%#v\n", fsmData)
	log.Info("started FSM")
}

//StopFSM command to send the SSM to stop a specific FSM
func StopFSM(client *Client, fsmInputData interface{}) {
	var fsmData models.FSMAction
	var fsmDataReply models.Message
	if err := mapstructure.Decode(fsmInputData, &fsmData); err != nil {
		log.Error(err)
	}

	//TODO send command to the SSM
	fsmDataReply.Name = "fsm stop"
	fsmDataReply.Data = fsmData
	client.send <- fsmDataReply

	log.Infof("%#v\n", fsmData)
	log.Info("started FSM")
}

//AddFSM adding an FSM
func AddFSM(client *Client, fsmInputData interface{}) {

	var fsmToInsert models.FSM
	err := mapstructure.Decode(fsmInputData, &fsmToInsert)
	if err != nil {
		log.Error("Failed to map structure FSM")
	} else {
		response, err := r.DB("fsms").Table("fsm_psa").Insert(fsmToInsert).RunWrite(client.db.Connection)
		if err != nil {
			log.Errorf("Failed to inser FSM: %s", err.Error())
			return
		}
		log.Info(response)
	}

}

//UpdateFSM used to update the database with FSMs information received from SSM
func UpdateFSM(client *Client, fsmInputData interface{}) {
	res, err := r.DB("fsms").Table("fsm_psa").Filter(map[string]interface{}{
		"name": "test",
	}).Run(client.db.Connection)
	if err != nil {
		fmt.Print(err)
		return
	}
	// Scan query result into the person variable
	var results []interface{}
	err = res.All(&results)
	var fsmResults models.FSM
	err = mapstructure.Decode(results[0], &fsmResults)
	if err != nil {
		log.Error("Failed to find fsm", err.Error())
		return
	}
	log.Infof("%#v\n", fsmResults.ID)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return
	}
	var updateFSMwithValues models.FSM
	err = mapstructure.Decode(fsmInputData, &updateFSMwithValues)
	if err != nil {
		log.Error("Failed to map fsm to inputdata", err.Error())
	}
	resp, err := r.DB("fsms").Table("fsm_psa").Get(fsmResults.ID).Update(
		fsmInputData).RunWrite(client.db.Connection)
	log.Infof("%#v\n", resp)
	if err != nil {
		fmt.Print(err)
		return
	}
	// log.Infof("%d users", len(results))

}
