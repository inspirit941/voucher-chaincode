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
	  "sync"
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
	 Used				bool // Voucher 연산에 이미 쓰였는지 아닌지 파악할 변수
 }
 type CourseStatistics struct {
	 CourseId	string
	 AvgFocusRate	map[int]float64
	 // 해당 lecture 개수
	 Count		map[int]int
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
	 } else if function == "updateCourseStatistics" {
		 return s.updateCourseStatistics(APIstub, args)
	 } else if function == "CalculateVoucher" { 
		 return s.CalculateVoucher(APIstub, args)
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
 // courseid, lecture_number값, 그 외에는 focus_rate 같은 통계처리에 필요한 값을 인자로 받는다
 // https://blog.golang.org/go-maps-in-action. map의 동시성 처리 문제
 func (s *SmartContract) updateCourseStatistics(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	 if len(args) != 3 {
		 return shim.Error("Incorrect number of arguments in updateCourseStatistics.")
	 }
	 courseid := args[0]
	 lecture_number, err1 := strconv.Atoi(args[1])

	 if err1 != nil {
		 return shim.Error(err1.Error())
	 }
	 focus_rate, err2 := strconv.ParseFloat(args[2], 64)
	 if err2 != nil {
		 return shim.Error(err2.Error())
	 }

	 course_statistics, _ := APIstub.GetState(courseid)
	 // 해당 Asset이 존재하지 않는 경우, Asset을 생성한다.
	 if course_statistics == nil {
		var course = CourseStatistics{CourseId: courseid, 
			AvgFocusRate : make(map[int]float64), 
			Count : make(map[int]int) }
			
		course.AvgFocusRate[lecture_number] = focus_rate
		course.Count[lecture_number] = 1

		course_statisticsAsBytes, _ := json.Marshal(course)
		APIstub.PutState(courseid, course_statisticsAsBytes)
		return shim.Success(course_statisticsAsBytes)

	} 

	// 해당 Asset이 이미 존재할 경우
	course := CourseStatistics{}
	json.Unmarshal(course_statistics, &course)

	// 해당 asset에 값 업데이트 -> 동시성 처리 작업
	var counter = struct{
		sync.RWMutex
		course CourseStatistics
	}{course: course}

	// 1. 해당 lecture_number의 focus_rate 평균값 바꿔주기.
	counter.Lock()
	rate_sum := (counter.course.AvgFocusRate[lecture_number] * float64(counter.course.Count[lecture_number])) + focus_rate
	counter.course.Count[lecture_number] += 1
	counter.course.AvgFocusRate[lecture_number] = (rate_sum / float64(counter.course.Count[lecture_number]))
	counter.Unlock()
	
	// 값 저장하고 state에 저장
	course_statisticsAsBytes, _ := json.Marshal(course)
	APIstub.PutState(courseid, course_statisticsAsBytes)
	return shim.Success(course_statisticsAsBytes)

 }
 func (s *SmartContract) createLecture(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//  var jsonResp string
	
	 
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
		
		// 해당 courseid, lecture number의 통계값 업데이트
		lectureAsBytes, _ := json.Marshal(lecture)
		APIstub.PutState(LectureId, lectureAsBytes)
		return s.updateCourseStatistics(APIstub, []string{courseid, strconv.Itoa(lecture_number),
			 fmt.Sprintf("%f", focus_rate)})
	 	// return shim.Success(lectureAsBytes)
		
		
	 } else {
		// 해당 asset이 존재하지 않을 경우 
		var lecture = Lecture{Sid: sid, Courseid: courseid, 
			Lecture_fin_date : []int{lecture_fin_date}, 
			Lecture_number : []int{lecture_number},
			Focus_rate : []float64{focus_rate}, 
			Used : false }
		lectureAsBytes, _ := json.Marshal(lecture)
		APIstub.PutState(LectureId, lectureAsBytes)
		
		// sid, courseid, 고유 assetid (lecture의 경우 1)도 key값이어야 하므로, 
		// compositekey 함수 사용해서 key 생성
		// 20.02.22 현재 CompositeKey가 작동하는지는 아직 확인하지 못함
		indexName := "Sid~CourseId~AssetNum"
		SidCourseIdAssetNum, err := APIstub.CreateCompositeKey(indexName, []string{lecture.Sid, lecture.Courseid, "1"})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		APIstub.PutState(SidCourseIdAssetNum, value)

		// 해당 courseid, lecture number의 통계값 업데이트
		// s.updateCourseStatistics(APIstub, []string{courseid, strconv.Itoa(lecture_number), fmt.Sprintf("%f", focus_rate)})
		return s.updateCourseStatistics(APIstub, []string{courseid, strconv.Itoa(lecture_number), fmt.Sprintf("%f", focus_rate)})
	 	// return shim.Success(lectureAsBytes)

		
	 }
	 
 }

// https://medium.com/@kctheservant/putstate-and-getstate-the-api-in-chaincode-dealing-with-the-state-in-the-ledger-part-2-839f89ecbad4 
// 사용자가 해당 강좌를 완강한 뒤 Voucher를 계산받고자 하는 경우.
 func (s *SmartContract) CalculateVoucher(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	 // sid, courseid 두 개만 인자로 받을 경우
	 
	 if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
	 sid := args[0]
	 courseid := args[1] 
	 // Lecture 가져오기
	 LectureId := sid + "_" + courseid + "_1"
	 lecturejson, err := APIstub.GetState(LectureId)
	// 데이터가 없는 경우
	 if err != nil {
		 return shim.Error(err.Error())
	 }

	 // 데이터 가져와서 json으로 변형하기
	 var lecture = Lecture{}
	 json.Unmarshal(lecturejson, &lecture)

	 // 1. 중복 수강한 강좌의 경우, Focus_rate이 가장 높은 것만 가져오기
	 lecture_max_focus := make(map[int]float64)
	 for i := 0; i < len(lecture.Lecture_number); i++ {
		val, ok := lecture_max_focus[lecture.Lecture_number[i]]
		if ok {
			if val < lecture.Focus_rate[i] {
				lecture_max_focus[lecture.Lecture_number[i]] = lecture.Focus_rate[i]
			}
		} else {
			lecture_max_focus[lecture.Lecture_number[i]] = lecture.Focus_rate[i]
		}
	 }
	 // focus rate 평균 구하기
	 rate_sum := float64(0)
	 for _, value := range lecture_max_focus {
		 rate_sum += value
	 }
	 // 모든 강의 focus_rate 평균값
	 total_AvgFocusRate := rate_sum / float64(len(lecture_max_focus))
	 
	 // 해당 course의 focus rate 평균과 비교해서 클 경우 예컨대 크레딧 1.2배 제공
	 if total_AvgFocusRate > 90.0 {
		 return shim.Error("total_AvgFocusRate")
	 } else {
		 return shim.Error("fails")
	 }

	//  return shim.Success(nil)
	 
	 


	//  resultiterator, err := APIstub.getStateByPartialCompositeKey("Sid~CourseId~Assetnum", []string{sid, courseid})
	//  if err != nil {
	// 	 return shim.Error("Cannot get Asset from Sid and CourseId.")
	//  }
	//  defer resultiterator.close()
	//  var i int
	//  for i = 0; resultiterator.HasNext(); i++ {
	// 	 responseRange, err := resultiterator.Next() 
	// 	 if err != nil {
	// 		 return shim.Error(err.Error())
	// 	 }
	// 	_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}	 
		

	//  }
	//  // 1. Focus_rate 평균값 구하기



 }
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 