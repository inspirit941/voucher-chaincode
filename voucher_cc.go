/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 5 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Record structure, with 11 properties.  Structure tags are used by encoding/json library
type Wallet struct {
	Address string `json:address`
	Name    string `json:name`
	Balance int    `json:balance`
	Org     string `json:org`
}

/*
 * The Init method is called when the Smart Contract "opprecord" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("ex02 Init")
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "opprecord"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryBalance" {
		return s.queryRecord(APIstub, args)
	} else if function == "transfer" {
		return s.transfer(APIstub, args)
	} else if function == "createAccount" {
		return s.createAccount(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryBalance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var address = args[0]
	walletAsBytes, _ := APIstub.GetState(address)
	// WalletData, err := json.Unmarshal(walletAsBytes)
	return shim.Success(walletAsBytes)
}

func (s *SmartContract) transfer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 11")
	}

	var fromAddress = args[0]
	var toAddress = args[1]
	var balance = strconv.Atoi(args[2])

	walletAsBytes, _ := APIstub.GetState(fromAddress)
	toWalletAsBytes, _ := APIstub.GetState(toAddress)

	WalletFrom := Wallet{}
	WalletTo := Wallet{}

	err = json.Unmarshal(walletAsBytes, &WalletFrom) //unmarshal it aka JSON.parse()
	err = json.Unmarshal(toWalletAsBytes, &WalletTo)

	if WalletFrom.Balance < balance {
		return shim.Error("잔액이 부족합니다.")
	}

	WalletFrom.Balance -= balance
	WalletTo.Balance += balance

	fromWalletasBytes, _ := json.Marshal(WalletFrom)
	toWalletasBytes, _ := json.Marshal(WalletTo)

	err = APIstub.PutState(fromAddress, fromWalletasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(toAddress, toWalletasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(fromWalletasBytes)
}

// marshall - json을 바이트 형태로 변환
// unmarshall - 바이트를 json으로 변환
func (s *SmartContract) createWallet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var address =  args[0]
	var name = args[1]
	var balance, _ = strconv.Atoi(args[2])
	var org = args[3]
	
	var newWallet = Wallet{Address : address, Name : name, Balance : balance, Org : org}

	walletAsBytes, err := json.Marshal(newWallet)
	if err != nil {
		return shim.Error("Json_to_Bytes Error.")
	} else if walletAsBytes == nil {
		return shim.Error("Marble does not exits.")
	}

	APIstub.PutState(address, walletAsBytes)
	
	return shim.Success(walletAsBytes)
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
