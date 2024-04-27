package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var clients = make(map[int]client)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	wsc *websocket.Conn
	cmp int64
}

type socketMsgT struct {
	Type int    `json:"type"`
	Msg  string `json:"msg"`
}

func SandByCmp(reciver int, msgType int, msg string) {

	for _, c := range clients {

		if c.cmp != int64(reciver) {
			break
		}

		msg := socketMsgT{
			Type: msgType,
			Msg:  msg,
		}

		if err := c.wsc.WriteJSON(msg); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func SandMsg(reciver int, msgType int, msg string) {

	if conn, ok := clients[reciver]; ok {

		msg := socketMsgT{
			Type: msgType,
			Msg:  msg,
		}

		if err := conn.wsc.WriteJSON(msg); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func HandleWS(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	newClient := client{
		wsc: conn,
		cmp: c.MustGet("belongCmp").(int64),
	}
	clients[c.MustGet("UserID").(int)] = newClient
	fmt.Println("New Conn: ", c.MustGet("UserID"))
	defer conn.Close()
	defer delete(clients, c.MustGet("UserID").(int))
	for {
	}
}
