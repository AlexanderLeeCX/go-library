/**
 * @Author: Lee
 * @Description:
 * @File:  connect
 * @Version: 1.0.0
 * @Date: 2021/11/4 9:47 下午
 */

package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Data struct {
	Ip       string   `json:"ip"`
	Room     string   `json:"room"`
	User     string   `json:"user"`
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
}

type connection struct {
	ws   *websocket.Conn
	sc   chan []byte
	data *Data
}

var wu = &websocket.Upgrader{ReadBufferSize: 512,
	WriteBufferSize: 512, CheckOrigin: func(r *http.Request) bool { return true }}

func myws(w http.ResponseWriter, r *http.Request) {
	ws, err := wu.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{sc: make(chan []byte, 256), ws: ws, data: &Data{}}
	rChan <- c

	go c.writer()
	c.reader()
	defer func() {
		h := hubMap[c.data.Room]
		if h != nil {
			c.data.Type = "logout"
			h.userList = del(h.userList, c.data.User)
			c.data.UserList = h.userList
			c.data.Content = c.data.User
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
			rChan <- c
			h.u <- c
		}
	}()
}

func (c *connection) writer() {
	for message := range c.sc {
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	c.ws.Close()
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		json.Unmarshal(message, &c.data)
		fmt.Println(string(message))
		h := hubMap[c.data.Room]
		if h == nil {
			h = NewHub()
			hubMap[c.data.Room] = h
			go run(h)
		}
		switch c.data.Type {
		case "login":
			h.c[c] = true
			c.data.User = c.data.Content
			c.data.From = c.data.User
			h.userList = append(h.userList, c.data.User)
			c.data.UserList = h.userList
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
		case "user":
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
		case "logout":
			c.data.Type = "logout"
			h.userList = del(h.userList, c.data.User)
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
			rChan <- c
		default:
			fmt.Print("========default================")
		}
	}
}

func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(n_slice)
	return n_slice
}

func run(h *hub) {
	for {
		select {
		case c := <-h.u:
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.sc)
			}
		case data := <-h.b:
			for c := range h.c {
				select {
				case c.sc <- data:
				default:
					delete(h.c, c)
					close(c.sc)
				}
			}
		}
	}
}
