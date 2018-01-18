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
		w.LeaveRoom(session,nil)
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
	delete(w.players, s.ID())
	log.Println("Leave &v",s.ID())
	return w.Broadcast("leave",s.ID())
}

const totalCount = 5
var curRound = 0
var curRoundCommitCount = 0
var curQuestion protocol.Question
var curQuestionTimeUnix int64
var questions []protocol.Question
func (w *Room) StartCompetition(s * session.Session,msg *protocol.StartRequest) error {

	if curQuestion.Id != 0{
		log.Println("playing")
		log.Println(curQuestion)
		return nil
	}
	if w.ownPlayerId != s.ID(){
		log.Println("you aren't owner")
		return nil
	}
	go NewAQuestion(w)

	return nil
}

func NewAQuestion(w *Room) {
	curRound ++
	curRoundCommitCount = 0
	log.Println("Start")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	curRandom := r.Intn(len(questions))
	log.Println(curRandom)
	curQuestion = questions[curRandom]
	curQuestionTimeUnix = time.Now().Unix()
	w.Broadcast("questionNotify", &curQuestion)

	go CountDown(w)
}

func CountDown(w *Room) {
	nowRound := curRound
	time.Sleep(10 * time.Second)
	if nowRound == curRound && curRound < totalCount {
		NewAQuestion(w)
	}
}

func(w *Room) CommitAnswer(s *session.Session,playerAnswer *protocol.PlayerAnswer) error {
	log.Println("CommitAnswer &v", s.ID())
	if curQuestion.Id == 0 {
		return nil
	}
	if w.players[s.ID()].Id == 0 {
		return nil
	}

	notify := new(protocol.AnswerNotify)
	notify.Id = s.ID()
	isRight := curQuestion.IRightIdx == playerAnswer.AnswerIndex
	notify.IsRight = isRight
	if isRight{
		notify.Score = CalculateScore(curQuestionTimeUnix)
		w.players[s.ID()].Score += notify.Score
	}
	w.Broadcast("AnswerNotify",notify)


	curRoundCommitCount ++
	if curRound < totalCount{
		if curRoundCommitCount >=len(w.players) {
		go NewAQuestion(w)
	}
	}else {
		//结束了
		var maxScoreId int64 = 0
		var maxScore int64 = 0
		for _, value := range w.players {
			if value.Score > maxScore {
				maxScore = value.Score
				maxScoreId = value.Id
			}
		}

		w.Broadcast("endCompetition", &protocol.EndCompetition{Id:maxScoreId})
	}

	return nil
}

var totalTime int64= 10
func CalculateScore(lastTime int64) int64 {
	leftTime := totalTime - (time.Now().Unix() - lastTime)
	log.Println("left time &v",leftTime)
	score := 200 * leftTime / totalTime
	log.Println("score & v",score)

	return score
}