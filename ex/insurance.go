/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
//	"bufio"
//	"bytes"
//	"crypto/aes"
//	"crypto/cipher"
//	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
//	"image"
//	"image/gif"
//	"image/jpeg"
//	"image/png"
//	"io"
//	"net/http"
	"os"
//	"os/exec"
//	"runtime"
	"strconv"
//	"strings"
//	"time"


)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var gopath string
var ccPath string

//////////////////////////////////////////////////////////////////////////////////////////////////
// The following array holds the list of tables that should be created
// The deploy/init deletes the tables and recreates them every time a deploy is invoked
//////////////////////////////////////////////////////////////////////////////////////////////////
var insuranceTables = []string{"UserTable", "Contracts"}

//////////////////////////////////////////////////////////////////////////////////////////////////
// The recType is a mandatory attribute. The original app was written with a single table
// in mind. The only way to know how to process a record was the 70's style 80 column punch card
// which used a record type field. The array below holds a list of valid record types.
// This could be stored on a blockchain table or an application
//////////////////////////////////////////////////////////////////////////////////////////////////
var recType = []string{"CONTRACT", "USER"}


/////////////////////////////////////////////////////////////
// Create Buyer, Seller , Auction House, Authenticator
// Could establish valid UserTypes -
// AH (Auction House)
// TR (Buyer or Seller)
// AP (Appraiser)
// IN (Insurance)
// BK (bank)
// SH (Shipper)
/////////////////////////////////////////////////////////////
type UserObject struct {
	UserID    string
	RecType   string // Type = USER
	Name      string
	UserType  string // Auction House (AH), Bank (BK), Buyer or Seller (TR), Shipper (SH), Appraiser (AP)
	Address   string
	Phone     string
	Email     string
	Bank      string
	AccountNo string
	RoutingNo string
}

type ContractObject struct {
	ContractID    string
	RecType   string // Type = USER
	Status    string
	Price    string 
	CreateDate   string
	CustomerName     string
	CustomerAge     string
	Seller     string
	DoctorName      string
	DoctorComment      string

}


/////////////////////////////////////////////////////////////////////////////////////////////////////
// A Map that holds TableNames and the number of Keys
// This information is used to dynamically Create, Update
// Replace , and Query the Ledger
// In this model all attributes in a table are strings
// The chain code does both validation
// A dummy key like 2016 in some cases is used for a query to get all rows
//
//              "UserTable":        1, Key: UserID
//              "ItemTable":        1, Key: ItemID
//              "UserCatTable":     3, Key: "2016", UserType, UserID
//              "ItemCatTable":     3, Key: "2016", ItemSubject, ItemID
//              "AuctionTable":     1, Key: AuctionID
//              "AucInitTable":     2, Key: Year, AuctionID
//              "AucOpenTable":     2, Key: Year, AuctionID
//              "TransTable":       2, Key: AuctionID, ItemID
//              "BidTable":         2, Key: AuctionID, BidNo
//              "ItemHistoryTable": 4, Key: ItemID, Status, AuctionHouseID(if applicable),date-time
//
/////////////////////////////////////////////////////////////////////////////////////////////////////

func GetNumberOfKeys(tname string) int {
	TableMap := map[string]int{
		"UserTable":        1,
		"Contracts":        2,
	}
	return TableMap[tname]
}
//////////////////////////////////////////////////////////////
// Invoke Functions based on Function name
// The function name gets resolved to one of the following calls
// during an invoke
//
//////////////////////////////////////////////////////////////
func InvokeFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	InvokeFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
		"PostCustomerRequest":      PostCustomerRequest,
//		"PostDoctorAnswer": PostDoctorAnswer,
//		"PostSellerAnswer": PostSellerAnswer,
	}
	return InvokeFunc[fname]
}

//////////////////////////////////////////////////////////////
// Query Functions based on Function name
//
//////////////////////////////////////////////////////////////
func QueryFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	QueryFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
		"GetContract":           GetContract,
//		"GetListOfContractByStatus":           GetListOfContractByStatus,
	}
	return QueryFunc[fname]
}




