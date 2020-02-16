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

- Lecture

강의를 수강 완료했을 때 생성하는 Asset.

```go
type Lecture struct {
     Sid               string `json:Sid` // 학생 고유번호
	 Courseid          string `json:Courseid` // 수업 고유번호
	 Lecture_fin_date  []int `json:Lecture_fin_date` // 강의를 마친 시간. unix timestamp
	 Lecture_number    []int	`json:Lecture_number` // 해당 수업의 몇 번째 강의를 들었는지
	 Focus_rate        []float64 `json:Focus_rate` // focus 시간 / 전체 시간 비율
}
```

Lecture Asset의 primary key는 Sid+Courseid+"_1" 형태로 정의한다.



### ChainCode

input 인자 변동 가능성 있음

- data_to_fabric(Sid, CourseId, Lecture_fin_date, Lecture_number, total_lectureTime, focus_lectureTime, 기타)

웹 서버 측에서 블록체인에 데이터를 전송하는 함수. 웹 서버 측에서는 블록체인에 데이터를 전송해야 할 경우 이 함수만 사용할 수 있도록.

	* total_lecture_time : 사용자가 동영상을 재생한 총 시간. 유닉스 (Unix) 시간초를 string 형태로 전달
	* focus_lecture_time : blur로 처리되지 않은 재생시간. 유닉스 (Unix) 시간초를 string 형태로 전달

퀴즈가 있는 강의일 경우 "맞춘 정답 리스트나 개수" 등의 값도 추가할 수 있을 것으로 생각함

- createLecture(Sid, CourseId, Lecture_fin_date, lecture_number, focus_rate)
	* 전부 string. data_to_fabric으로 처리한 데이터를 토대로 Lecture Asset을 생성하는 함수

