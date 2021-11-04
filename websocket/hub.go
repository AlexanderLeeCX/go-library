/**
 * @Author: Lee
 * @Description:
 * @File:  hub
 * @Version: 1.0.0
 * @Date: 2021/11/4 9:47 下午
 */

package websocket

import "encoding/json"

type hub struct {
	c        map[*connection]bool
	b        chan []byte
	u        chan *connection
	userList []string
}

func NewHub() *hub {
	return &hub{
		c:        make(map[*connection]bool),
		u:        make(chan *connection),
		b:        make(chan []byte),
		userList: []string{},
	}
}

var (
	hubMap = make(map[string]*hub)
	rChan  = make(chan *connection)
)

func runChan() {
	for {
		select {
		case c := <-rChan:
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			//c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			c.sc <- data_b
		}
	}
}