/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a master Object of the Item
// Since the Owner Changes hands, a record has to be written for each
// Transaction with the updated Encryption Key of the new owner
// Example
//./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostItem", "Args":["1000", "ARTINV", "Shadows by Asppen", "Asppen Messer", "20140202", "Original", "Landscape" , "Canvas", "15 x 15 in", "sample_7.png","$600", "100"]}'
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PostCustomerRequest(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	itemObject, err := CreateItemObject(args[0:])
	if err != nil {
		fmt.Println("PostItem(): Cannot create item object \n")
		return nil, err
	}

	// Check if the Owner ID specified is registered and valid
	//ownerInfo, err := ValidateMember(stub, itemObject.CurrentOwnerID)
	//fmt.Println("Owner information  ", ownerInfo, itemObject.CurrentOwnerID)
	//if err != nil {
	//	fmt.Println("PostItem() : Failed Owner information not found for ", itemObject.CurrentOwnerID)
	//	return nil, err
	//}

	// Convert Item Object to JSON
	buff, err := ContracttoJSON(itemObject) //
	if err != nil {
		fmt.Println("PostItem() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("PostItem(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		keys := []string{args[0]}
		err = UpdateLedger(stub, "Contracts", keys, buff)
		if err != nil {
			fmt.Println("PostItem() : write error while inserting record\n")
			return buff, err
		}


	}

	return buff, nil
}

func CreateItemObject(args []string) (ContractObject, error) {

	var err error
	var bufferErrTxt bytes.Buffer
	var myItem ContractObject

	// Check there are 12 Arguments provided as per the the struct - two are computed
	if len(args) != 10 {
		fmt.Println("CreateItemObject(): Incorrect number of arguments. Expecting 12 ", len(args))
		bufferErrTxt.WriteString("CreateItemObject(): Incorrect number of arguments. Expecting 10 ")
		bufferErrTxt.WriteString(args[0])

		return myItem, errors.New(bufferErrTxt.String())
	}

	// Validate ItemID is an integer

	_, err = strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("CreateItemObject(): ART ID should be an integer create failed! ")
		
		return myItem, errors.New("CreateItemObject(): ART ID should be an integer create failed!")
	}


	myItem = ContractObject{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9]}

	fmt.Println("CreateItemObject(): Item Object created: ID# ", myItem.ContractID)

	// Code to Validate the Item Object)
	// If User presents Crypto Key then key is used to validate the picture that is stored as part of the title
	// TODO

	return myItem, nil
}

