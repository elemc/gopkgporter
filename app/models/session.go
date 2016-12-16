package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// Session is a struct for sessions
type Session struct {
	gorm.Model
	SessionID     string
	SessionUser   *User
	SessionUserID uint
	Expiration    time.Time
}

// CreateSession function create and return new session
func CreateSession() (session *Session) {
	session = new(Session)
	session.CreateNewUUID()
	session.Expiration = time.Now().Add(time.Hour * 336) // two week
	return
}

// CreateNewUUID function create and set new UUID in SessionID
func (s *Session) CreateNewUUID() {
	unix32bits := uint32(time.Now().Unix())
	buff := make([]byte, 32)
	numReader, err := rand.Read(buff)
	if numReader != len(buff) || err != nil {
		revel.ERROR.Printf("Error in CreateNewUUID: %s", err)
		return
	}

	s.SessionID = fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}
