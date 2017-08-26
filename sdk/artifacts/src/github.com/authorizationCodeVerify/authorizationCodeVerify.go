/*
This Chaincode stores the accounts and passwords

*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("authorizationCodeVerify")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Init function is used for creating the Administrator account
// ============================================================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### authorizationCodeVerify Init ###########")

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

	logger.Info("########### authorizationCodeVerify Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	// Handle different functions
	if function == "add" { //add a new account
		return t.Add(stub, args)
	}
	if function == "verify" {
		if len(args) == 4 { //verify the request of the  hespital
			return t.P2VerifyQuery(stub, args)
		} else if len(args) == 5 { //verify the request of the  celeres
			return t.P3VerifyQuery(stub, args)
		}
	}
	if function == "test" { //add a new Administrator account
		return t.Test(stub, args)
	}

	return shim.Success([]byte("Received unknown function invocation"))
}

// ============================================================================================================================
// Add function is used for adding an new Authorization code
// 2 input
// "account","authorization code"
// ============================================================================================================================
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	account := args[0]
	code := args[1]

	err := stub.PutState(account, []byte(code))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// P2VerifyQuery function is used for verifying the hospital's authorization code and the patient's authorization code
// 4 input
// "hospitalID","hospital authorization code","patientID","patient authorization code",
// ============================================================================================================================
func (t *SimpleChaincode) P2VerifyQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	p2 := args[0] //hospital
	p2Code := args[1]
	p1 := args[2] //patient
	p1Code := args[3]
	p2Test, err := stub.GetState(p2)
	p1Test, errs := stub.GetState(p1)
	//test if the account has been existed
	if err != nil || p2Test == nil {
		return shim.Success([]byte("error in reading hospital's code"))
	}
	if errs != nil || p1Test == nil {
		return shim.Success([]byte("error in reading patient's code"))
	}
	if p2Code != string(p2Test) {
		return shim.Success([]byte("hospital's code is error"))
	} else if p1Code != string(p1Test) {
		return shim.Success([]byte("patient's code is error"))
	}
	return shim.Success([]byte("ok"))

}

// ============================================================================================================================
// P2VerifyQuery function is used for verifying the hospital's authorization code and the patient's authorization code
// 5 input
// "requester ID","hospitalID","hospital authorization code","patientID","patient authorization code",
// ============================================================================================================================
func (t *SimpleChaincode) P3VerifyQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//p3 := args[0]
	p2 := args[1]
	p2Code := args[2]
	p1 := args[3]
	p1Code := args[4]
	p2Test, err := stub.GetState(p2)
	p1Test, errs := stub.GetState(p1)
	//test if the account has been existed
	if err != nil || p2Test == nil {
		return shim.Success([]byte("error in reading hospital's code"))
	}
	if errs != nil || p1Test == nil {
		return shim.Success([]byte("error in reading patient's code"))
	}
	if p2Code != string(p2Test) {
		return shim.Success([]byte("hospital's code is error"))
	} else if p1Code != string(p1Test) {
		return shim.Success([]byte("patient's code is error"))
	}
	return shim.Success([]byte("ok"))

}

func (t *SimpleChaincode) Test(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Println("query is running " + function)
	account := args[0]
	// Handle different functions
	code, err := stub.GetState(account) //get the var from chaincode state

	if err != nil || code == nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + account + "\"}"
		return shim.Success([]byte(jsonResp))
	}

	return shim.Success([]byte(code))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
