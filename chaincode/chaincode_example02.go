package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "reflect"
	"time"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	Offers []Offer
}

const CONTRACTSTATEKEY string = "ContractStateKey"  
// store contract state - only version in this example
const MYVERSION string = "1.0"

// ************************************
// offer and contract state 
// ************************************

type Offer struct {
	OfferID        		string       `json:"offerID,omitempty"`        	 // all offers must have an ID, primary key of contract
    Startlocation       Geolocation  `json:"startlocation,omitempty"`       // start location
	Endlocation         Geolocation  `json:"endlocation,omitempty"`         // end location
	startLocationStr    string		  `json:"startlocationstr,omitempty"`    // start location string
	endLocationStr      string		  `json:"endlocationstr,omitempty"`      // end location string
	startTime	   		time	      `json:"starttime,omitempty"`		     // earliest start of asset
	endTime		   		time         `json:"endtime,omitempty"`	         // latest end of asset
	isBid		   		boolean      `json:"isbid,omitempty"`               // is bid offer or not
    Carrier        		[]string       `json:"carrier,omitempty"`             // list of carrier (if assigned size = 1)
	Owner        		string       `json:"owner,omitempty"`               // the owner
}

type ContractState struct {
    Version string `json:"version"`
}

type Geolocation struct {
    Latitude    *float64 `json:"latitude,omitempty"`
    Longitude   *float64 `json:"longitude,omitempty"`
}

var contractState = ContractState{MYVERSION}


// ************************************
// deploy callback mode 
// ************************************
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    var stateArg ContractState
    var err error
    if len(args) != 1 {
        return nil, errors.New("init expects one argument, a JSON string with tagged version string")
    }
    err = json.Unmarshal([]byte(args[0]), &stateArg)
    if err != nil {
        return nil, errors.New("Version argument unmarshal failed: " + fmt.Sprint(err))
    }
    if stateArg.Version != MYVERSION {
        return nil, errors.New("Contract version " + MYVERSION + " must match version argument: " + stateArg.Version)
    }
    contractStateJSON, err := json.Marshal(stateArg)
    if err != nil {
        return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
    }
    err = stub.PutState(CONTRACTSTATEKEY, contractStateJSON)
    if err != nil {
        return nil, errors.New("Contract state failed PUT to ledger: " + fmt.Sprint(err))
    }
    return nil, nil
}

// ************************************
// deploy and invoke callback mode 
// ************************************
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    // Handle different functions
    if function == "createOffer" {
        // create OfferID
        return t.createOffer(stub, args)
    } else if function == "updateOffer" {
        // create OfferID
        return t.updateOffer(stub, args)
    } else if function == "deleteOffer" {
        // Deletes an Offer by ID from the ledger
        return t.deleteOffer(stub, args)
    }
    return nil, errors.New("Received unknown invocation: " + function)
}

// ************************************
// query callback mode 
// ************************************
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    // Handle different functions
    if function == "readOffer" {
        // gets the state for an OfferID as a JSON struct
        return t.readOffer(stub, args)
    } else if function =="readOfferObjectModel" {
        return t.readOfferObjectModel(stub, args)
    }  else if function == "readOfferSamples" {
		// returns selected sample objects 
		return t.readOfferSamples(stub, args)
	} else if function == "readOfferSchemas" {
		// returns selected sample objects 
		return t.readOfferSchemas(stub, args)
	}
    return nil, errors.New("Received unknown invocation: " + function)
}

/**********main implementation *************/

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple Chaincode: %s", err)
    }
}

/*****************Offer CRUD INTERFACE starts here************/

/****************** 'deploy' methods *****************/

/******************** createOffer ********************/

func (t *SimpleChaincode) createOffer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    _,erval:=t. createOrUpdateOffer(stub, args)
    return nil, erval
}

//******************** updateOffer ********************/

func (t *SimpleChaincode) updateOffer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
     _,erval:=t. createOrUpdateOffer(stub, args)
    return nil, erval
}


//******************** deleteOffer ********************/

func (t *SimpleChaincode) deleteOffer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var OfferID string // Offer ID
    var err error
    var stateIn Offer

    // validate input data for number of args, Unmarshaling to Offer state and obtain Offer id
    stateIn, err = t.validateInput(args)
    if err != nil {
        return nil, err
    }
    OfferID = *stateIn.OfferID
    // Delete the key / Offer from the ledger
    err = stub.DelState(OfferID)
    if err != nil {
        err = errors.New("DELSTATE failed! : "+ fmt.Sprint(err))
       return nil, err
    }
    return nil, nil
}

/******************* Query Methods ***************/

//********************readOffer********************/

