package utils

// import (
// 	"errors"
// 	"fmt"
// 	"main/service"
// 	"net/http"
// 	"os"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"

// 	"github.com/golang-jwt/jwt/v5"
// )

// var clients = make(map[int]client)

// var upgrader = websocket.Upgrader{
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// type client struct {
// 	wsc *websocket.Conn
// 	cmp int64
// }

// type socketMsgT struct {
// 	Type int    `json:"type"`
// 	Msg  string `json:"msg"`
// }

// func SandByCmp(reciver int, msgType int, msg string) {

// 	for _, c := range clients {

// 		if c.cmp != int64(reciver) {
// 			break
// 		}

// 		msg := socketMsgT{
// 			Type: msgType,
// 			Msg:  msg,
// 		}

// 		if err := c.wsc.WriteJSON(msg); err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}
// }

// func SandMsg(reciver int, msgType int, msg string) {

// 	if conn, ok := clients[reciver]; ok {

// 		msg := socketMsgT{
// 			Type: msgType,
// 			Msg:  msg,
// 		}

// 		if err := conn.wsc.WriteJSON(msg); err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}
// }

// func HandleWS(c *gin.Context) {

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// fmt.Println("New Conn: ", c.MustGet("UserID"))
// 	defer conn.Close()
// 	// defer delete(clients, c.MustGet("UserID").(int))
// 	for {
// 		_, msg, err := conn.ReadMessage()

// 		if err != nil {
// 			return
// 		}
// 		id, _, cmp, err := getUD(string(msg))
// 		if err != nil {
// 			fmt.Print(err)
// 			return
// 		}
// 		fmt.Println(id)
// 		fmt.Println(cmp)

// 		newClient := client{
// 			wsc: conn,
// 			cmp: cmp,
// 		}
// 		clients[id] = newClient
// 	}
// }

// func getUD(rtoken string) (int, int16, int64, error) {
// 	fmt.Print(rtoken)
// 	if rtoken == "" {
// 		return 0, 0, 0, errors.New("1")

// 	}

// 	token, err := jwt.Parse(rtoken, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("accessToken")), nil
// 	})

// 	if err != nil {
// 		fmt.Print(err)
// 		return 0, 0, 0, errors.New("2")
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 		res, err := claims.GetAudience()

// 		if err != nil {
// 			return 0, 0, 0, errors.New("3")
// 		}

// 		id, err := strconv.Atoi(res[0])
// 		if err != nil {

// 			return 0, 0, 0, errors.New("4")
// 		}
// 		info, err := service.UserServ.GetSeed(&service.UserServImpl{}, int64(id))
// 		if err != nil {
// 			fmt.Print("QQ")
// 			return 0, 0, 0, errors.New("5")
// 		}

// 		issuer, err := claims.GetIssuer()
// 		if err != nil {

// 			return 0, 0, 0, errors.New("6")
// 		}
// 		if info.String != issuer {

// 			return 0, 0, 0, errors.New("7")
// 		}

// 		userInfo, err := service.UserServ.GetUserById( int64(id))
// 		if err != nil {

// 			return 0, 0, 0, errors.New("8")
// 		}

// 		return id, userInfo.Role, userInfo.Belongcmp, nil

// 	} else {
// 		return 0, 0, 0, errors.New("9")
// 	}

// }
