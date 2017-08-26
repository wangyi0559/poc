// ============================================================================================================================
// This smart contract is used for managing the accounts
// It includes adding an account ,deleting an account,changing an
// account's password and verifying the account's password
// ============================================================================================================================
// 本智能合约用于存储医院和医疗分析机构对病人病历的申请状态
// 功能包括：添加记录查询状态
// ============================================================================================================================
package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("account")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type RequestState struct {
	CeleresID  string `json:"celeresID"`  //user who requested for the medical record
	HospitalID string `json:"hospitalID"` //the ID of the hospital which stores this medical record
	Hcode      string `json:"hcode"`      //the authorization code of the hospital
	PatientID  string `json:"patientID"`  // the ID of the patient who owns the medical record
	Pcode      string `json:"pcode"`      //the authorization code of the patient
}

// ============================================================================================================================
// the Init function is used for deploying the chaincode and setting the Administrator account
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### requestRecord Init ###########")

	_, args := stub.GetFunctionAndParameters()
	var ID string    // Administrator's ID
	var IDval string // Password of the Administrator ID
	var err error

	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments. "))
	}
	// Initialize the chaincode
	ID = args[0]
	IDval = args[1]

	// Write the password to the ledger
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

	logger.Info("########### requestRecord Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "add" {
		return t.Add(stub, args)
	}
	if function == "verify" {
		return t.VerifyQuery(stub, args)
	}
	if function == "test" { //add a new Administrator account
		return t.Test(stub, args)
	}

	return shim.Success([]byte("Received unknown function invocation"))
}

// ============================================================================================================================
// Add function is used for add an new request record of the hospital or the celeres
// 6 input
// "requestID","requester","hospitalID","hospital authorization code","patient","hpatient authorization code"
// ============================================================================================================================
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	requestID := args[0]
	var requestState RequestState
	requestState.CeleresID = args[1]
	requestState.HospitalID = args[2]
	requestState.Hcode = args[3]
	requestState.PatientID = args[4]
	requestState.Pcode = args[5]
	JsonRequestState, _ := json.Marshal(requestState)

	// add the account
	err := stub.PutState(requestID, JsonRequestState)
	if err != nil {
		return shim.Success([]byte("Failed to add the record"))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// VerifyQuery function is used for replying the state of the medical record request
// 1 input
// "requestID"
// ============================================================================================================================
func (t *SimpleChaincode) VerifyQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		shim.Success([]byte("Incorrect number of arguments. "))
	}
	requestID := args[0]
	result, err := stub.GetState(requestID)
	//test if the account has been existed
	if err != nil || result == nil {
		return shim.Success([]byte("error in reading request record"))
	}

	return shim.Success([]byte(result))
}

func (t *SimpleChaincode) Test(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	account := args[0]
	// Handle different functions
	password, err := stub.GetState(account) //get the var from chaincode state
	if err != nil || password == nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + account + "\"}"
		return shim.Success([]byte(jsonResp))
	}

	return shim.Success([]byte(password))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
