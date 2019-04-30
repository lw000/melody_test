// melody_test project main.go
package main

import (
	"demo/melody_test/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"log"
	"net/http"
)

type ChatData struct {
	Uid     string `json:"uid"`
	Msg     string `json:"msg"`
	GroupId int64  `json:"group_id"`
}

func main() {
	engine := gin.Default()
	m := melody.New()
	engine.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "templates/index.html")
	})

	engine.GET("/ws", func(c *gin.Context) {
		er := m.HandleRequest(c.Writer, c.Request)
		if er != nil {
			log.Println(er)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		uuid := utils.UUID()
		s.Set("uuid", uuid)
		log.Println(uuid, s.Request.RemoteAddr, "join")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var rdata ChatData
		if er := json.Unmarshal(msg, &rdata); er != nil {
			er = s.CloseWithMsg([]byte("error"))
			if er != nil {
				log.Println(er)
			}
			return
		}
		var er error
		//er = m.BroadcastFilter(msg, func(session *melody.Session) bool {
		//	return s == session
		//})

		er = m.Broadcast(msg)
		if er != nil {
			log.Println(string(msg))
		}
	})

	m.HandleDisconnect(func(s *melody.Session) {
		v, exists := s.Get("uuid")
		if exists {
			log.Println(v, s.Request.RemoteAddr, "leave")
		}
	})

	log.Panic(engine.RunTLS(":6789", "crs/server.crt", "crs/server.key"))
}
