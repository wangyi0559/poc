// ============================================================================================================================
// This smart contract is used for storing and verifying the original medical record
// It includes creating the Administrator account ,adding the original medical record
// and verifying the aoriginal medical record
// ============================================================================================================================

package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("account")

type SimpleChaincode struct {
}

type RecordDetail struct {
	DbIP       string `json:"dbIP"`       //user who created the open trade order
	RecordHash string `json:"recordHash"` //
}

type RecordDiseaseIndexDetail struct {
	RecordID string `json:"recordID"`
	DbIP     string `json:"dbIP"` //user who created the open trade order

}

type RecordDiseaseIndex struct {
	RecordDiseaseIndexDetails []RecordDiseaseIndexDetail `json:"recordDiseaseIndexDetails"` //user who created the open trade order
}

// ============================================================================================================================
// Init function is used for creating the Administrator account
// ============================================================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### account Init ###########")

	_, args := stub.GetFunctionAndParameters()

	var ID string    // Administrator's ID
	var IDval string // Password of the Administrator ID
	var err error

	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments. "))
	}
	ID = args[0]
	IDval = args[1]
	// Write the state to the ledger
	err = stub.PutState(ID, []byte(IDval))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke function is the entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	logger.Info("########### account Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	// Handle different functions
	if function == "add" { //add a new account
		return t.Add(stub, args)
	}
	if function == "query" { //deletes an account from its state
		return t.Get(stub, args)
	}
	if function == "verify" {
		return t.VerifyRecordHash(stub, args)

	}

	return shim.Success([]byte("Received unknown function invocation"))
}

// ============================================================================================================================
// Add function is used for add an new medical record
// 4 input
// "medical record ID","medical record db address","medical record hash"
// ============================================================================================================================
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	// write the record details into the world state and the chain
	recordID := args[0]
	var recordDetail RecordDetail
	recordDetail.DbIP = args[1]
	recordDetail.RecordHash = args[2]
	JsonRecordDetail, _ := json.Marshal(recordDetail)
	recordTest, err := stub.GetState(recordID)

	//test if the account has been existed
	if recordTest != nil {
		return shim.Success([]byte("the record is existed"))
	}
	// add the hash
	err = stub.PutState(recordID, JsonRecordDetail)
	if err != nil {
		return shim.Success([]byte("Failed to add the record"))
	}

	// write the disease index of the record details into the world state and the chain
	disease := args[3]
	var recordDiseaseIndexDetail RecordDiseaseIndexDetail
	var recordDiseaseIndexDetails RecordDiseaseIndex
	recordDiseaseIndexDetail.RecordID = args[0]
	recordDiseaseIndexDetail.DbIP = args[1]
	recordDiseaseIndexAsBytes, errs := stub.GetState(disease)
	if errs != nil {
		return shim.Success([]byte("Failed to get disease index details"))
	}
	json.Unmarshal(recordDiseaseIndexAsBytes, &recordDiseaseIndexDetails)
	recordDiseaseIndexDetails.RecordDiseaseIndexDetails = append(recordDiseaseIndexDetails.RecordDiseaseIndexDetails, recordDiseaseIndexDetail)
	jsonAsBytes, _ := json.Marshal(recordDiseaseIndexDetails)
	err = stub.PutState(disease, jsonAsBytes)
	return shim.Success(nil)
}

// ============================================================================================================================
// Get function is used for getting the medical record details
// 1 input
// "medical record ID"
// ============================================================================================================================
func (t *SimpleChaincode) Get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	recordID := args[0]
	// Handle different functions
	recordDetail, err := stub.GetState(recordID) //get the var from chaincode state
	if err != nil || recordDetail == nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + recordID + "\"}"
		return shim.Success([]byte(jsonResp))
	}
	return shim.Success([]byte(recordDetail))
}

// ============================================================================================================================
// VerifyRecordHash function is used for verifying the medical record in this smart contract
// 2 input
// "medical record ID","medical record hash"
// ============================================================================================================================
func (t *SimpleChaincode) VerifyRecordHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments. "))
	}
	recordID := args[0]
	recordHash := args[1]
	recordTest, err := stub.GetState(recordID)
	var JsonRecordTest RecordDetail
	json.Unmarshal(recordTest, &JsonRecordTest)
	ver := []byte("ok")
	jsonResp := "{\"Error\":\"Failed to get state for " + recordID + "\"}"
	//test if the account has been existed
	if err != nil || recordTest == nil {
		return shim.Success([]byte(jsonResp))
	}

	// verify
	if recordHash == string(JsonRecordTest.RecordHash) {
		return shim.Success(ver)

	} else {
		return shim.Success([]byte("Failed to verify the record"))
	}

}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
