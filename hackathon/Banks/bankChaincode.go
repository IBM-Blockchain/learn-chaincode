/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"time"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var bitgramIndexStr = "_bitgramindex"			//name for the key/value that will store a list of all known bitgrams
var bitgramTradesStr = "_alltrades"				//name for the key/value that will store all open trades

// Name, Address, FinacialScore, GovScore, SocialScore, EmployerScore, QualificationScore, KYCScore, Namespace(the first transacting bank)


/* order of arguments
'IN_BITGRAM_ID'
'IN_SMART_ID_NAMESPACE'
'IN_SMART_ID_NAME'
'IN_SMART_ID_ADDR'
'IN_SMART_ID_FINSC'
'IN_SMART_ID_GOVSC'
'IN_SMART_ID_SOCSC'
'IN_SMART_ID_EMPSC'
'IN_SMART_ID_PUBLICSC'
*/

type Bitgram struct{
	Name string `json:"name"`					   //the fieldtags are needed to keep case from bouncing around
	Address string `json:"address"`
	FS string `json:"fs"`
    GS string `json:"gs"`
	ES string `json:"es"`
	QS string `json:"qs"`
	KS string `json:"ks"`
	NS string `json:"ns"`
}


type Trade struct{
	Bitgram string `json:"bitgram"`					//user who created the open trade order
	Timestamp int64 `json:"timestamp"`			//utc timestamp of creation
	Tobank string `json:"tobank"`		               //bought by
}

type AllTrades struct{
	BitgramTrades []Trade `json:"open_trades"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(bitgramIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades)								//clear the open trade struct
	err = stub.PutState(bitgramTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "writeBitgramIdentity" {			    //updates the bitgram identity record to the chaincode state
		return t.writeBitgramIdentity(stub, args)
	//} else if function == "shareIdentity" {							//shares a new trade order
		//return t.shareIdentity(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)		   //error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================

func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// writeBitgramIdentity - create/update a new bitgram identity, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) writeBitgramIdentity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	
    str := `{ "IN_SMART_ID_NAMESPACE":"` + args[1] + `","IN_SMART_ID_NAME":"`+ args[2] + `","IN_SMART_ID_ADDR":"`+ args[3] + `","IN_SMART_ID_FINSC":"`+ args[4] + `","IN_SMART_ID_GOVSC":"`+ args[5] + `","IN_SMART_ID_SOCSC":"`+ args[6] +`","IN_SMART_ID_EMPSC":"`+ args[7] +`","IN_SMART_ID_PUBLICSC":"`+ args[8] +`"}`


	err = stub.PutState(args[0], []byte(str))								//store bitgram with id as key
	if err != nil {
		return nil, err
	}
		
	//get the bitgram index
	bitgramsAsBytes, err := stub.GetState(bitgramIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get bitgram index")
	}
	var bitgramIndex []string
	json.Unmarshal(bitgramsAsBytes, &bitgramIndex)							//un stringify it aka JSON.parse()
	
	//append
	bitgramIndex = append(bitgramIndex, args[0])								//add bitgram name to index list
	fmt.Println("! bitgram index: ", bitgramIndex)
	jsonAsBytes, _ := json.Marshal(bitgramIndex)
	err = stub.PutState(bitgramIndexStr, jsonAsBytes)						//store name of bitgram

	fmt.Println("- end init bitgram")
	return nil, nil
}



// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