func GetContract(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//var err error

	// Get the Objects and Display it
	Avalbytes, err := QueryLedger(stub, "Contracts", args)
	if Avalbytes == nil {
		fmt.Println("GetTransaction() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if err != nil {
		fmt.Println("GetTransaction() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetTransaction() : Response : Successfull")
	return Avalbytes, nil
}
//////////////////////////////////////////////////////////
// Converts an ART Object to a JSON String
//////////////////////////////////////////////////////////
func ContracttoJSON(ar ContractObject) ([]byte, error) {

	ajson, err := json.Marshal(ar)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}
//////////////////////////////////////////////////////////
// Converts an User Object to a JSON String
//////////////////////////////////////////////////////////
func JSONtoUser(user []byte) (UserObject, error) {

	ur := UserObject{}
	err := json.Unmarshal(user, &ur)
	if err != nil {
		fmt.Println("JSONtoUser error: ", err)
		return ur, err
	}
	fmt.Println("JSONtoUser created: ", ur)
	return ur, err
}
//////////////////////////////////////////////////////////
// Converts an Contract Object to a JSON String
//////////////////////////////////////////////////////////
func JSONtoContract(contract []byte) (ContractObject, error) {

	ur := ContractObject{}
	err := json.Unmarshal(contract, &ur)
	if err != nil {
		fmt.Println("JSONtoContart error: ", err)
		return ur, err
	}
	fmt.Println("JSONtoContartc created: ", ur)
	return ur, err
}

////////////////////////////////////////////////////////////////////////////
// Open a Ledgers if one does not exist
// These ledgers will be used to write /  read data
// Use names are listed in aucTables {}
// THIS FUNCTION REPLACES ALL THE INIT Functions below
//  - InitUserReg()
//  - InitAucReg()
//  - InitBidReg()
//  - InitItemReg()
//  - InitItemMaster()
//  - InitTransReg()
//  - InitAuctionTriggerReg()
//  - etc. etc.
////////////////////////////////////////////////////////////////////////////
func InitLedger(stub shim.ChaincodeStubInterface, tableName string) error {

	// Generic Table Creation Function - requires Table Name and Table Key Entry
	// Create Table - Get number of Keys the tables supports
	// This version assumes all Keys are String and the Data is Bytes
	// This Function can replace all other InitLedger function in this app such as InitItemLedger()

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
		fmt.Println("Auction_Application: Failed creating Table ", tableName)
		return errors.New("Auction_Application: Failed creating Table " + tableName)
	}

	var columnDefsForTbl []*shim.ColumnDefinition

	for i := 0; i < nKeys; i++ {
		columnDef := shim.ColumnDefinition{Name: "keyName" + strconv.Itoa(i), Type: shim.ColumnDefinition_STRING, Key: true}
		columnDefsForTbl = append(columnDefsForTbl, &columnDef)
	}

	columnLastTblDef := shim.ColumnDefinition{Name: "Details", Type: shim.ColumnDefinition_BYTES, Key: false}
	columnDefsForTbl = append(columnDefsForTbl, &columnLastTblDef)

	// Create the Table (Nil is returned if the Table exists or if the table is created successfully
	err := stub.CreateTable(tableName, columnDefsForTbl)

	if err != nil {
		fmt.Println("Auction_Application: Failed creating Table ", tableName)
		return errors.New("Auction_Application: Failed creating Table " + tableName)
	}

	return err
}

////////////////////////////////////////////////////////////////////////////
// Open a User Registration Table if one does not exist
// Register users into this table
////////////////////////////////////////////////////////////////////////////
func UpdateLedger(stub shim.ChaincodeStubInterface, tableName string, keys []string, args []byte) error {

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
	}

	var columns []*shim.Column

	for i := 0; i < nKeys; i++ {
		col := shim.Column{Value: &shim.Column_String_{String_: keys[i]}}
		columns = append(columns, &col)
	}

	lastCol := shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(args)}}
	columns = append(columns, &lastCol)

	row := shim.Row{columns}
	ok, err := stub.InsertRow(tableName, row)
	if err != nil {
		return fmt.Errorf("UpdateLedger: InsertRow into "+tableName+" Table operation failed. %s", err)
	}
	if !ok {
		return errors.New("UpdateLedger: InsertRow into " + tableName + " Table failed. Row with given key " + keys[0] + " already exists")
	}

	fmt.Println("UpdateLedger: InsertRow into ", tableName, " Table operation Successful. ")
	return nil
}
////////////////////////////////////////////////////////////////////////////
// Query a User Object by Table Name and Key
////////////////////////////////////////////////////////////////////////////
func QueryLedger(stub shim.ChaincodeStubInterface, tableName string, args []string) ([]byte, error) {

	var columns []shim.Column
	nCol := GetNumberOfKeys(tableName)
	for i := 0; i < nCol; i++ {
		colNext := shim.Column{Value: &shim.Column_String_{String_: args[i]}}
		columns = append(columns, colNext)
	}

	row, err := stub.GetRow(tableName, columns)
	fmt.Println("Length or number of rows retrieved ", len(row.Columns))

	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving data " + args[0] + ". \"}"
		fmt.Println("Error retrieving data record for Key = ", args[0], "Error : ", jsonResp)
		return nil, errors.New(jsonResp)
	}

	//fmt.Println("User Query Response:", row)
	//jsonResp := "{\"Owner\":\"" + string(row.Columns[nCol].GetBytes()) + "\"}"
	//fmt.Println("User Query Response:%s\n", jsonResp)
	Avalbytes := row.Columns[nCol].GetBytes()

	// Perform Any additional processing of data
	fmt.Println("QueryLedger() : Successful - Proceeding to ProcessRequestType ")
	err = ProcessQueryResult(stub, Avalbytes, args)
	if err != nil {
		fmt.Println("QueryLedger() : Cannot create object  : ", args[1])
		jsonResp := "{\"QueryLedger() Error\":\" Cannot create Object for key " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////
// Return the right Object Buffer after validation to write to the ledger
// var recType = []string{"ARTINV", "USER", "BID", "AUCREQ", "POSTTRAN", "OPENAUC", "CLAUC"}
/////////////////////////////////////////////////////////////////////////////////////////////

func ProcessQueryResult(stub shim.ChaincodeStubInterface, Avalbytes []byte, args []string) error {

	// Identify Record Type by scanning the args for one of the recTypes
	// This is kind of a post-processor once the query fetches the results
	// RecType is the style of programming in the punch card days ..
	// ... well

	var dat map[string]interface{}

	if err := json.Unmarshal(Avalbytes, &dat); err != nil {
		panic(err)
	}

	var recType string
	recType = dat["RecType"].(string)
	switch recType {

	case "USER":
		ur, err := JSONtoUser(Avalbytes) //
		if err != nil {
			return err
		}
		fmt.Println("ProcessRequestType() : ", ur)
		return err
	case "CONTRACT":
		ur, err := JSONtoContract(Avalbytes) //
		if err != nil {
			return err
		}
		fmt.Println("ProcessRequestType() : ", ur)
		return err

	default:

		return errors.New("Unknown")
	}
	return nil

}


////////////////////////////////////////////////////////////////////////////
// Open a User Registration Table if one does not exist
// Register users into this table
////////////////////////////////////////////////////////////////////////////
func DeleteFromLedger(stub shim.ChaincodeStubInterface, tableName string, keys []string) error {
	var columns []shim.Column

	//nKeys := GetNumberOfKeys(tableName)
	nCol := len(keys)
	if nCol < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
		return errors.New("DeleteFromLedger failed. Must include at least key values")
	}

	for i := 0; i < nCol; i++ {
		colNext := shim.Column{Value: &shim.Column_String_{String_: keys[i]}}
		columns = append(columns, colNext)
	}

	err := stub.DeleteRow(tableName, columns)
	if err != nil {
		return fmt.Errorf("DeleteFromLedger operation failed. %s", err)
	}

	fmt.Println("DeleteFromLedger: DeleteRow from ", tableName, " Table operation Successful. ")
	return nil
}



























