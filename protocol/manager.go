package protocol

// LoginRequest represent a login request
type LoginRequest struct {
	Username  string `json:"username"`
	Cipher    string `json:"cipher"`
	Timestamp int    `json:"timestamp"`
}

// LoginResponse represent a login response
type LoginResponse struct {
	Status int    `json:"status"`
	ID     int64  `json:"id"`
	Error  string `json:"error"`
}

type StartRequest struct{
	ID int64
}

type Question struct {
	Id int `json:"id"`
	SDesc string `json:"sDesc"`
	SDesc1 string `json:"sDesc1"`
	SDesc2 string `json:"sDesc2"`
	SDesc3 string `json:"sDesc3"`
	SDesc4 string `json:"sDesc4"`
	IRightIdx int32 `json:"iRightIdx"`
}

type AnswerNotify struct {
	Id      int64  `json:"id"`
	IsRight bool `json:"isRight"`
	Score   int64  `json:"score"`
}

type PlayerAnswer struct {
	AnswerIndex int32 `json:"answerIndex"`
}

type Player struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Score int64  `json:"score"`
}

type PlayerLoginInfo struct {
	Name string `json:"name"`
}

type EndCompetition struct {
	Id int64 `json:"id"'`
}