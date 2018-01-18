package main

import (
	"github.com/lonnng/nano/component"
	"github.com/lonnng/nano/session"
	"./protocol"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"

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

