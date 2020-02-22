# 바우처 체인코드

요구사항
- 영어도서관 데이터 Invoke 및 Query
- 토큰 전송 및 보관

## OPP

|API Endpoints|Parameters|Description|
|---|---|--|
|`GET``http://169.56.126.28:8080/api/query/{SID}`|SID:학생번호|학생의 정보를 불러옵니다.|
|`POST``169.56.126.28:8080/api/invoke/create`|`{"data":[ "232170435",  "hyokeun kim",  "A",  "4340", "26", "55", "32", "0", "65", "44", "21"]}`| 학생의 학습 정보를 등록 및 업데트 합니다. `Sid: "232170431", FullName: "hyokeun kim", Level: "A", StarEarned: "4340", Logins: "26", Listen: "55", Read: "32", Worksheet: "0", Quiz: "65", PassedQuizCount: "44", PracticeRecording: "21"`|

API 결과 반환 값

``` json
{"response":"{\"FullName\":\"hyokeun kim\",\"Level\":\"A\",\"Listen\":\"55\" 
\"Logins\":\"26\",\"PassedQuizCount\":\"44\",\"PracticeRecording\":\"21\",
\"Quiz\":\"65\",\"Read\":\"32\",\"Sid\":\"232170433\",\"StarEarned\":\"4340\",
\"Worksheet\":\"0\"}"}
```

## Voucher (Credit)



## Voucher credit 제공량 연산을 위한 데이터들

### Asset

특이점 : struct 내에 정의한 항목은 대문자로 정의하지 않으면 Chaincode에서 반영되지 않는다.

- Lecture : 강의를 수강 완료했을 때 생성하는 Asset.

primary key : Sid+Courseid+"_1" 형태로 정의한다.

```go
type Lecture struct {
    Sid		string `json:Sid` // 학생 고유번호
	Courseid		string `json:Courseid` // 수업 고유번호
	Lecture_fin_date	[]int `json:Lecture_fin_date` // 강의를 마친 시간. unix timestamp
	Lecture_number		[]int `json:Lecture_number` // 해당 수업의 몇 번째 강의를 들었는지
	Focus_rate		[]float64 `json:Focus_rate` // focus 시간 / 전체 시간 비율
	Used		bool // 이미 Voucher 계산에 사용된 lecture인지 아닌지.
}
```

- CourseStatistics: 해당 Courseid에서 CalculateVoucher 함수에 사용할 각종 값을 저장하는 Asset.

primary key : Courseid
Lecture Asset의 primary key는 Sid+ "\_" + Courseid + "\_1" 형태로 정의한다.


* AvgFocusRate: hashmap 형태의 자료구조. key값으로는 lecture_number를, value값은 해당 lecture_number 수강 완료생의 평균 Focus rate 값.
* Count : hashmap 형태의 자료구조. key값은 lecture_number, value값은 해당 lecture_number를 수강한 학생의 수

집계할 값이 많아질 경우 추가 가능.

```go
type CourseStatistics struct {
	CourseId	string // 강좌 id
	AvgFocusRate	map[int]float64 // 강좌 내 각 강의별 (lecture_number) 모든 수강생의 Focus Rate 평균값
	Count		map[int]int // 강좌 내 각 강의별 (lecture_number) lecture 개수
 }
```


### ChainCode

용어 헷갈림 방지를 위해, '강좌' 와 '강의' 용어 정의
* 강좌: 한 강의에서 순서대로 정의된 '수업'
* 강의: 여러 개의 강좌를 포함하고 있는 하나의 완결된 '주제'

ex) 파이썬 웹 프로그래밍 '강의'가 있다면, 1번 - python이란? '강좌', 2번 - 웹 프로그래밍 개요 '강좌'가 있는 것.

강좌는 아래에서 lecture_number로, 강의는 CourseId에 대응된다.

*아래 함수들은 input 인자나 내부 로직이 변동될 가능성이 있음.*


웹 서버 측에서 블록체인에 데이터를 전송하는 함수. 웹 서버 측에서는 블록체인에 데이터를 전송해야 할 경우 이 함수만 사용할 수 있도록.
- dataToFabric(Sid, CourseId, Lecture_fin_date, Lecture_number, total_lectureTime, focus_lectureTime, 기타)
	* total_lecture_time : 사용자가 동영상을 재생한 총 시간. 유닉스 (Unix) 시간초를 string 형태로 전달
	* focus_lecture_time : blur로 처리되지 않은 재생시간. 유닉스 (Unix) 시간초를 string 형태로 전달

퀴즈가 있는 강의일 경우 "맞춘 정답 리스트나 개수" 등의 값도 추가할 수 있을 것으로 생각함

강좌 수강을 마친 뒤, 수강생의 해당 강의 수강정보를 저장하는 함수
- createLecture(Sid, CourseId, Lecture_fin_date, lecture_number, focus_rate)
	* 전부 string. data_to_fabric으로 처리한 데이터를 토대로 Lecture Asset을 생성하는 함수
	* CompositeKey의 작동여부는 20.02.22 기준으로 아직 미확인
	

강좌 수강을 마친 뒤 createLecture함수를 실행할 때, 해당 강의를 수강한 학생들이 만들어낸 데이터를 저장하는 함수
- updateCourseStatistics(Courseid, lecture_number, focus_rate)
	* createLecture 함수 마지막에 실행됨. 해당 Courseid의 CourseStatistics Asset을 생성하고, CourseStatistics의 AvgFocusRate값을 업데이트한다.

	* map 자료구조는 동시성을 지원하지 않는다고 해서, go 언어의 sync.RWMutex를 사용함
	참고: https://blog.golang.org/go-maps-in-action
	
모든
- CalculateVoucher(Sid, CouresId, 기타...) : 하나의 lecture를 마친 뒤 Voucher를 계산받는 함수.
	* 현재의 logic (변경 가능.)
		-	수강생(Sid)이 수강한 강의 (CourseId) 전체의 평균 AvgFocusRate가 해당 강의의 모든 수강생 평균 AvgFocusRate보다 클 경우 Voucher 1.2배 지급


(개별 강의 하나하나마다 만들지, 강의 전체를 듣고 강좌별로 작업할지 현재 분명하지 않아서, 일단 개별 강의 수강을 마친 뒤 함수 시행. 강의 전체를 듣고 해당 강좌를 대상으로 수행하려면, 각 강좌별 몇 개의 강의가 있는지 등등)