func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// TODO - Include all initialization to be complete before Invoke and Query
	// Uses aucTables to delete tables if they exist and re-create them

	//myLogger.Info("[Trade and Auction Application] Init")
	fmt.Println("[Trade and Auction Application] Init")
	var err error

	for _, val := range insuranceTables {
		err = stub.DeleteTable(val)
		if err != nil {
			return nil, fmt.Errorf("Init(): DeleteTable of %s  Failed ", val)
		}
		err = InitLedger(stub, val)
		if err != nil {
			return nil, fmt.Errorf("Init(): InitLedger of %s  Failed ", val)
		}
	}
	// Update the ledger with the Application version
	err = stub.PutState("version", []byte(strconv.Itoa(1)))
	if err != nil {
		return nil, err
	}

	fmt.Println("Init() Initialization Complete  : ", args)
	return []byte("Init(): Initialization Complete"), nil

/*
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	fmt.Printf("test JCO3")

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
*/
}

// Transaction makes payment of X units from A to B

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var buff []byte

	if ChkReqType(args) == true {

		InvokeRequest := InvokeFunction(function)
		if InvokeRequest != nil {
			buff, err = InvokeRequest(stub, function, args)
		}
	} else {
		fmt.Println("Invoke() Invalid recType : ", args, "\n")
		return nil, errors.New("Invoke() : Invalid recType : " + args[0])
	}

	return buff, err
/*
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	fmt.Printf("test JCO2")
	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Invalid transaction amount, expecting a integer value")
	}
	// Aval = Aval - X
	// Bval = Bval + X

	Aval = Aval - 10 * X
	Bval = Bval + 100* X
	fmt.Printf("JCO Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
*/
}


/////////////////////////////////////////////////////////////////
// This function checks the incoming args stuff for a valid record
// type entry as per the declared array recType[]
// The assumption is that rectType can be anywhere in the args or struct
// not necessarily in args[1] as per my old logic
// The Request type is used to process the record accordingly
/////////////////////////////////////////////////////////////////
func ChkReqType(args []string) bool {
	for _, rt := range args {
		for _, val := range recType {
			if val == rt {
				return true
			}
		}
	}
	return false
}


// Deletes an entity from state
/*
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}
*/
// Query callback representing the query of a chaincode

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var buff []byte
	fmt.Println("ID Extracted and Type = ", args[0])
	fmt.Println("Args supplied : ", args)

	if len(args) < 1 {
		fmt.Println("Query() : Include at least 1 arguments Key ")
		return nil, errors.New("Query() : Expecting Transation type and Key value for query")
	}

	QueryRequest := QueryFunction(function)
	if QueryRequest != nil {
		buff, err = QueryRequest(stub, function, args)
	} else {
		fmt.Println("Query() Invalid function call : ", function)
		return nil, errors.New("Query() : Invalid function call : " + function)
	}

	if err != nil {
		fmt.Println("Query() Object not found : ", args[0])
		return nil, errors.New("Query() : Object not found : " + args[0])
	}
	return buff, err
/*
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
*/
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	fmt.Printf("test JCO1")
	
	gopath = os.Getenv("GOPATH")
    ccPath = fmt.Sprintf("%s/src/github.com/hyperledger/fabric/auction/art/artchaincode/", gopath)

	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}


}
