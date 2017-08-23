// ============================================================================================================================
// This smart contract is used for managing the accounts
// It includes adding an account ,deleting an account,changing an
// account's password and verifying the account's password
// ============================================================================================================================
// 本智能合约用于用户账户管理
// 功能包括：增加账号 删除账号 修改管理账号 对登录账号进行密码验证
// ============================================================================================================================

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("account")

type SimpleChaincode struct {
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

	// Initialize the chaincode
	ID = args[0]
	IDval = args[1]
	logger.Info("ID = %s, IDval = %s\n", ID, IDval)

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

	if function == "delete" { //deletes an account from its state
		return t.Delete(stub, args)
	}
	if function == "add" { //add a new account
		return t.Add(stub, args)
	}
	if function == "edit" { //change the password of the account
		return t.Edit(stub, args)
	}
	if function == "reset" { // reset the password of the account
		return t.Reset(stub, args)
	}
	if function == "verify" { //verify the account and password
		return t.Verify(stub, args)
	}
	if function == "ifAccountExisted" { // reply if the account is existed
		return t.IfAccountExisted(stub, args)
	}
	if function == "test" { //a function for testing
		return t.Test(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'add', 'edit','reset', 'verify','ifAccountExisted', 'test' . But got: %v", function)
	return shim.Success([]byte(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'add', 'edit','reset', 'verify','ifAccountExisted', 'test' . But got: %v", function)))
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
// "account","password"
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
// Edit function is used for changing the account's password
// 3 input
// "account","old password","new password"
// ============================================================================================================================
func (t *SimpleChaincode) Edit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	account := args[0]
	oldPassword := args[1]
	newPassword := args[2]
	accountTest, err := stub.GetState(account)

	//test if the account has been existed
	if err != nil {
		return shim.Success([]byte("Failed to get state"))
	}
	if accountTest == nil {
		return shim.Success([]byte("account not found"))
	}
	if oldPassword != string(accountTest) {
		return shim.Success([]byte("old password is wrong"))
	}

	// edit the account
	err = stub.PutState(account, []byte(newPassword))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Reset function is used for resetting the account's password in case of forgetting the password
// 1 input
// "account"
// ============================================================================================================================
func (t *SimpleChaincode) Reset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Success([]byte("Incorrect number of arguments."))
	}
	account := args[0]
	password := args[0]
	accountTest, err := stub.GetState(account)

	//test if the account has been existed
	if err != nil {
		return shim.Success([]byte("account is not here"))
	}
	if accountTest == nil {
		return shim.Success([]byte("account not found"))
	}

	// reset the account's password
	err = stub.PutState(account, []byte(password))
	if err != nil {
		return shim.Success([]byte(err.Error()))
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Reset function is used for verifying the password,if the password is correct,there will be a reply of "ok".
// 2 input
// "account","password"
// ============================================================================================================================
func (t *SimpleChaincode) Verify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Success([]byte("Incorrect number of arguments. "))
	}
	account := args[0]
	password := args[1]
	accountTest, err := stub.GetState(account)
	//test if the account has been existed
	if err != nil || accountTest == nil {
		return shim.Success([]byte("account not found"))
	}

	// verify
	if password == string(accountTest) {
		return shim.Success([]byte("ok"))

	} else {
		return shim.Success([]byte("The password is not correct"))
	}

}

// ============================================================================================================================
// Reset function is used for making sure the existing of the account ,if the account is existed,there will be a reply of "ok".
// 2 input
// "account","password"
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