func (t *SimpleChaincode) readOffer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var OfferID string // Offer ID
    var err error
    var state Offer

     // validate input data for number of args, Unmarshaling to Offer state and obtain Offer id
    stateIn, err:= t.validateInput(args)
    if err != nil {
        return nil, errors.New("Offer does not exist!")
    }
    OfferID = *stateIn.OfferID
        // Get the state from the ledger
    OfferBytes, err:= stub.GetState(OfferID)
    if err != nil  || len(OfferBytes) ==0{
        err = errors.New("Unable to get Offer state from ledger")
        return nil, err
    } 
    err = json.Unmarshal(OfferBytes, &state)
    if err != nil {
         err = errors.New("Unable to unmarshal state data obtained from ledger")
        return nil, err
    }
    return OfferBytes, nil
}

//*************readOfferObjectModel*****************/

func (t *SimpleChaincode) readOfferObjectModel(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var state Offer = Offer{}

    // Marshal and return
    stateJSON, err := json.Marshal(state)
    if err != nil {
        return nil, err
    }
    return stateJSON, nil
}
//*************readOfferSamples*******************/

func (t *SimpleChaincode) readOfferSamples(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte(samples), nil
}
//*************readOfferSchemas*******************/

func (t *SimpleChaincode) readOfferSchemas(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte(schemas), nil
}

// ************************************
// validate input data : common method called by the CRUD functions
// ************************************
func (t *SimpleChaincode) validateInput(args []string) (stateIn Offer, err error) {
    var OfferID string // Offer ID
    var state Offer = Offer{} // The calling function is expecting an object of type Offer

    if len(args) !=1 {
        err = errors.New("Incorrect number of arguments. Expecting a JSON strings with mandatory OfferID")
        return state, err
    }
    jsonData:=args[0]
    OfferID = ""
    stateJSON := []byte(jsonData)
    err = json.Unmarshal(stateJSON, &stateIn)
    if err != nil {
        err = errors.New("Unable to unmarshal input JSON data")
        return state, err
        // state is an empty instance of Offer state
    }      
    // was OfferID present?
    // The nil check is required because the Offer id is a pointer. 
    // If no value comes in from the json input string, the values are set to nil
    
    if stateIn.OfferID !=nil { 
        OfferID = strings.TrimSpace(*stateIn.OfferID)
        if OfferID==""{
            err = errors.New("OfferID not passed")
            return state, err
        }
    } else {
        err = errors.New("Offer id is mandatory in the input JSON data")
        return state, err
    }
    
    
    stateIn.OfferID = &OfferID
    return stateIn, nil
}
//******************** createOrUpdateOffer ********************/

func (t *SimpleChaincode) createOrUpdateOffer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var OfferID string                 // Offer ID                    // used when looking in map
    var err error
    var stateIn Offer
    var stateStub Offer
   

    // validate input data for number of args, Unmarshaling to Offer state and obtain Offer id

    stateIn, err = t.validateInput(args)
    if err != nil {
        return nil, err
    }
    OfferID = *stateIn.OfferID
    // Partial updates introduced here
    // Check if Offer record existed in stub
    OfferBytes, err:= stub.GetState(OfferID)
    if err != nil || len(OfferBytes)==0{
        // This implies that this is a 'create' scenario
         stateStub = stateIn // The record that goes into the stub is the one that cme in
    } else {
        // This is an update scenario
        err = json.Unmarshal(OfferBytes, &stateStub)
        if err != nil {
            err = errors.New("Unable to unmarshal JSON data from stub")
            return nil, err
            // state is an empty instance of Offer state
        }
          // Merge partial state updates
        stateStub, err =t.mergePartialState(stateStub,stateIn)
        if err != nil {
            err = errors.New("Unable to merge state")
            return nil,err
        }
    }
    stateJSON, err := json.Marshal(stateStub)
    if err != nil {
        return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
    }
    // Get existing state from the stub
    
  
    // Write the new state to the ledger
    err = stub.PutState(OfferID, stateJSON)
    if err != nil {
        err = errors.New("PUT ledger state failed: "+ fmt.Sprint(err))            
        return nil, err
    } 
	t.Offers = append(t.Offers, OfferID)
    return nil, nil
}
/*********************************  internal: mergePartialState ****************************/	
 func (t *SimpleChaincode) mergePartialState(oldState Offer, newState Offer) (Offer,  error) {
     
    old := reflect.ValueOf(&oldState).Elem()
    new := reflect.ValueOf(&newState).Elem()
    for i := 0; i < old.NumField(); i++ {
        oldOne:=old.Field(i)
        newOne:=new.Field(i)
        if ! reflect.ValueOf(newOne.Interface()).IsNil() {
            oldOne.Set(reflect.Value(newOne))
        } 
    }
    return oldState, nil
 }