/**
 * @Author: Lee
 * @Description:
 * @File:  server
 * @Version: 1.0.0
 * @Date: 2021/11/4 9:47 下午
 */

package websocket

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func StartServer() {

	router := mux.NewRouter()
	go runChan()
	router.HandleFunc("/ws", myws)
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}
