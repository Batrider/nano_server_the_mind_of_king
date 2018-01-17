package main

import (
	"github.com/lonnng/nano/component"
	"github.com/lonnng/nano"
	"github.com/google/uuid"
	"github.com/lonnng/nano/session"
	"./protocol"
	"log"
)

type Room struct {
	ownPlayerId int64
	component.Base
	*nano.Group
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

var playerId int64 = 0

func (w *Room) Enter(s *session.Session,msg []byte) error{
	w.Group.Add(s)
	log.Println("Enter &v",s.ID())
	playerId ++
	return s.Response(&protocol.EnterWorldResponse{ID:s.ID() + playerId})
}

func(w *Room) CommitAnswer(s *session.Session,msg []byte) error{
	log.Println("Update &v",s.ID())

	return w.Broadcast("update",msg)

}

func(w *Room) LeaveRoom(s *session.Session,msg []byte) error{
	w.Leave(s)
	log.Println("Leave &v",s.ID())

	return w.Broadcast("leave",msg)
}