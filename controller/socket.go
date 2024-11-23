package controller

import (
	"encoding/json"
	"errors"
	"log"
	"main/service"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type SocketCtrl interface {
	TestSocket(c *gin.Context)
}

type SocketCtrlImpl struct {
	svc *service.AppService
}

var (
	clients  sync.Map // Thread-safe map for managing clients
	upgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	conn *websocket.Conn
	cmp  int64
	role int64
}

type SocketMessage struct {
	Type int    `json:"type"`
	Msg  string `json:"msg"`
}

// SendMessage sends a message to a specific client by user ID.
func SendMessage(userID int, msgType int, msg string) {
	if client, ok := clients.Load(userID); ok {
		c := client.(*Client)
		message := SocketMessage{
			Type: msgType,
			Msg:  msg,
		}
		if err := c.conn.WriteJSON(message); err != nil {
			log.Printf("Error sending message to client %d: %v", userID, err)
			clients.Delete(userID) // Clean up the client on error
		}
	}
}

func SendMessageCmpToDriver(cmpID int, msgType int, msg string) {
	message := SocketMessage{
		Type: msgType,
		Msg:  msg,
	}
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		if client.cmp == int64(cmpID) && client.role >= 300 {
			if err := client.conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message to client %d: %v", key, err)
				clients.Delete(key) // Clean up the client on error
			}
		}
		return true
	})
}

// SendMessageByCmp sends a message to all clients belonging to a specific company.
func SendMessageByCmp(cmpID int, msgType int, msg string) {
	message := SocketMessage{
		Type: msgType,
		Msg:  msg,
	}
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		if client.cmp == int64(cmpID) {
			if err := client.conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message to client %d: %v", key, err)
				clients.Delete(key) // Clean up the client on error
			}
		}
		return true
	})
}

// ParseUserDetails extracts user details from the JWT token.
func ParseUserDetails(ctrl *SocketCtrlImpl, tokenString string) (int, int16, int64, error) {
	if tokenString == "" {
		return 0, 0, 0, errors.New("missing token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("accessToken")), nil
	})
	if err != nil {
		return 0, 0, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, 0, errors.New("invalid token claims")
	}

	audience, err := claims.GetAudience()
	if err != nil {
		return 0, 0, 0, err
	}
	userID, err := strconv.Atoi(audience[0])
	if err != nil {
		return 0, 0, 0, err
	}

	info, err := ctrl.svc.UserServ.GetSeed(int64(userID))
	if err != nil {
		return 0, 0, 0, err
	}

	issuer, err := claims.GetIssuer()
	if err != nil || info.Seed.String != issuer {
		return 0, 0, 0, errors.New("invalid token issuer")
	}

	userInfo, err := ctrl.svc.UserServ.GetUserById(int64(userID))
	if err != nil {
		return 0, 0, 0, err
	}

	return userID, userInfo.Role, userInfo.Belongcmp, nil
}

// TestSocket handles new WebSocket connections and manages client registration.
func (s *SocketCtrlImpl) TestSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	var userID int
	defer func() {
		if userID > 0 {
			clients.Delete(userID)
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if string(msg) == "ping" {
			response := map[string]string{"type": "pong"}
			jsonResponse, _ := json.Marshal(response)
			if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
				log.Printf("Ping response error: %v", err)
				break
			}
			continue
		}

		id, role, cmp, err := ParseUserDetails(s, string(msg))
		if err != nil {
			log.Printf("Error parsing user details: %v", err)
			break
		}

		userID = id
		clients.Store(userID, &Client{conn: conn, cmp: cmp, role: int64(role)})
	}
}

func SocketCtrlInit(svc *service.AppService) *SocketCtrlImpl {
	return &SocketCtrlImpl{
		svc: svc,
	}
}
