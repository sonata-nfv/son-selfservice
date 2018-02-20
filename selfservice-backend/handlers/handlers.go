package handlers

import (
	"fmt"
	"net/http"
	"reflect"

	r "github.com/GoRethink/gorethink"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/zanetworker/son-selfservice/selfservice-backend/communication"
	"github.com/zanetworker/son-selfservice/selfservice-backend/models"
)


// The WebSocket connection to the SSM
var ssmConnection *websocket.Conn

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func SSMConnect(w http.ResponseWriter, r *http.Request) {
     log.Info("Received SSM connection!")

     conn, err := upgrader.Upgrade(w, r, nil)
     if err != nil {
        log.Error(err)
        return
     }
     ssmConnection = conn
}

func sendMessage(request models.Message) (response models.Message) {
    log.Infof("Sending data: %#v", request)
    err := ssmConnection.WriteJSON(request)
    if err != nil {
        log.Error("ERROR on send: ", err)
        return
    }

    log.Info("Reading reply...")
    err = ssmConnection.ReadJSON(&response)
    if err != nil {
        log.Error("ERROR on read: ", err)
        return
    }
    log.Infof("Received: %#v", response)
    return response
}

//StartServiceBasic command to send the SSM for starting the basic tier service
func StartServiceBasic(client *communication.Client, serviceInputData interface{}) {
	log.Info("Starting service BASIC")

	var serviceRequest models.Message
	var serviceReply models.Message

	serviceRequest.Name = "basic"
	serviceReply.Data = ""

	// Send the message to the SSM
	serviceReply = sendMessage(serviceRequest)

	client.Send <- serviceReply
	log.Info("Started service BASIC")
}

//StopServiceBasic command to send the SSM for starting the basic tier service
func StopServiceBasic(client *communication.Client, serviceInputData interface{}) {
	log.Info("Stopping Service BASIC")

	var serviceRequest models.Message
	var serviceReply models.Message

	serviceRequest.Name = "basic stop"
	serviceReply.Data = ""

        // Send the message to the SSM
        serviceReply = sendMessage(serviceRequest)

	client.Send <- serviceReply
	log.Info("Stopped service BASIC")
}

//StartServiceAnon command to send the SSM for starting the Anon tier service
func StartServiceAnon(client *communication.Client, serviceInputData interface{}) {
	log.Info("Starting Service ANON")

	var serviceRequest models.Message
	var serviceReply models.Message

	serviceRequest.Name = "anon"
	serviceReply.Data = ""

        // Send the message to the SSM
        serviceReply = sendMessage(serviceRequest)

	client.Send <- serviceReply
	log.Info("Started service ANON")
}

//StopServiceAnon command to send the SSM for starting the Anon tier service
func StopServiceAnon(client *communication.Client, serviceInputData interface{}) {
	log.Info("Stopping Service ANON")

	var serviceRequest models.Message
	var serviceReply models.Message

	serviceRequest.Name = "anon stop"
	serviceReply.Data = ""

        // Send the message to the SSM
        serviceReply = sendMessage(serviceRequest)

	client.Send <- serviceReply
	log.Info("Stopped service ANON")
}

//StartFSM command to send the SSM to start a specific FSM
func StartFSM(client *communication.Client, fsmInputData interface{}) {
	var fsmData models.FSMAction
	var fsmDataRequest models.Message
	var fsmDataReply models.Message

	if err := mapstructure.Decode(fsmInputData, &fsmData); err != nil {
		log.Error(err)
	}

	fsmDataRequest.Name = "fsm start"
	fsmDataRequest.Data = fsmData

        // Send the message to the SSM
        //fsmDataReply = sendMessage(fsmDataRequest)

	client.Send <- fsmDataReply
	log.Info("Started FSM")
}

//StopFSM command to send the SSM to stop a specific FSM
func StopFSM(client *communication.Client, fsmInputData interface{}) {
	var fsmData models.FSMAction
	var fsmDataRequest models.Message
	var fsmDataReply models.Message

	if err := mapstructure.Decode(fsmInputData, &fsmData); err != nil {
		log.Error(err)
	}

	fsmDataRequest.Name = "fsm stop"
	fsmDataRequest.Data = fsmData

        // Send the message to the SSM
        //fsmDataReply = sendMessage(fsmDataRequest)

	client.Send <- fsmDataReply
	log.Info("Stopped FSM")
}

//AddFSMs adds FSMs if they don't exists
func AddFSMs(client *communication.Client, fsmsInputData interface{}) {
	log.Info("Adding FSMs")
	var fsms models.FSMs

	switch reflect.TypeOf(fsmsInputData).Kind() {
	case reflect.Slice:
		log.Info("Its a slice")
		fsmsInputDataValue := reflect.ValueOf(fsmsInputData)
		// log.Infof("%#v\n", fsmsInputDataValue)

		// if err := mapstructure.Decode(fsmsInputDataValue, &fsms); err != nil {
		// 	log.Error(err.Error())
		// 	return
		// }

		test := fsmsInputDataValue.Interface().([]interface{})
		for _, vlaue := range test {
			var fsm models.FSM
			err := mapstructure.Decode(vlaue, &fsm)
			if err != nil {
				log.Error(err.Error())
			}
			fsms.FSMList = append(fsms.FSMList, fsm)
		}

		for _, fsm := range fsms.FSMList {
			err := AddFSM(client.DB.Connection, fsm)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}

}

//UpdateFSM used to update the database with FSMs information received from SSM
func UpdateFSM(client *communication.Client, fsmInputData interface{}) {
	res, err := r.DB("fsms").Table("fsm_psa").Filter(map[string]interface{}{
		"name": "test",
	}).Run(client.DB.Connection)
	if err != nil {
		fmt.Print(err)
		return
	}
	// Scan query result into the person variable
	var results []interface{}
	err = res.All(&results)
	if len(results) == 0 {
		log.Infof("FSM does not exist in the database")
	}
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
	resp, err := r.DB("fsms").Table("fsm_psa").Get(fsmResults.ID).Update(fsmInputData).RunWrite(client.DB.Connection)
	log.Infof("%#v\n", resp)
	if err != nil {
		fmt.Print(err)
		return
	}

}
