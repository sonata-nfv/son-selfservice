import {UPDATE_FSMS_BASIC, UPDATE_FSMS_ANON} from '../actions'
//e.g., state
// {
//     id: "12345",
//     name: "Firewall",
//     state: "stopped"
// }
const INIT_STATE = {
  basic:[
    {"name": "VPN",           "id": "1", "state": "stopped"},
    {"name": "TOR",           "id": "2", "state": "stopped"},
  ],
  anon:[
    {"name": "VPN",           "id": "1", "state": "stopped"},
    {"name": "Proxy",         "id": "2", "state": "stopped"},
    {"name": "TOR",           "id": "3", "state": "stopped"}
  ]
};



export default (state = INIT_STATE, action) => {
  switch (action.type){
    case UPDATE_FSMS_BASIC:
    console.log("[reducer] update fsms basic" +action.payload)
      const stateToReturnBasic = {
        basic:
          state.basic.filter((fsm)=>{
            return fsm.id !== action.payload.id
      })
    }
      stateToReturnBasic.basic.push(createStateObjectBasic(action.payload))
      const stateToReturnBasicOrdered = {
        basic: stateToReturnBasic.basic.sort(function(a, b){
          return a.id-b.id
        }),
        anon: INIT_STATE.anon
        }
      return stateToReturnBasicOrdered

    case UPDATE_FSMS_ANON:
      const stateToReturnAnon ={
        anon:
        state.anon.filter((fsm)=>{
          return fsm.id !== action.payload.id
      })
    }
      stateToReturnAnon.anon.push(createStateObjectAnon(action.payload))

      //order the state
      const stateToReturnAnonOrdered ={
        anon: stateToReturnAnon.anon.sort(function(a, b){
            return a.id-b.id
          }),
          basic: INIT_STATE.basic
      }
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
