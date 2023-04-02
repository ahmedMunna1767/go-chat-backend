package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uuid.UUID
	Conn *websocket.Conn
	Pool *Pool
	Wg   *sync.WaitGroup
}

func (c *Client) Read() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		switch messageType {
		case websocket.BinaryMessage:
			c.Wg.Add(1)
			go c.send(c.Wg, message)
		case websocket.CloseMessage:
			c.Pool.RemoveConnection(c)
			closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "ta-ta goodbye!")
			err := c.Conn.WriteMessage(websocket.CloseMessage, closeMessage)
			if err != nil {
				log.Println(err)
			}
			return
		case websocket.PingMessage:
			err := c.Conn.WriteMessage(websocket.PongMessage, message)
			if err != nil {
				log.Println(err)
				return
			}
		default:
		}
	}
}

func (c *Client) send(wg *sync.WaitGroup, message []byte) {
	defer wg.Done()

	if len(message) < 16 {
		log.Println("invalid receiver: length is less than 16")
		return
	}
	receiverID, err := uuid.FromBytes(message[:16])
	if err != nil {
		log.Println(err)
		return
	}
	payload := message[16:]
	receiver, err := c.Pool.GetClientByID(receiverID)
	if err != nil {
		log.Println(err)
		return
	}
	payload = append(c.ID[:], payload...)
	err = receiver.Conn.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		log.Println(err)
		return
	}
}
