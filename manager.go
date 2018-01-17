package main

import (
	"github.com/lonnng/nano/component"
	"github.com/lonnng/nano/session"
	"./protocol"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"math/rand"
	"time"
	"log"
)

type Manager struct {
	component.Base
}

func NewManager() *Manager{
	return &Manager{}
}

func (m *Manager)LoadQuestions() error {
	//准备题目
	raw, error := ioutil.ReadFile("./question.json")
	if error != nil {
		fmt.Println(error.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &questions)
	return nil
}


func (m *Manager) Login(s *session.Session,msg *protocol.LoginRequest) error {

	id := s.ID()
	s.Bind(id)

	return s.Response(protocol.LoginResponse{
		Status: protocol.LoginStatusSucc,
		ID:     id,
	})
}

const totalCount = 5
var questions []protocol.Question
func (m *Manager) StartCompetition(s * session.Session,msg *protocol.StartRequest) error {


	go StartRound(1)
	return nil
}

func (m *Manager) ShowQuestions(){
	log.Println("hello")
	go StartRound(1)
}

func StartRound(roundIndex int32) error {
	if roundIndex > totalCount {
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	curRandom := r.Intn(len(questions))
	log.Println(curRandom)
	q:=questions[curRandom]

	log.Println(q.SDesc)

	time.Sleep(time.Second * 10)

	roundIndex +=1
	go StartRound(roundIndex)

	return nil
}