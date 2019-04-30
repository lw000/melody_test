package main

import (
	"demo/melody_test/test/config"
	"demo/melody_test/ws"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	cfg *config.JsonConfig
)

func chatMsg(uid int) ([]byte, error) {
	suid := strconv.Itoa(uid)
	m := make(map[string]interface{})
	m["uid"] = suid
	m["msg"] = strings.Repeat(suid, 2)
	m["group_id"] = 1

	data, er := json.Marshal(m)
	if er != nil {
		log.Println(er)
		return nil, er
	}

	return data, nil
}

func TestMessage(fc *ws.FastWsClient, uid int) {
	tickHeartBeat := time.NewTicker(time.Second * time.Duration(45))
	tickSend := time.NewTicker(time.Millisecond * time.Duration(cfg.Millisecond))
	for {
		select {
		case <-tickHeartBeat.C:
			er := fc.Ping()
			if er != nil {
				log.Println(er)
				return
			}
		case <-tickSend.C:
			data, er := chatMsg(uid)
			if er != nil {
				log.Println(er)
				return
			}

			er = fc.SendMessage(data)
			if er != nil {
				log.Println(er)
				return
			}
		}
	}
}

func main() {
	var er error
	cfg, er = config.LoadJsonConfig("./conf/conf.json")
	if er != nil {
		log.Panic(er)
	}

	for i := 1; i <= cfg.Count; i++ {
		cli := &ws.FastWsClient{}
		cli.HandleConnected(func() {
			log.Printf("connected [uid=%d]", i)
		})

		cli.HandleDisConnected(func() {
			log.Println("disconnected")
		})

		cli.HandleMessage(func(data []byte) {
			log.Println(string(data))
		})
		er = cli.Create(cfg.Scheme, cfg.Host, cfg.Path)
		if er != nil {
			log.Println(er)
			continue
		}

		if cfg.Send {
			go TestMessage(cli, i)
			go cli.Run()
		}
		time.Sleep(time.Microsecond * time.Duration(20))
	}
	select {}
}
