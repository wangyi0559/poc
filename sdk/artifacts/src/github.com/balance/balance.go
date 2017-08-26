// ============================================================================================================================
// This smart contract is used for managing the accounts
// It includes adding an account ,deleting an account,changing an
// account's password and verifying the account's password
// ============================================================================================================================
// 本智能合约用于用户余额管理
// 功能包括：增加账号 删除账号 账户转账
// ============================================================================================================================

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

var logger = shim.NewLogger("balance")

type SimpleChaincode struct {
}

// ============================================================================================================================
// Init function is used for creating the Administrator account
// ============================================================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### balance Init ###########")

	_, args := stub.GetFunctionAndParameters()

	var ID string    // Administrator's ID
	var IDval string // balance of the Administrator ID
	var err error

	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments. "))
	}

	// Initialize the chaincode
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

	logger.Info("########### balance Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	// Handle different functions
	if function == "delete" { //deletes an account from its state
		return t.Delete(stub, args)
	}
	if function == "add" { //add a new account
		return t.Add(stub, args)
	}
	if function == "transfer" { //change the password of the account
		return t.Transfer(stub, args)
	}
	if function == "ifAccountExisted" { // reply if the account is existed
		return t.IfAccountExisted(stub, args)
	}
	if function == "test" { //a function for testing
		return t.Test(stub, args)
	}

	return shim.Success([]byte("Received unknown function invocation"))
}

// ============================================================================================================================
// Delete function is used for deleting an account
// 1 input
// "account"
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	account := args[0]
	err := stub.DelState(account) //remove the key from chaincode state
	if err != nil {
		return shim.Success([]byte("Failed to delete account"))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Add function is used for adding a new accounts
// 2 input
// "account","balance"
// ============================================================================================================================
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	account := args[0]
	password := args[1]
	accountTest, err := stub.GetState(account)

	//test if the account has been existed
	if err != nil {
		return shim.Success([]byte("Failed to get state"))
	}
	if accountTest != nil {
		return shim.Success([]byte("the ccount is existed"))
	}
	// add the account
	err = stub.PutState(account, []byte(password))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Edit function is used for transfer money from account1 to account2
// 3 input
// "account","old password","new password"
// ============================================================================================================================
func (t *SimpleChaincode) Transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	var A, B string
	var Aval, Bval int
	var X int
	A = args[0]
	B = args[1]

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Success([]byte("Failed to get state"))
	}
	if Avalbytes == nil {
		return shim.Success([]byte("Entity not found"))
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Success([]byte("Failed to get state"))
	}
	if Bvalbytes == nil {
		return shim.Success([]byte("Entity not found"))
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Success([]byte("Invalid transaction amount, expecting a integer value"))
	}
	Aval = Aval - X
	Bval = Bval + X

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// IfAccountExisted function is used for see if the account has been existed.
// 1 input
// "account"
// ============================================================================================================================
func (t *SimpleChaincode) IfAccountExisted(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Success([]byte("Incorrect number of arguments. "))
	}
	account := args[0]
	// Handle different functions
	password, err := stub.GetState(account) //get the var from chaincode state
	if err != nil || password == nil {
		return shim.Success([]byte("account not found"))
	}

	return shim.Success([]byte("ok"))
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
