import {UPDATE_FSMS_BASIC, UPDATE_FSMS_ANON} from '../actions'
//e.g., state
// {
//     id: "12345",
//     name: "Firewall",
//     state: "stopped"
// }
const INIT_STATE = {
  basic:[
    {"name": "VPN",   "id": "1", "state": "started"},
    {"name": "TOR",   "id": "3", "state": "stopped"},
    {"name": "FW",    "id": "4", "state": "started"},
  ],
  anon:[
    {"name": "VPN",   "id": "1", "state": "started"},
    {"name": "Proxy", "id": "2", "state": "stopped"},
    {"name": "TOR",   "id": "3", "state": "stopped"},
    {"name": "FW",    "id": "4", "state": "started"}
  ]
};



export default (state = INIT_STATE, action) => {
  switch (action.type){
    case UPDATE_FSMS_BASIC:
    console.log("[reducer] update fsms basic" +action.payload)
    console.log(action.payload)
      const stateToReturnBasic = {
        basic:
          state.basic.filter((fsm)=>{
            return fsm.id !== action.payload.id
      }),
	anon:
          state.anon.filter((fsm)=>{
            return fsm.id !== action.payload.id
      }),
    }
      if (action.payload.id != "2") {
	      stateToReturnBasic.basic.push(createStateObjectBasic(action.payload))
	}
      stateToReturnBasic.anon.push(createStateObjectAnon(action.payload))
      const stateToReturnBasicOrdered = {
        basic: stateToReturnBasic.basic.sort(function(a, b){
          return a.id-b.id
        }),
        anon: stateToReturnBasic.anon.sort(function(a, b){
          return a.id-b.id
        })
        }
      console.log(stateToReturnBasicOrdered)
      return stateToReturnBasicOrdered

    case UPDATE_FSMS_ANON:
    console.log("[reducer] update fsms anon" +action.payload)
    console.log(action.payload)
      const stateToReturnAnon ={
        anon:
        state.anon.filter((fsm)=>{
          return fsm.id !== action.payload.id
      }),
        basic:
          state.basic.filter((fsm)=>{
            return fsm.id !== action.payload.id
      })
    }
      stateToReturnAnon.anon.push(createStateObjectAnon(action.payload))
      if (action.payload.id != "2") {
	      stateToReturnAnon.basic.push(createStateObjectBasic(action.payload));
      }

      //order the state
      const stateToReturnAnonOrdered ={
        anon: stateToReturnAnon.anon.sort(function(a, b){
            return a.id-b.id
          }),
        basic: stateToReturnAnon.basic.sort(function(a, b){
          return a.id-b.id
        })
      }
      console.log(stateToReturnAnonOrdered)
      return stateToReturnAnonOrdered

    default:
      return state;
  }
}


const createStateObjectBasic= (dataForObject) => {
  return {
      id: dataForObject.id,
      name: dataForObject.name,
      state: dataForObject.state
  }
}

const createStateObjectAnon= (dataForObject) => {
  return {
      id: dataForObject.id,
      name: dataForObject.name,
      state: dataForObject.state
  }
}
