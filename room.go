package main

import (
	"github.com/lonnng/nano/component"
	"github.com/lonnng/nano"
	"github.com/google/uuid"
	"github.com/lonnng/nano/session"
	"./protocol"
	"log"
	"time"
	"math/rand"
	"encoding/json"
)

type Room struct {
	ownPlayerId int64
	component.Base
	*nano.Group
	players map[int64]*protocol.Player
}

func NewRoom() *Room {
	return &Room{
		Group:nano.NewGroup(uuid.New().String()),
	}
}

func (w *Room) Init(){
	nano.OnSessionClosed(func(session *session.Session) {
		w.Leave(session)
		w.Broadcast("leave",&protocol.LeaveWorldResponse{ID:session.ID()})
	})
}

func (w *Room) Enter(s *session.Session,msg *protocol.PlayerLoginInfo) error{
	if w.Group.Count() == 0 {
		w.ownPlayerId = s.ID()
	}
	w.Group.Add(s)

	player:= new(protocol.Player)
	player.Id = s.ID()
	player.Name = msg.Name
	player.Score = 0
	w.players[s.ID()] = player

	log.Println("Enter &v",s.ID(),msg.Name)
	return s.Response(&protocol.EnterWorldResponse{ID:s.ID()})
}

func(w *Room) LeaveRoom(s *session.Session,msg []byte) error{
	w.Leave(s)
	if s.ID() == w.ownPlayerId {
		if len(w.Members())>0 {
			w.ownPlayerId = w.Members()[0]
		}
	}
	log.Println("Leave &v",s.ID())
	return w.Broadcast("leave",msg)
}

const totalCount = 5
var curQuestion protocol.Question
var curQuestionTimeUnix int64
var questions []protocol.Question
func (w *Room) StartCompetition(s * session.Session,msg *protocol.StartRequest) error {

	go StartRound(w)
	return nil
}

func StartRound(w *Room)  {
	log.Println("Start")

	for i := 0; i < totalCount; i++ {

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		curRandom := r.Intn(len(questions))
		log.Println(curRandom)
		curQuestion = questions[curRandom]
		curQuestionTimeUnix = time.Now().Unix()
		data, error := json.Marshal(curQuestion)
		if error == nil {
			log.Println("questionNotify")
			w.Broadcast("questionNotify", data)
		}
		time.Sleep(time.Second * 10)
	}

	log.Println("End")

}

func(w *Room) CommitAnswer(s *session.Session,playerAnswer *protocol.PlayerAnswer) error {
	log.Println("CommitAnswer &v", s.ID())
	if &curQuestion == nil {
		return nil
	}
	notify := new(protocol.AnswerNotify)
	notify.Id = s.ID()
	isRight := curQuestion.IRightIdx == playerAnswer.AnswerIndex
	notify.IsRight = isRight
	if isRight{

	}

	notify.Score = CalculateScore(curQuestionTimeUnix)
	return nil
	//
}

var totalTime int64= 10
func CalculateScore(lastTime int64) int64 {
	leftTime := totalTime - (time.Now().Unix() - lastTime)
	score := 200 * leftTime / totalTime
	return score
}