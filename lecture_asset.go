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
	//  "time"
	//  "math"
	 
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer"
 )
 
 // Define the Smart Contract structure
 type SmartContract struct {
 }
 
 // Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
 // type Car struct {
 // 	Make   string `json:"make"`
 // 	Model  string `json:"model"`
 // 	Colour string `json:"colour"`
 // 	Owner  string `json:"owner"`
 // }
 
 // 
 type Lecture struct {
	 Sid               string `json:Sid`
	 Courseid          string `json:Courseid`
	 Lecture_fin_date  []int `json:Lecture_fin_date`
	 Lecture_number    []int	`json:Lecture_number`
	 Focus_rate        []float64 `json:Focus_rate`
 }
 
 /*
  * The Init method is called when the Smart Contract "opprecord" is instantiated by the blockchain network
  * Best practice is to have any Ledger initialization in separate function -- see initLedger()
  */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 fmt.Println("Lecture_stat Init")
	 return shim.Success(nil)
 }
 
 /*
  * The Invoke method is called as a result of an application request to run the Smart Contract "opprecord"
  * The calling application program has also specified the particular smart contract function to be called, with arguments
  */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested Smart Contract function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	//  fmt.Println(function)
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "queryRecord" {
		 return s.queryRecord(APIstub, args)
	//  } else if function == "initLedger" {
	// 	 return s.initLedger(APIstub)
	 } else if function == "createLecture" {
		 return s.createLecture(APIstub, args)
	 } else if function == "dataToFabric" {
		 return s.dataToFabric(APIstub, args)
	 }
	 //  else if function == "queryAllRecordss" {
	 // 	return s.queryAllCars(APIstub)
	 // }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) queryRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 var prime_key = args[0]
	 recordAsBytes, _ := APIstub.GetState(prime_key)
	 return shim.Success(recordAsBytes)
 }

 func (s *SmartContract) dataToFabric(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string

	// args ...string?
	// input 개수가 늘어나면 변경해야 함
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments.")
	}
	sid := args[0]
	courseid := args[1]
	lecture_fin_date, err1 := strconv.Atoi(args[2])
	lecture_number, err2 := strconv.Atoi(args[3])
	total_lectureTime, err3 := strconv.Atoi(args[4])
	focus_lectureTime, err4 := strconv.Atoi(args[5])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		jsonResp = "{\"Error\":\"string to integer / float conversion Failed" + "\"}"
		return shim.Error(jsonResp)
	 }


	// focus rate -> focus_time / total unix (단순비율)
	focus_rate := float64(focus_lectureTime) / float64(total_lectureTime)

	// total_time := time.Unix(total_lectureTime)
	// focus_time := time.Unix(focus_lectureTime)
	// 소수점 둘째 자리에서 반올림. 자릿수 바꾸려면 100을 바꾸면 됨
	// focus_rate := focus_time / total_time

	// Lecture Asset 생성
	// s.createLecture(APIstub, []string{sid, courseid, strconv.Itoa(lecture_fin_date),
	// 	strconv.Itoa(lecture_number), fmt.Sprintf("%.2f", focus_rate)})


	return s.createLecture(APIstub, []string{sid, courseid, strconv.Itoa(lecture_fin_date),
		strconv.Itoa(lecture_number), fmt.Sprintf("%f", focus_rate)})

 }
 
 func (s *SmartContract) createLecture(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	 var jsonResp string
	
	 
	 if len(args) != 5 {
		 return shim.Error("Incorrect number of arguments. Expecting 5")
	 }
	 sid := args[0]
	 courseid := args[1]
	 lecture_fin_date, err1 := strconv.Atoi(args[2])
	 lecture_number, err2 := strconv.Atoi(args[3])
	 focus_rate, err3 := strconv.ParseFloat(args[4], 64)
	 if err1 != nil {
		 return shim.Error("lecture_fin_date")
	 } else if err2 != nil {
		 return shim.Error("lecture_number")
	 } else if err3 != nil {
		 return shim.Error("focus_rate")
	 } 
	 if err1 != nil || err2 != nil || err3 != nil {
		jsonResp = "{\"Error\":\"string to integer / float conversion Failed" + "\"}"
		return shim.Error(jsonResp)
	 }

	 LectureId := sid + "_" + courseid + "_1"
	 lecturejson, _ := APIstub.GetState(LectureId)

	 // 이미 Lecture가 존재하는 경우
	 // https://github.com/wikibook/blockchain/blob/master/13%EC%9E%A5_%EB%A6%AC%EC%8A%A4%ED%8A%B819_chaincode_counter.go
	 // http://golang.site/go/article/13-Go-%EC%BB%AC%EB%A0%89%EC%85%98---Slice
	 if lecturejson != nil {
		lecture := Lecture{}
		json.Unmarshal(lecturejson, &lecture)
		// 해당 Asset에 값 업데이트.
		lecture.Lecture_fin_date = append(lecture.Lecture_fin_date, lecture_fin_date)
		lecture.Lecture_number = append(lecture.Lecture_number, lecture_number)
		lecture.Focus_rate = append(lecture.Focus_rate, focus_rate)

		// Event?
		
		//
		lectureAsBytes, _ := json.Marshal(lecture)
		APIstub.PutState(LectureId, lectureAsBytes)
		return shim.Success(lectureAsBytes)
	 } else {
		// 해당 asset이 존재하지 않을 경우 
		var lecture = Lecture{Sid: sid, Courseid: courseid, 
			Lecture_fin_date : []int{lecture_fin_date}, 
			Lecture_number : []int{lecture_number},
			Focus_rate : []float64{focus_rate}}
		lectureAsBytes, _ := json.Marshal(lecture)
		APIstub.PutState(LectureId, lectureAsBytes)
		return shim.Success(lectureAsBytes)
	 }
 }
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 